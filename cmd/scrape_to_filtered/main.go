package main

import (
	"encoding/json"
	"gowine/internal/shared"
	"os"

	"go.uber.org/zap"
)

var logger = shared.CreateSugaredLogger()

// Process scraped data
func main() {
	var scrapedProducts []*shared.Product

	// Load scraped products from JSON, if it exists
	file, err := os.Open("json/scraped_products.json")
	if err != nil {
		logger.Fatal("Failed to open file", zap.String("file", file.Name()))
	}
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&scrapedProducts)
		if err != nil {
			logger.Fatal("Failed to decode scraped products", zap.Error(err))
		}
		err = file.Close()
		if err != nil {
			logger.Warn("Failed to close file", zap.Error(err))
		}
	}

	// Filter for products which has a lower apertif price than vinmonopolet price
	// Price needs to be x00 NOK or less, and discount needs to be 20% or more
	// This is the concept of gowine
	var gowineProducts []*shared.Product
	for _, product := range scrapedProducts {
		if product.ApertifPrice <= 1000 && product.GetAbsolutePercentDifference() >= 20 {
			gowineProducts = append(gowineProducts, product)
		}
	}

	// Save gowine products to JSON
	file, err = os.Create("json/gowine_products.json")
	if err != nil {
		logger.Fatal("Failed to create file", zap.Error(err))
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(gowineProducts)
	if err != nil {
		logger.Error("Failed to encode gowine products", zap.Error(err))
	}

	err = file.Close()
	if err != nil {
		logger.Warn("Failed to close file", zap.Error(err))
	}
}
