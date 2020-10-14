// Package ips implements API functions residing under /provisioning/ips.
// This path contains methods for querying ips in a VLAN.
package ips

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anexia-it/go-anxcloud/pkg/client"
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
func GetFree(ctx context.Context, location, vlan string, c client.Client) ([]IP, error) {
	url := fmt.Sprintf(
		"https://%s%s/%s/%s",
		client.DefaultHost,
		pathPrefix,
		location,
		vlan,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create ips request: %w", err)
	}
	httpResponse, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get ips: %w", err)
	}

	responsePayload := response{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode ip get response: %w", err)
	}

	return responsePayload.Data, err
}
