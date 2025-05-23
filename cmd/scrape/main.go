package main

import (
	"encoding/json"
	"gowine/internal/apertif"
	"gowine/internal/shared"
	"gowine/internal/vinmonopolet"
	"log"
	"os"
	"slices"
	"sync"

	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}
}

func main() {
	products := vinmonopolet.GetWines()
	if len(products) == 0 {
		log.Fatal("No products retrieved. Exiting.")
	}

	var wg sync.WaitGroup
	var scrapedProducts []*shared.Product
	var validMutex sync.Mutex

	var expiredProducts []*shared.Product
	var expiredMutex sync.Mutex

	// Load expired products from JSON, if it exists
	file, err := os.Open("json/expired_products.json")
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&expiredProducts)
		if err != nil {
			log.Fatalf("Failed to decode expired products: %s", err.Error())
		}
		err = file.Close()
		if err != nil {
			log.Printf("Failed to close file: %s", err.Error())
		}
	}

	// Log expired products
	if len(expiredProducts) > 0 {
		log.Printf("Found %d expired products, filtering...", len(expiredProducts))
		log.Printf("Amount of products before filtering: %d", len(products))
		for _, product := range expiredProducts {
			for i, p := range products {
				if p.Basic.ProductId == product.Basic.ProductId {
					products = slices.Delete(products, i, i+1)
					break
				}
			}
		}
		log.Printf("Amount of products after filtering: %d", len(products))
	}

	// Load scraped products from JSON, if it exists
	file, err = os.Open("json/scraped_products.json")
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&scrapedProducts)
		if err != nil {
			log.Fatalf("Failed to decode scraped products: %s", err.Error())
		}
		err = file.Close()
		if err != nil {
			log.Printf("Failed to close file: %s", err.Error())
		}
	}

	// Log scraped products
	if len(scrapedProducts) > 0 {
		log.Printf("Found %d pre-scraped products, adding to scraped...", len(scrapedProducts))
	}

	log.Printf("Starting to scrape %d products", len(products))

	loadingBar := progressbar.Default(int64(len(products)))

	// Limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 20)

	for _, product := range products {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(wine shared.Product) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Scrape data from both sources
			vinmonopolet.ScrapeVinmonopolet(&wine)
			apertif.ScrapeApertif(&wine)

			// Too many requests to vivino, TODO: fix, delay or something
			// vivino.ScrapeVivino(&wine)

			if wine.VinmonopoletPrice == -1 {
				//log.Printf("%s: product expired, check expired_products.json", wine.Basic.ProductId)
				expiredMutex.Lock()
				expiredProducts = append(expiredProducts, &wine)
				expiredMutex.Unlock()
			} else {
				validMutex.Lock()
				scrapedProducts = append(scrapedProducts, &wine)
				validMutex.Unlock()

				//log.Printf("%s: finished scraping", wine.Basic.ProductId)
				loadingBar.Add(1)
			}
		}(product)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Printf("All scraping done, processing results.")

	// Filter products with complete pricing
	filteredProducts := filterCompleteProducts(scrapedProducts)
	priceDifferenceProducts := filterPriceDifferences(filteredProducts)

	// Save results to JSON
	saveToJSON("json/scraped_products.json", priceDifferenceProducts)

	// Save expired products to JSON
	saveToJSON("json/expired_products.json", expiredProducts)

	log.Printf("Saved %d products with price differences to scraped_products.json", len(priceDifferenceProducts))
	log.Printf("Saved %d expired products to expired_products.json", len(expiredProducts))
}

// Filters products that have valid prices from both sources
func filterCompleteProducts(products []*shared.Product) []*shared.Product {
	var filtered []*shared.Product
	for _, product := range products {
		if product.VinmonopoletPrice != 0 && product.ApertifPrice != 0 {
			filtered = append(filtered, product)
		} else {
			log.Printf("Product %s, art.nr %s has missing prices, skipping", product.Basic.ProductShortName, product.Basic.ProductId)
		}
	}
	return filtered
}

// Filters products with a price difference
func filterPriceDifferences(products []*shared.Product) []*shared.Product {
	var filtered []*shared.Product
	for _, product := range products {
		if product.VinmonopoletPrice != product.ApertifPrice {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

// Saves products to a JSON file
func saveToJSON(filename string, products []*shared.Product) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(products); err != nil {
		log.Fatalf("Failed to encode products to JSON: %s", err)
	}
}
