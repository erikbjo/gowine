package vinmonopolet

type Product struct {
	Basic struct {
		ProductId        string `json:"productId"`
		ProductShortName string `json:"productShortName"`
	} `json:"basic"`
	LastChanged struct {
		Date string `json:"date"`
		Time string `json:"time"`
	} `json:"lastChanged"`
}

func (p *Product) GetVinmonopoletUrl() string {
	return "https://www.vinmonopolet.no/p/" + p.Basic.ProductId
}

type Products struct {
	Products []Product `json:"products"`
}
