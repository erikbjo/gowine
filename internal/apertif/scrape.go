package apertif

import (
	"gowine/internal/shared"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

var logger = shared.CreateSugaredLogger()

func ScrapeApertif(wine *shared.Product, retryNumber int) {
	if retryNumber > 5 {
		return
	}

	// TODO: new collector for each wine?
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*aperitif.no*",
		Parallelism: 2,
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})
	if err != nil {
		logger.Error("Error setting colly limits", zap.Error(err))
	}

	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext,
	})

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

		re := regexp.MustCompile(`[^0-9]`)
		price := re.ReplaceAllString(priceText, "")
		if len(price) > 2 {
			price = price[:len(price)-2]
		}

		if p, err := strconv.Atoi(price); err == nil {
			wine.ApertifPrice = p
		}
	})

	// Visit the page
	err = c.Visit(wine.GetApertifUrl())
	if err != nil {
		ScrapeApertif(wine, retryNumber+1)
	}
}
