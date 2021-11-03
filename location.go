package gosquare

//Location is the struct for a Square Location
type Location struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CountryCode  string `json:"country"`
	BusinessName string `json:"business_name"`
}

//Locations defines a list of Location
type Locations []Location

type LocationList struct {
	Locations	Locations `json:"locations"`
}
