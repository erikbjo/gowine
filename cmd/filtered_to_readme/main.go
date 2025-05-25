package main

import (
	"encoding/json"
	"gowine/internal/shared"
	"os"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

var logger = shared.CreateSugaredLogger()

func main() {
	// Take the filtered products from the scrape process and write them to the README
	// in a nice format

	// Load gowine products from JSON
	file, err := os.Open("json/gowine_products.json")
	if err != nil {
		logger.Fatal("Failed to open file", zap.Error(err))
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Warn("Failed to close gowine_products file", zap.Error(err))
		}
	}()

	// Read the products
	var filteredProducts []*shared.Product
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&filteredProducts)
	if err != nil {
		logger.Fatal("Failed to decode gowine products", zap.Error(err))
	}

	// Prepare the README
	readme, err := os.Create("README.md")
	if err != nil {
		logger.Fatal("Failed to create README", zap.Error(err))
	}
	defer func() {
		err = readme.Close()
		if err != nil {
			logger.Warn("Failed to close readme file", zap.Error(err))
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

	// helper function
	write := func(description string, content string) {
		if _, err := readme.WriteString(content); err != nil {
			logger.Fatal("write failed", zap.String("part", description), zap.Error(err))
		}
	}

	write("header", "# Go(d) wine\n\n")
	write("intro", "Liste over produkter som blir billigere etter månedsskiftet:\n\n")
	write("white wines", whiteWines.String())
	write("red wines", redWines.String())
	write("other products", otherProducts.String())
}
