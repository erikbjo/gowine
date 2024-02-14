package utils

// Wine is a struct that represents a wine from Vinmonopolet
type Wine struct {
	Basic       Basic
	LastChanged LastChanged
}

// Basic is a struct that represents the basic information of a wine from Vinmonopolet
type Basic struct {
	ProductId        string `json:"productId"`
	ProductShortName string `json:"productShortName"`
}

// LastChanged is a struct that represents the last changed information of a wine from Vinmonopolet
type LastChanged struct {
	Date string `json:"date"`
	Time string `json:"time"`
}
