package vivino

import (
	"gowine/internal/shared"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
)

func ScrapeVivino(wine *shared.Product) {
	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	})

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error while visiting Vivino: %s\n", err)
	})

	// Scrape ratingz
	firstMatch := true
	c.OnHTML(".text-inline-block.light.average__number", func(e *colly.HTMLElement) {
		if firstMatch {
			ratingText := e.Text[:3]
			log.Printf("%s rating: %s", wine.Basic.ProductId, ratingText)
			firstMatch = false
			wine.VivinoScore = ratingText
		}
	})

	// TODO: Link to the wine

	// Visit the page
	url := wine.GetVivinoUrl()
	err := c.Visit(url)
	if err != nil {
		log.Printf("Error while executing colly: %s", err.Error())
	}
}
