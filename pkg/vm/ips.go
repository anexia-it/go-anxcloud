package vm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

const (
	ipPathPrefix = "/api/vsphere/v1/provisioning/ips.json"
)

// IP defines informationen corresponding to the IP of a VLAN.
type IP struct {
	Identifier string `json:"identifier"`
	Text       string `json:"text"`
	Prefix     string `json:"prefix"`
}

// IPResponse is the response from the API regarding IP queries.
type IPResponse struct {
	Data []IP `json:"data"`
}

// GetFreeIPs returns the freen IPs on a VLAN.
func GetFreeIPs(ctx context.Context, location, vlan string, c client.Client) ([]IP, error) {
	url := fmt.Sprintf(
		"https://%s%s/%s/%s",
		client.DefaultHost,
		ipPathPrefix,
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

	responsePayload := IPResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode ip get response: %w", err)
	}

	return responsePayload.Data, err
}
