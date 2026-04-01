package availabilityzones

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.anx.io/go-anxcloud/pkg/client"
)

type AvailabilityZone struct {
	Identifier     string   `json:"identifier"`
	Name           string   `json:"name"`
	CpuCategories  []string `json:"cpu_categories"`
	DiskCategories []string `json:"disk_categories"`
}

// API contains methods for location querying.
type API interface {
	List(ctx context.Context, locationID string) ([]AvailabilityZone, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new location API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}

func (a api) List(ctx context.Context, locationID string) (
	zones []AvailabilityZone, err error) {

	url := fmt.Sprintf(
		"%s/api/vsphere/v1/provisioning/location.json/%s/availability_zone",
		a.client.BaseURL(),
		locationID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("could not create ListAvailabilityZones request: %w", err)
		return
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		err = fmt.Errorf("could not execute ListAvailabilityZones request: %w", err)
		return
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != 200 {
		err = fmt.Errorf("could not execute ListAvailabilityZones request, got response %s", httpResponse.Status)
		return
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&zones)
	if err != nil {
		err = fmt.Errorf("could not decode ListAvailabilityZones response: %w", err)
		return
	}

	return
}
