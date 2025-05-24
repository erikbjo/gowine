package main

import (
	"encoding/json"
	"fmt"
	"gowine/internal/apertif"
	"gowine/internal/shared"
	"gowine/internal/vinmonopolet"
	"os"
	"slices"
	"sync"

	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

var logger = shared.CreateSugaredLogger()

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading .env file: %s", err.Error())
	}
}

func main() {
	products, err := vinmonopolet.GetWines()
	if len(products) == 0 {
		logger.Fatal("Got zero wines from vinmonopolet")
	}
	if err != nil {
		logger.Fatalf("Could not get vinmonopolet wines: %s", err.Error())
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
			logger.Fatalf("Failed to decode expired products: %s", err.Error())
		}
		err = file.Close()
		if err != nil {
			logger.Warnf("Failed to close file: %s", err.Error())
		}
	}

	// Log expired products
	if len(expiredProducts) > 0 {
		logger.Infof("Found %d expired products, filtering...", len(expiredProducts))
		logger.Infof("Amount of products before filtering: %d", len(products))
		for _, product := range expiredProducts {
			for i, p := range products {
				if p.Basic.ProductId == product.Basic.ProductId {
					products = slices.Delete(products, i, i+1)
					break
				}
			}
		}
		logger.Infof("Amount of products after filtering: %d", len(products))
	}

	// Load scraped products from JSON, if it exists
	file, err = os.Open("json/scraped_products.json")
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&scrapedProducts)
		if err != nil {
			logger.Fatalf("Failed to decode scraped products: %s", err.Error())
		}
		err = file.Close()
		if err != nil {
			logger.Infof("Failed to close file: %s", err.Error())
		}
	}

	// TODO: Should not read pre-scraped, but maybe move it to json/log/month dir
	// Log scraped products
	if len(scrapedProducts) > 0 {
		logger.Infof("Found %d pre-scraped products, adding to scraped...", len(scrapedProducts))
	}

	logger.Infof("Starting to scrape %d products", len(products))

	loadingBar := progressbar.Default(int64(len(products)))

	// Limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 25)

	for _, product := range products {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(wine shared.Product) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// TODO: move functions to structs
			vinmonopolet.ScrapeVinmonopolet(&wine, 0)
			apertif.ScrapeApertif(&wine, 0)

			// Too many requests to vivino, TODO: fix, delay or something
			// vivino.ScrapeVivino(&wine)

			if wine.VinmonopoletPrice == -1 {
				expiredMutex.Lock()
				expiredProducts = append(expiredProducts, &wine)
				expiredMutex.Unlock()
			} else {
				validMutex.Lock()
				scrapedProducts = append(scrapedProducts, &wine)
				validMutex.Unlock()
			}
			_ = loadingBar.Add(1)
		}(product)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Infof("All scraping done, processing results.")

	// Filter products with complete pricing
	filteredProducts := filterCompleteProducts(scrapedProducts)
	priceDifferenceProducts := filterPriceDifferences(filteredProducts)

	// Save results to JSON
	err = saveToJSON("json/scraped_products.json", priceDifferenceProducts)
	if err != nil {
		logger.Infof("Failed to save scraped products to json: %s", err.Error())
	}

	// Save expired products to JSON
	err = saveToJSON("json/expired_products.json", expiredProducts)
	if err != nil {
		logger.Infof("Failed to save expired products to json: %s", err.Error())
	}

	logger.Infof("Saved %d products with price differences to scraped_products.json", len(priceDifferenceProducts))
	logger.Infof("Saved %d expired products to expired_products.json", len(expiredProducts))
}

// Filters products that have valid prices from both sources
func filterCompleteProducts(products []*shared.Product) []*shared.Product {
	var filtered []*shared.Product
	for _, product := range products {
		if product.VinmonopoletPrice != 0 && product.ApertifPrice != 0 {
			filtered = append(filtered, product)
		} else {
			logger.Infof("Product %s, art.nr %s has missing prices, skipping", product.Basic.ProductShortName, product.Basic.ProductId)
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
func saveToJSON(filename string, products []*shared.Product) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", filename, err.Error())
	}

	defer func() {
		err := file.Close()
		if err != nil {
			logger.Warnf("Failed to close file: %s", err.Error())
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(products); err != nil {
		return fmt.Errorf("failed to encode products to file %s: %s", filename, err.Error())
	}

	return nil
}
