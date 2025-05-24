package vinmonopolet

import (
	"encoding/json"
	"fmt"
	"gowine/internal/shared"
	"net/http"
	"os"
)

const (
	apiUrl = "https://apis.vinmonopolet.no/products/v0/details-normal"
	START  = "0"
)

// GetWines returns all wines from Vinmonopolet
func GetWines() ([]shared.Product, error) {
	client := &http.Client{}

	defer client.CloseIdleConnections()

	r, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating vinmonopolet api request: %s", err.Error())
	}

	apiKey := os.Getenv("VINMONOPOLETAPIKEY")
	if apiKey == "" {
		return nil, fmt.Errorf("VINMONOPOLETAPIKEY is net set")
	}

	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Ocp-Apim-Subscription-Key", apiKey)

	//r.URL.RawQuery = "start=" + START //+ "&maxResults=100000"
	//r.URL.RawQuery = "maxResults=10000"

	res, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error doing vinmonopolet request: %s", err.Error())
	}

	decoder := json.NewDecoder(res.Body)
	var mp []shared.Product

	err = decoder.Decode(&mp)
	if err != nil {
		return nil, fmt.Errorf("error decoding vinmonopolet response: %s", err.Error())
	}

	return mp, nil
}
