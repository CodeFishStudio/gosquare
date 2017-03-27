package gosquare

//Location is the struct for a Square Location
type Location struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"emai"`
	CountryCode  string `json:"country_code"`
	BusinessName string `json:"business_name"`
}

//Locations defines a list of Location
type Locations []Location
