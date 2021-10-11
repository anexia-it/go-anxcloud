// Package ips implements API functions residing under /provisioning/ips.
// This path contains methods for querying ips in a VLAN.
package ips

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/vsphere/v1/provisioning/ips.json"
)

// IP defines informationen corresponding to the IP of a VLAN.
type IP struct {
	Identifier string `json:"identifier"`
	Text       string `json:"text"`
	Prefix     string `json:"prefix"`
}

type response struct {
	Data []IP `json:"data"`
}

// GetFree returns information about the free IPs on a VLAN.
func (a api) GetFree(ctx context.Context, location, vlan string) ([]IP, error) {
	url := fmt.Sprintf(
		"%s%s/%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		location,
		vlan,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create ips request: %w", err)
	}
	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get ips: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute ip get request, got response %s", httpResponse.Status)
	}

	responsePayload := response{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)

	if err != nil {
		return nil, fmt.Errorf("could not decode ip get response: %w", err)
	}

	return responsePayload.Data, err
}
