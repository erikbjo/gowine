package main

import (
	"github.com/joho/godotenv"
	"gowine/vinmonopolet"
	"log"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file" + err.Error())
	}
}

func main() {
	// Testing
	products := vinmonopolet.GetWines()

	// Test first
	if len(products) > 0 {
		wine := products[0]
		log.Println(wine)
		vinmonopolet.ScrapeVinmonopolet(&wine)
		log.Println(wine)
	}

	// Flow:
	// 1. Fetch wines from Vinmonopolet
	// 2. Scrape Vinmonopolet for more details
	// 3. Filter and validate wines
	// 5. Scrape Apertif for prices
	// 6. Save wines to database

	// On next run, only update wines missing prices
}
