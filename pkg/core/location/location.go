// Package location implements API functions residing under /core/location.json.
// This path contains methods for querying existing locations.
package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/core/v1/location.json"
)

// Location is the metadata of a single location.
type Location struct {
	Code        string `json:"code"`
	CityCode    string `json:"city_code"`
	Country     string `json:"country"`
	ID          string `json:"identifier"`
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	Name        string `json:"name"`
	CountryName string `json:"country_name"`
}

type listResponse struct {
	Data struct {
		Data []Location `json:"data"`
	} `json:"data"`
}

func (a api) List(ctx context.Context, page, limit int, search string) ([]Location, error) {
	url := fmt.Sprintf(
		"%s%s?page=%d&limit=%d&search=%s",
		a.client.BaseURL(),
		pathPrefix, page, limit, search,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create location list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute location list request: %w", err)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode location list response: %w", err)
	}

	return responsePayload.Data.Data, nil
}
