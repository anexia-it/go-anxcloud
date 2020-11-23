package disktype

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiskType represents a disk type that may be used for VMs.
type DiskType struct {
	ID          string `json:"id"`
	StorageType string `json:"storage_type"`
	Bandwidth   int    `json:"bandwidth"`
	Latency     int    `json:"latency"`
}

const (
	pathPrefix = "/api/vsphere/v1/provisioning/disk_type.json"
)

func (a api) List(ctx context.Context, locationID string) ([]DiskType, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix, locationID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create disk type list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute disk type list request: %w", err)
	}
	var responsePayload []DiskType
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode disk type list response: %w", err)
	}

	return responsePayload, err
}
