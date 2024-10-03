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
		log.Fatal("Error loading .env file" + err.Error())
	}
}

func main() {
	products := vinmonopolet.GetWines()

	var wg sync.WaitGroup
	resultChan := make(chan *shared.Product, len(products))

	var scrapedProducts []*shared.Product
	log.Printf("Starting to scrape %d products", len(products))

	// Limit the number of concurrent goroutines
	semaphore := make(chan struct{}, 10)

	for _, product := range products {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(wine *shared.Product) {
			defer wg.Done()
			defer func() {
				<-semaphore
				log.Printf("Scraped %s", wine.Basic.ProductShortName)
			}()
			defer func() { resultChan <- wine }() // Ensure the result is always sent

			log.Printf("Starting to scrape %s, art.nr %s", wine.Basic.ProductShortName, wine.Basic.ProductId)
			vinmonopolet.ScrapeVinmonopolet(wine, resultChan)
			apertif.ScrapeApertif(wine, resultChan)
		}(&product)
	}

	// Close the result channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
		log.Printf("All goroutines are done")
	}()

	// Collect results
	for wine := range resultChan {
		scrapedProducts = append(scrapedProducts, wine)
	}

	// Filter out products that are missing prices
	var filteredProducts []*shared.Product
	for _, product := range scrapedProducts {
		if product.VinmonopoletPrice != 0 && product.ApertifPrice != 0 {
			filteredProducts = append(filteredProducts, product)
			log.Printf("Added %s to filtered products", product.Basic.ProductShortName)
		}
	}

	scrapedProducts = filteredProducts

	// Filter for only wines with a price difference
	var priceDifferenceProducts []*shared.Product
	for _, product := range scrapedProducts {
		if product.VinmonopoletPrice != product.ApertifPrice {
			priceDifferenceProducts = append(priceDifferenceProducts, product)
		}
	}

	scrapedProducts = priceDifferenceProducts

	// Save results to JSON
	file, err := os.Create("scraped_products.json")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(scrapedProducts); err != nil {
		log.Fatalf("Failed to encode products to JSON: %s", err)
	}

	fmt.Println("Scraping and saving to JSON completed")
}
