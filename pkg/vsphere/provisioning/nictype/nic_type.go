// Package nictype implements API functions residing under /provisioning/nictype.
// This path contains methods for querying existing nic types.
package nictype

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/vsphere/v1/provisioning/nic_type.json"
)

// List queries the API for nic types.
func (a api) List(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create nic types list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute nic types list request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute nic types list request, got response %s", httpResponse.Status)
	}

	var responsePayload []string
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode nic types list response: %w", err)
	}

	return responsePayload, err
}
