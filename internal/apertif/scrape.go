package apertif

import (
	"gowine/internal/shared"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

var logger = shared.CreateSugaredLogger()
var baseCollector *colly.Collector

func init() {
	baseCollector = colly.NewCollector()
	baseCollector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	err := baseCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*aperitif.no*",
		Parallelism: 7,
		Delay:       2 * time.Second,
		RandomDelay: 3 * time.Second,
	})
	if err != nil {
		logger.Error("Error setting colly limits", zap.Error(err))
	}

	baseCollector.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext,
	})
}

func ScrapeApertif(wine *shared.Product, retryNumber int) {
	if retryNumber > 5 {
		logger.Warn("Skipping wine due to 5+ retries", zap.String("ProductId", wine.Basic.ProductId))
		return
	}

	c := baseCollector.Clone()

	// Scrape price
	// Apertif may return multiple results, so checking product id
	c.OnHTML("li.product-list-element", func(e *colly.HTMLElement) {
		indexText := e.ChildText(".index") // Returns e.g. "(13601102)"
		if !strings.Contains(indexText, wine.Basic.ProductId) {
			return
		}

		priceText := e.ChildText(".price")
		if priceText == "" {
			return
		}

		priceText = strings.ReplaceAll(priceText, "Pris: ", "") // "67.30 kr"
		priceText = strings.ReplaceAll(priceText, " kr", "")    // "67.30"
		priceText = strings.ReplaceAll(priceText, " ", "")      // Handle thousands separators if any
		parts := strings.Split(priceText, ".")                  // Split by decimal

		if p, err := strconv.Atoi(parts[0]); err == nil {
			wine.ApertifPrice = p
		}
	})

	// Visit the page
	err := c.Visit(wine.GetApertifUrl())
	if err != nil {
		ScrapeApertif(wine, retryNumber+1)
	}
}
