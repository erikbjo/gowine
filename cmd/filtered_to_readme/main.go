package main

import (
	"encoding/json"
	"gowine/internal/shared"
	"os"
	"sort"
	"strconv"
	"strings"
)

var logger = shared.CreateSugaredLogger()

func main() {
	// Take the filtered products from the scrape process and write them to the README
	// in a nice format

	// Load gowine products from JSON
	file, err := os.Open("json/gowine_products.json")
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Warnf("Failed to close gowine_products file: %s", err.Error())
		}
	}()

	// Read the products
	var filteredProducts []*shared.Product
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&filteredProducts)
	if err != nil {
		logger.Fatalf("Failed to decode gowine products: %s", err.Error())
	}

	// Prepare the README
	readme, err := os.Create("README.md")
	if err != nil {
		logger.Fatalf("Failed to create README: %s", err.Error())
	}
	defer func() {
		err = readme.Close()
		if err != nil {
			logger.Warnf("Failed to close readme file: %s", err.Error())
		}
	}()

	// Sort the products by discount
	sort.Slice(filteredProducts, func(i, j int) bool {
		return filteredProducts[i].GetDiscount() > filteredProducts[j].GetDiscount()
	})

	whiteWines := strings.Builder{}
	redWines := strings.Builder{}
	otherProducts := strings.Builder{}

	whiteWines.WriteString("## Hvite viner\n\n| Navn (Vivino link) | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Land |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	redWines.WriteString("## Røde viner\n\n| Navn (Vivino link) | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Land |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	otherProducts.WriteString("## Andre produkter\n\n| Navn (Vivino link) | ID | Gammel pris | Ny pris | Delta | Rabatt % | Score | Type |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n")

	for _, product := range filteredProducts {
		switch product.Type {
		case "Hvitvin":
			whiteWines.WriteString("| " + product.GetVivinoMarkdownUrl() + " | " + product.GetVinmonopoletMarkdownUrl() + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.GetPriceDelta()) + " | " + strconv.Itoa(product.GetDiscount()) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Country + " |\n")
		case "Rødvin":
			redWines.WriteString("| " + product.GetVivinoMarkdownUrl() + " | " + product.GetVinmonopoletMarkdownUrl() + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.GetPriceDelta()) + " | " + strconv.Itoa(product.GetDiscount()) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Country + " |\n")
		default:
			otherProducts.WriteString("| " + product.GetVivinoMarkdownUrl() + " | " + product.GetVinmonopoletMarkdownUrl() + " | " + strconv.Itoa(product.VinmonopoletPrice) + " | " + strconv.Itoa(product.ApertifPrice) + " | " + strconv.Itoa(product.GetPriceDelta()) + " | " + strconv.Itoa(product.GetDiscount()) + " | " + strconv.Itoa(product.ApertifScore) + " | " + product.Type + " |\n")
		}
	}

	// Write the header
	readme.WriteString("# Go(d) wine\n\n")
	readme.WriteString("Liste over produkter som blir billigere etter månedsskiftet:\n\n")

	readme.WriteString(whiteWines.String())
	readme.WriteString(redWines.String())
	readme.WriteString(otherProducts.String())
}
