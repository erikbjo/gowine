package main

import (
	"encoding/json"
	"gowine/internal/shared"
	"log"
	"os"
)

// Process scraped data
func main() {
	var scrapedProducts []*shared.Product

	// Load scraped products from JSON, if it exists
	file, err := os.Open("scraped_products.json")
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&scrapedProducts)
		if err != nil {
			log.Fatalf("Failed to decode scraped products: %s", err)
		}
		file.Close()
	}

	// Filter for products which has a lower apertif price than vinmonopolet price
	// Price needs to be 400 NOK or less, and discount needs to be 20% or more
	// This is the concept of gowine
	var gowineProducts []*shared.Product
	for _, product := range scrapedProducts {
		if product.ApertifPrice <= 400 {
			discount := 100 - (product.ApertifPrice * 100 / product.VinmonopoletPrice)
			if discount >= 40 {
				gowineProducts = append(gowineProducts, product)
			}
		}
	}

	// Save gowine products to JSON
	file, err = os.Create("gowine_products.json")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(gowineProducts)
	if err != nil {
		log.Fatalf("Failed to encode gowine products: %s", err)
	}

	file.Close()
}