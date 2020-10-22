// Package location implements API functions residing under /provisioning/location.
// This path contains methods for querying existing locations.
package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/vsphere/v1/provisioning/location.json"
)

// Location is the metadata of a single location.
type Location struct {
	Code        string `json:"code"`
	Country     string `json:"country"`
	ID          string `json:"id"`
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	Name        string `json:"name"`
	CountryName string `json:"country_name"`
}

type response struct {
	Data []Location `json:"data"`
}

// All queries the API for known location.
//
// ctx is attached to the request and will cancel it on cancelation.
// definition contains the definition of the VM to be created.
func (a api) All(ctx context.Context) ([]Location, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create location list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute location list request: %w", err)
	}
	var responsePayload response
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode location list response: %w", err)
	}

	return responsePayload.Data, err
}
