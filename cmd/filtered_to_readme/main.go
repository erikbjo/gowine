package main

import (
	"encoding/json"
	"gowine/internal/shared"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Take the filtered products from the scrape process and write them to the README
	// in a nice format

	// Load gowine products from JSON
	file, err := os.Open("json/gowine_products.json")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Read the products
	var filteredProducts []*shared.Product
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&filteredProducts)
	if err != nil {
		log.Fatalf("Failed to decode gowine products: %s", err)
	}

	// Write the products to the README
	readme, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("Failed to create README: %s", err)
	}
	defer readme.Close()

	// Sort the products by discount
	sort.Slice(filteredProducts, func(i, j int) bool {
		return filteredProducts[i].Discount > filteredProducts[j].Discount
	})

	// Write the header
	readme.WriteString("# Go(d) wine\n\n")
	readme.WriteString("Liste over produkter som blir billigere etter månedsskiftet:\n\n")

	// Write the products
	// One table for white wines, one for red wines, and one for everything else
	readme.WriteString("## Hvite viner\n\n")
	readme.WriteString("| Navn | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Land |\n")
	readme.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, product := range filteredProducts {
		if product.Type == "Hvitvin" {
			info := "| " + product.Basic.ProductShortName + " | " + product.Basic.ProductId + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.Difference) + " | " + strconv.Itoa(product.Discount) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Country + " |\n"
			readme.WriteString(info)
		}
	}

	readme.WriteString("## Røde viner\n\n")
	readme.WriteString("| Navn | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Land |\n")
	readme.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, product := range filteredProducts {
		if product.Type == "Rødvin" {
			info := "| " + product.Basic.ProductShortName + " | " + product.Basic.ProductId + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.Difference) + " | " + strconv.Itoa(product.Discount) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Country + " |\n"
			readme.WriteString(info)
		}
	}

	readme.WriteString("## Andre produkter\n\n")
	readme.WriteString("| Navn | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Type |\n")
	readme.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, product := range filteredProducts {
		if product.Type != "Rødvin" && product.Type != "Hvitvin" {
			info := "| " + product.Basic.ProductShortName + " | " + product.Basic.ProductId + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.Difference) + " | " + strconv.Itoa(product.Discount) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Type + " |\n"
			readme.WriteString(info)
		}
	}
}
