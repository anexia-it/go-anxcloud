// Package search implements API functions residing under /search.
// This path contains methods to search VMs.
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anexia-it/go-anxcloud/pkg/client"
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
// client is the HTTP to be used for the request.
func ByName(ctx context.Context, name string, c client.Client) ([]VM, error) {
	url := fmt.Sprintf(
		"https://%s%s/%s",
		client.DefaultHost,
		pathPrefix,
		name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create VM search request: %w", err)
	}

	httpResponse, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute VM search request: %w", err)
	}
	var responsePayload response
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode VM search response: %w", err)
	}

	return responsePayload.Data, err
}
