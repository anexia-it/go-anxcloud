// Package search implements API functions residing under /search.
// This path contains methods to search VMs.
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	pathPrefix = "/api/vsphere/v1/search/by_name.json"
)

// VM is a single VM and its metatadata.
type VM struct {
	Name            string `json:"name"`
	Identifier      string `json:"identifier"`
	LocationCode    string `json:"location_code"`
	LocationCountry string `json:"location_country"`
	LocationName    string `json:"location_name"`
	PrimaryIPv4     string `json:"ip_v4_primary"`
	PrimaryIPv6     string `json:"ip_v6_primary"`
	OSName          string `json:"os_name"`
	OSFamily        string `json:"os_family"`
	Tags            string `json:"tags"`
}

type response struct {
	Data []VM `json:"data"`
}

// ByName returns VMs that matches the given name.
//
// ctx is attached to the request and will cancel it on cancelation.
// name is the name search string. It may contain wildcards as stated in the API docs.
func (a api) ByName(ctx context.Context, name string) ([]VM, error) {
	params := url.Values{}
	params.Add("name", name)
	url := fmt.Sprintf(
		"%s%s?%s",
		a.client.BaseURL(),
		pathPrefix,
		params.Encode(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create VM search request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute VM search request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute VM search request, got response %s", httpResponse.Status)
	}

	var responsePayload response
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return nil, fmt.Errorf("could not decode VM search response: %w", err)
	}

	return responsePayload.Data, err
}
