package shared

type Product struct {
	Basic struct {
		ProductId        string `json:"productId"`
		ProductShortName string `json:"productShortName"`
	} `json:"basic"`
	LastChanged struct {
		Date string `json:"date"`
		Time string `json:"time"`
	} `json:"lastChanged"`
	Volume            int     `json:"volume"`
	VinmonopoletPrice int     `json:"vinmonopolet_price"`
	ApertifPrice      int     `json:"apertif_price"`
	ApertifScore      int     `json:"apertif_score"`
	Type              string  `json:"type"`
	Country           string  `json:"country"`
	Grape             string  `json:"grape"`
	Alcohol           float64 `json:"alcohol"`
	Difference        int     `json:"difference"`
	Discount          int     `json:"discount"`
}

func (p *Product) GetVinmonopoletUrl() string {
	return "https://www.vinmonopolet.no/p/" + p.Basic.ProductId
}

func (p *Product) GetApertifUrl() string {
	return "https://www.aperitif.no/pollisten?query=" + p.Basic.ProductId
}
