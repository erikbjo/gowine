package apertif

import (
	"fmt"
	"github.com/gocolly/colly"
	"gowine/internal/shared"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func ScrapeApertif(wine *shared.Product) {
	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	})

	// Error handling
	// c.OnError(func(r *colly.Response, err error) {
	// 	log.Printf("\tError while visiting Apertif: %s", err)
	// })

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
			fmt.Printf("Failed to scrape score for %s", wine.Basic.ProductId)
		} else {
			wine.ApertifScore = score
		}
	})

	// Visit the page
	url := wine.GetApertifUrl()
	err := c.Visit(url)
	if err != nil {
		log.Println("Error while visiting Apertif: " + err.Error())
	}

	// time.Sleep(time.Millisecond * 2000)
}