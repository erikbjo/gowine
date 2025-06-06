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

	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext,
	})

	// Scrape price
	c.OnHTML(".price", func(e *colly.HTMLElement) {
		priceText := e.Text
		re := regexp.MustCompile(`[^0-9]`)
		price := re.ReplaceAllString(priceText, "")
		if len(price) > 2 {
			price = price[:len(price)-2] // Remove last 2 digits
		}
		wine.ApertifPrice, _ = strconv.Atoi(price)
	})

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
	url := wine.GetApertifUrl()
	err := c.Visit(url)
	if err != nil {
		ScrapeApertif(wine, retryNumber+1)
	}
}
