package v1

// anxcloud:object

// Location describes a Anexia site where resources can be deployed.
type Location struct {
	Identifier  string  `json:"identifier" anxcloud:"identifier"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	CountryCode string  `json:"country"`
	CityCode    string  `json:"city_code"`
	Latitude    *string `json:"lat"`
	Longitude   *string `json:"lon"`
}
