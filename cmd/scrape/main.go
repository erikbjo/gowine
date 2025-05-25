package main

import (
	"encoding/json"
	"gowine/internal/apertif"
	"gowine/internal/shared"
	"gowine/internal/vinmonopolet"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

var logger = shared.CreateSugaredLogger()

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}
}

func main() {
	err := cleanStaleFiles()
	if err != nil {
		logger.Error("Failed to clean stale files", zap.Error(err))
	}

	products, err := vinmonopolet.GetWines()
	if len(products) == 0 {
		logger.Fatal("Got zero wines from vinmonopolet")
	}
	if err != nil {
		logger.Fatal("Failed to get vinmonopolet wines", zap.Error(err))
	}

	var wg sync.WaitGroup
	var scrapedProducts []*shared.Product
	var validMutex sync.Mutex

	var expiredProducts []*shared.Product
	var expiredMutex sync.Mutex

	// Load expired products from JSON, if it exists
	file, err := os.Open("json/expired_products.json")
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&expiredProducts)
		if err != nil {
			logger.Fatal("Failed to decode expired products", zap.Error(err))
		}
		err = file.Close()
		if err != nil {
			logger.Warn("Failed to close file",
				zap.String("file", file.Name()),
				zap.Error(err))
		}
	}

	// Log expired products
	if len(expiredProducts) > 0 {
		logger.Info("Filtering expired products",
			zap.Int("amount", len(expiredProducts)))
		for _, product := range expiredProducts {
			for i, p := range products {
				if p.Basic.ProductId == product.Basic.ProductId {
					products = slices.Delete(products, i, i+1)
					break
				}
			}
		}
	}

	logger.Info("Starting to scrape products", zap.Int("amount", len(products)))

	loadingBar := progressbar.Default(int64(len(products)))

	// Limit the number of concurrent goroutines
	// Vinmonopolet are quick to throttle IPs from personal networks, use VPN etc with many goroutines
	semaphore := make(chan struct{}, 100)

	for _, product := range products {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(wine shared.Product) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// TODO: move functions to structs
			vinmonopolet.ScrapeVinmonopolet(&wine, 0)
			apertif.ScrapeApertif(&wine, 0)

			// Too many requests to vivino, TODO: fix, delay or something
			// vivino.ScrapeVivino(&wine)

			if wine.VinmonopoletPrice == -1 {
				expiredMutex.Lock()
				expiredProducts = append(expiredProducts, &wine)
				expiredMutex.Unlock()
			} else {
				validMutex.Lock()
				scrapedProducts = append(scrapedProducts, &wine)
				validMutex.Unlock()
			}
			_ = loadingBar.Add(1)
		}(product)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info("All scraping done, processing results.")

	// Filter products with complete pricing
	filteredProducts := filterCompleteProducts(scrapedProducts)
	priceDifferenceProducts := filterPriceDifferences(filteredProducts)

	// Save results to JSON
	err = saveToJSON("json/scraped_products.json", priceDifferenceProducts)
	if err != nil {
		logger.Error("Failed to save scraped products to json", zap.Error(err))
	}

	// Save expired products to JSON
	err = saveToJSON("json/expired_products.json", expiredProducts)
	if err != nil {
		logger.Error("Failed to save expired products to json", zap.Error(err))
	}

	logger.Info("Saved products to files",
		zap.Int("amountScraped", len(priceDifferenceProducts)),
		zap.Int("amountExpired", len(expiredProducts)))
}

// Filters products that have valid prices from both sources
func filterCompleteProducts(products []*shared.Product) []*shared.Product {
	var filtered []*shared.Product
	for _, product := range products {
		if product.VinmonopoletPrice != 0 && product.ApertifPrice != 0 {
			filtered = append(filtered, product)
		} else {
			logger.Info("Product har missing prices, skipping",
				zap.String("productId", product.Basic.ProductId),
				zap.String("productName", product.Basic.ProductShortName))
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
func saveToJSON(filename string, products []*shared.Product) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(products); err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

// cleanStaleFiles moves old files (scraped and gowine products) to the log dir for the current month
func cleanStaleFiles() error {
	staleFiles := []string{"json/scraped_products.json", "json/gowine_products.json"}

	for _, fileString := range staleFiles {
		var products []*shared.Product

		file, err := os.Open(fileString)
		if err == nil {
			decoder := json.NewDecoder(file)
			err = decoder.Decode(&products)
			if err != nil {
				return err
			}

			if len(products) > 0 {
				logger.Warn("Found stale scraped products, moving to archive", zap.String("file", fileString))

				split := strings.Split(file.Name(), "/")
				dirPath := split[0] + "/log/" + strings.ToLower(time.Now().UTC().Format("Jan")) + "/"

				err = os.MkdirAll(dirPath, os.ModePerm)
				if err != nil {
					return err
				}

				err = os.Rename(file.Name(), dirPath+split[1])
				if err != nil {
					return err
				}
			}

			err = file.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
