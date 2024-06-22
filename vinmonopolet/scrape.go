package vinmonopolet

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func ScrapeVinmonopolet(wine *Product) {
	c := colly.NewCollector()

	// Scrape price
	c.OnHTML(".product__price", func(e *colly.HTMLElement) {
		priceText := e.Text
		re := regexp.MustCompile(`[^0-9]`)
		price := re.ReplaceAllString(priceText, "")
		if len(price) > 2 {
			price = price[:len(price)-2] // Remove last 2 digits
		}
		wine.VinmonopoletPrice, _ = strconv.Atoi(price)
	})

	// Scrape volume
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
			} else {
				fmt.Printf("Failed to scrape alcohol for %s (%s)\n", wine.Basic.ProductShortName, wine.Basic.ProductId)
			}
		}
	})

	// Visit the page
	url := wine.GetVinmonopoletUrl()
	err := c.Visit(url)
	if err != nil {
		log.Println("Error in visiting page:", err.Error())
	}

	// c.Wait()
	// time.Sleep(time.Millisecond * 50)
}
