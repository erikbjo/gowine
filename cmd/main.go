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
	products := vinmonopolet.GetWines("100", "1200")
	for _, product := range products {
		log.Println(product.GetVinmonopoletUrl())
	}
}
