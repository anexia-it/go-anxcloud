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
func (a api) List(ctx context.Context, page, limit int, locationCode, organization string) ([]Location, error) {

	url := fmt.Sprintf(
		"%s%s?page=%v&limit=%v",
		a.client.BaseURL(),
		pathPrefix, page, limit,
	)

	//TODO remove this after ANEXIA API has been fixed - as it looks the endpoint does not support empty strings
	if len(locationCode) > 0 {
		url = fmt.Sprintf("%s&%s=%s", url, "location_code", locationCode)
	}

	//TODO remove this after ANEXIA API has been fixed - as it looks the endpoint does not support empty strings
	if len(organization) > 0 {
		url = fmt.Sprintf("%s&%s=%s", url, "organization", organization)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create location list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute location list request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute location list request, got response %s", httpResponse.Status)
	}

	var responsePayload response
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return nil, fmt.Errorf("could not decode location list response: %w", err)
	}

	return responsePayload.Data, err
}
