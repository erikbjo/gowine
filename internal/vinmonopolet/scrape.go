package vinmonopolet

import (
	"gowine/internal/shared"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func ScrapeVinmonopolet(wine *shared.Product, retryNumber int) {
	if retryNumber > 5 {
		// cba
		return
	}

	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext,
	})

	// Sjekker om utgått, utgått hvis product-price-expired finnes
	c.OnHTML(".product-price-expired", func(e *colly.HTMLElement) {
		wine.VinmonopoletPrice = -1
		// log.Printf("Product %s is expired", wine.Basic.ProductId)
	})

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
		}
	})

	// Scrape volume, but only take the first element
	c.OnHTML(".amount", func(e *colly.HTMLElement) {
		volumeText := e.Text
		re := regexp.MustCompile(`[^0-9]`)
		volume := re.ReplaceAllString(volumeText, "")
		wine.Volume, _ = strconv.Atoi(volume)
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
			}
		}
	})

	// Visit the page
	url := wine.GetVinmonopoletUrl()
	err := c.Visit(url)
	if err != nil {
		ScrapeVinmonopolet(wine, retryNumber+1)
	}
}
