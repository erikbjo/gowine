package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"gowine/internal/apertif"
	"gowine/internal/shared"
	"gowine/internal/vinmonopolet"
	"log"
	"os"
	"sync"
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
	var mutex sync.Mutex

	log.Printf("Starting to scrape %d products", len(products))

	// Limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 10)

	for _, product := range products {
		semaphore <- struct{}{} // Reserve a slot
		wg.Add(1)

		go func(wine shared.Product) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the slot

			log.Printf("Starting to scrape %s, art.nr %s", wine.Basic.ProductShortName, wine.Basic.ProductId)

			// Scrape data from both sources
			vinmonopolet.ScrapeVinmonopolet(&wine)
			apertif.ScrapeApertif(&wine)

			// Safely append results
			mutex.Lock()
			scrapedProducts = append(scrapedProducts, &wine)
			mutex.Unlock()

			log.Printf("Finished scraping %s", wine.Basic.ProductShortName)
		}(product)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Printf("All scraping done, processing results.")

	// Filter products with complete pricing
	filteredProducts := filterCompleteProducts(scrapedProducts)
	priceDifferenceProducts := filterPriceDifferences(filteredProducts)

	// Save results to JSON
	saveToJSON("scraped_products.json", priceDifferenceProducts)
	fmt.Println("Scraping and saving to JSON completed")
}

// Filters products that have valid prices from both sources
func filterCompleteProducts(products []*shared.Product) []*shared.Product {
	var filtered []*shared.Product
	for _, product := range products {
		if product.VinmonopoletPrice != 0 && product.ApertifPrice != 0 {
			filtered = append(filtered, product)
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
