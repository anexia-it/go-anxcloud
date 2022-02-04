package v1

// anxcloud:object:hooks=ResponseDecodeHook

// Location describe a Anexia site where resources can be deployed.
type Location struct {
	Identifier  string   `json:"identifier" anxcloud:"identifier"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	CountryCode string   `json:"country"`
	CityCode    string   `json:"city_code"`
	Latitude    *float64 `json:"lat"`
	Longitude   *float64 `json:"lon"`
}
