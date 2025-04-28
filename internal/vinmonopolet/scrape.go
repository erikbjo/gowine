package vinmonopolet

import (
	"fmt"
	"github.com/gocolly/colly"
	"gowine/internal/shared"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ScrapeVinmonopolet(wine *shared.Product) {
	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	})

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error while visiting Vinmonopolet: %s\n", err)
	})

	// Sjekker om utgått, utgått hvis product-price-expired finnes
	c.OnHTML(".product-price-expired", func(e *colly.HTMLElement) {
		wine.VinmonopoletPrice = -1
		// log.Printf("Product %s is expired", wine.Basic.ProductId)
	})

	// Boolean flags to scrape only the first price and volume
	// priceScraped := false
	// volumeScraped := false

	// Scrape price, but only take the first element
	c.OnHTML(".product__price", func(e *colly.HTMLElement) {
		if wine.VinmonopoletPrice == 0 {
			priceText := e.Text
			re := regexp.MustCompile(`[^0-9]`)
			price := re.ReplaceAllString(priceText, "")
			if len(price) > 2 {
				price = price[:len(price)-2] // Remove last 2 digits
			}
			wine.VinmonopoletPrice, _ = strconv.Atoi(price)
			// log.Printf("Scraped price for %s: %d", wine.Basic.ProductId, wine.VinmonopoletPrice)
			// priceScraped = true // Set the flag to true after scraping the first price
		}
	})

	// Scrape volume, but only take the first element
	c.OnHTML(".amount", func(e *colly.HTMLElement) {
		volumeText := e.Text
		re := regexp.MustCompile(`[^0-9]`)
		volume := re.ReplaceAllString(volumeText, "")
		wine.Volume, _ = strconv.Atoi(volume)
		// log.Printf("Scraped volume for %s: %d", wine.Basic.ProductId, wine.Volume)
		// volumeScraped = true // Set the flag to true after scraping the first volume
	})

	// Scrape type
	c.OnHTML(".product__category-name", func(e *colly.HTMLElement) {
		wine.Type = e.Text
	})

	// Scrape country
	c.OnHTML(".product__region", func(e *colly.HTMLElement) {
		wine.Country = strings.Split(e.Text, ",")[0]
	})

	// Scrape grape
	c.OnHTML(".label-xWJ3XCYJ", func(e *colly.HTMLElement) {
		wine.Grape = e.Text
	})

	// Scrape alcohol
	c.OnHTML(".content-item-wLPXgMvT", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Alkohol") {
			alcoholText := e.ChildText("span")
			re := regexp.MustCompile(`(\d+,\d+|\d+)%`)
			match := re.FindStringSubmatch(alcoholText)
			if match != nil {
				alcohol := strings.Replace(match[1], ",", ".", 1)
				wine.Alcohol, _ = strconv.ParseFloat(alcohol, 64)
			} else {
				fmt.Printf("Failed to scrape alcohol for %s\n", wine.Basic.ProductId)
			}
		}
	})

	// Visit the page
	url := wine.GetVinmonopoletUrl()
	err := c.Visit(url)
	if err != nil {
		log.Println("Error in visiting Vinmonopolet:", err.Error())
	}
}
