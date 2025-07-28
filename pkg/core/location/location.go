// Package location implements API functions residing under /core/location.json.
// This path contains methods for querying existing locations.
package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	pathPrefix = "/api/core/v1/location"
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
	escapedquerry := url.QueryEscape(search)
	url := fmt.Sprintf(
		"%s%s.json?page=%d&limit=%d&search=%s",
		a.client.BaseURL(),
		pathPrefix, page, limit, escapedquerry,
	)

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

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	if err != nil {
		return nil, fmt.Errorf("could not decode location list response: %w", err)
	}

	return responsePayload.Data.Data, nil
}

func (a api) Get(ctx context.Context, identifier string) (Location, error) {
	url := fmt.Sprintf(
		"%s%s.json/%s",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Location{}, fmt.Errorf("could not create location get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Location{}, fmt.Errorf("could not execute location get request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Location{}, fmt.Errorf("could not execute location get request, got response %s", httpResponse.Status)
	}

	var location Location
	err = json.NewDecoder(httpResponse.Body).Decode(&location)
	if err != nil {
		return Location{}, fmt.Errorf("could not decode location get response: %w", err)
	}

	return location, nil
}

func (a api) GetByCode(ctx context.Context, code string) (Location, error) {
	url := fmt.Sprintf(
		"%s%s/by-code.json/%s",
		a.client.BaseURL(),
		pathPrefix,
		code,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Location{}, fmt.Errorf("could not create location get-by-code request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Location{}, fmt.Errorf("could not execute location get-by-code request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Location{}, fmt.Errorf("could not execute location get-by-code request, got response %s", httpResponse.Status)
	}

	var location Location
	err = json.NewDecoder(httpResponse.Body).Decode(&location)
	if err != nil {
		return Location{}, fmt.Errorf("could not decode location get-by-code response: %w", err)
	}

	return location, nil
}
