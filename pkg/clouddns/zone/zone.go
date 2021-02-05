// Package zone implements API functions residing under /zone.
// This path contains methods for querying and setting the DNS zones and records.
package zone

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	pathPrefix string = "api/clouddns/v1/zone.json"
)

type listResponse struct {
	Results []Zone `json:"results"`
}

type Record struct {
	Identifier string `json:"identifier"`
	Immutable bool `json:"immutable"`
	Name string `json:"name"`
	RData string `json:"rdata"`
	Region string `json:"region"`
	TTL string `json:"ttl"`
	Type string `json:"Type"`
}

type Revision struct {
	CreatedAt time.Time `json:"created_at"`
	Identifier string `json:"identifier"`
	ModifiedAt time.Time `json:"modified_at"`
	Records []Record `json:"records"`
}

type Zone struct {
	Name string `json:"name"`
	IsMaster bool `json:"master"`
	MasterNS string `json:"master_ns"`
	Customer string `json:"customer"`
	AdminEmail string `json:"admin_email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	NotifyAllowedIPs []string `json:"notify_allowed_ips"`
	IsEditable bool `json:"is_editable"`
	TTL int `json:"ttl"`
	ValidationLevel int `json:"validation_level"`
	Revisions []Revision `json:"revisions"`
}

// List Zones API method´
func (a api) List(ctx context.Context) ([]Zone, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create zone list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute zone list request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute zone list request, got response %s", httpResponse.Status)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode zone list respone: %w", err)
	}

	return responsePayload.Results, nil
}

// Get zone details API method
func (a api) Get(ctx context.Context, name string) (Zone, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Zone{}, fmt.Errorf("could not create zone get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Zone{}, fmt.Errorf("could not execute zone get request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Zone{}, fmt.Errorf("could not execute zone get request, got response %s", httpResponse.Status)
	}

	var responsePayload Zone
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Zone{}, fmt.Errorf("could not decode zone get respone: %w", err)
	}

	return responsePayload, nil
}

// create
// update zone
// delete zone

func (a api) Delete(ctx context.Context, name string) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create zone delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute zone delete request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute zone delete request, got response %s", httpResponse.Status)
	}

	return httpResponse.Body.Close()
}

// apply (changeset)
// import
