package main

import (
	"encoding/json"
	"gowine/internal/shared"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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

	// Prepare the README
	readme, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("Failed to create README: %s", err)
	}
	defer readme.Close()

	// Sort the products by discount
	sort.Slice(filteredProducts, func(i, j int) bool {
		return filteredProducts[i].Discount > filteredProducts[j].Discount
	})

	whiteWines := strings.Builder{}
	redWines := strings.Builder{}
	otherProducts := strings.Builder{}

	whiteWines.WriteString("## Hvite viner\n\n| Navn | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Land |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	redWines.WriteString("## Røde viner\n\n| Navn | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Land |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	otherProducts.WriteString("## Andre produkter\n\n| Navn | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Type |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n")

	for _, product := range filteredProducts {
		switch product.Type {
		case "Hvitvin":
			whiteWines.WriteString("| " + product.Basic.ProductShortName + " | " + product.GetVinmonopoletMarkdownUrl() + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.Difference) + " | " + strconv.Itoa(product.Discount) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Country + " |\n")
		case "Rødvin":
			redWines.WriteString("| " + product.Basic.ProductShortName + " | " + product.GetVinmonopoletMarkdownUrl() + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.Difference) + " | " + strconv.Itoa(product.Discount) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Country + " |\n")
		default:
			otherProducts.WriteString("| " + product.Basic.ProductShortName + " | " + product.GetVinmonopoletMarkdownUrl() + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.Difference) + " | " + strconv.Itoa(product.Discount) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Type + " |\n")
		}
	}

	// Write the header
	readme.WriteString("# Go(d) wine\n\n")
	readme.WriteString("Liste over produkter som blir billigere etter månedsskiftet:\n\n")

	readme.WriteString(whiteWines.String())
	readme.WriteString(redWines.String())
	readme.WriteString(otherProducts.String())
}