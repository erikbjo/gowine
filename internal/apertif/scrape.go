package apertif

import (
	"gowine/internal/shared"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

func ScrapeApertif(wine *shared.Product, retryNumber int) {
	if retryNumber > 5 {
		return
	}

	// TODO: new collector for each wine?
	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext,
	})

	// Scrape price
	// Apertif may return multiple results, so checking product id
	c.OnHTML("li.product-list-element", func(e *colly.HTMLElement) {
		attrID := e.Attr("data-product-id")
		if attrID != wine.Basic.ProductId {
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

	// TODO: remove apertif score, not used
	// Scrape score
	c.OnHTML(".number", func(e *colly.HTMLElement) {
		scoreText := e.Text
		score, err := strconv.Atoi(scoreText)
		if err != nil {
			// fmt.Printf("Failed to scrape score for %s", wine.Basic.ProductId)
		} else {
			wine.ApertifScore = score
		}
	})

	// Visit the page
	err := c.Visit(wine.GetApertifUrl())
	if err != nil {
		ScrapeApertif(wine, retryNumber+1)
	}
}
