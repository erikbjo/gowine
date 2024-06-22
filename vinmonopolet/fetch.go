package vinmonopolet

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiUrl = "https://apis.vinmonopolet.no/products/v0/details-normal"
)

func init() {
	// This package is used to fetch wines from Vinmonopolet
}

// GetWines returns all wines from Vinmonopolet
func GetWines() []Product {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	defer client.CloseIdleConnections()

	r, err1 := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err1 != nil {
		log.Println("Error in creating request:", err1.Error())
		return nil
	}

	apiKey := os.Getenv("VINMONOPOLETAPIKEY")
	if apiKey == "" {
		log.Println("VINMONOPOLETAPIKEY is not set")
		return nil
	}

	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Ocp-Apim-Subscription-Key", apiKey)

	r.URL.RawQuery = "maxResults=100" + "&start=1000"

	res, err2 := client.Do(r)
	if err2 != nil {
		log.Println("Error in response:", err2.Error())
		return nil
	}

	decoder := json.NewDecoder(res.Body)
	var mp []Product

	err := decoder.Decode(&mp)
	if err != nil {
		log.Println("Error during decoding: " + err.Error())
		return nil
	}

	return mp
}
