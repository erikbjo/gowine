package server

import (
	"encoding/json"
	"gowine/utils"
	"log"
	"net/http"
	"time"
)

// WineHandler
/*
Handle requests for /wine
*/
func WineHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleWineGetRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' are supported.", http.StatusNotImplemented)
		return
	}

}

/*
Handle GET request for /wine
*/
func handleWineGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	changedSince := r.URL.Query().Get("changedSince")
	if changedSince == "" {
		changedSince = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}

	// Get all wines from Vinmonopolet
	wines := getWines(w, changedSince)

	// Marshal wines to JSON
	marshaledWines, err := json.MarshalIndent(wines, "", "    ")
	if err != nil {
		log.Println("Error during JSON encoding: " + err.Error())
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}

	// Write wines to client
	_, err = w.Write(marshaledWines)
	if err != nil {
		log.Println("Failed to write response: " + err.Error())
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}

	log.Println("Wines returned to client.")
}

/*
Get all wines from vinmonopolet
*/
func getWines(w http.ResponseWriter, changedSince string) []utils.Wine {
	defer client.CloseIdleConnections()

	r, err1 := http.NewRequest(http.MethodGet, utils.VinMonopoletAPIURL, nil)
	if err1 != nil {
		log.Println("Error in creating request:", err1.Error())
		http.Error(w, "Error in creating request", http.StatusInternalServerError)
	}

	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Ocp-Apim-Subscription-Key", utils.VinMonopoletAPIKey)

	r.URL.RawQuery = "changedSince=" + changedSince

	res, err2 := client.Do(r)
	if err2 != nil {
		log.Println("Error in response:", err2.Error())
		http.Error(w, "Error in response", http.StatusInternalServerError)
	}

	mp := decodeJSON(w, res)

	return mp
}

/*
Decode JSON and return as map
*/
func decodeJSON(w http.ResponseWriter, res *http.Response) []utils.Wine {
	decoder := json.NewDecoder(res.Body)
	var mp []utils.Wine

	err := decoder.Decode(&mp)
	if err != nil {
		log.Println("Error during decoding: " + err.Error())
		http.Error(w, "Error during decoding", http.StatusBadRequest)
	}

	return mp
}
