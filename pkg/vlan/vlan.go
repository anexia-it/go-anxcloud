// Package vlan implements API functions residing under /vlan.
// This path contains methods for querying, creating and deleting of VLANs.
package vlan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/vlan/v1/vlan.json"
)

// Summary describes some attributes of a VLAN.
type Summary struct {
	Identifier          string `json:"identifier"`
	Name                string `json:"name"`
	CustomerDescription string `json:"description_customer"`
}

// Location is the metadata of a single location.
type Location struct {
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"country"`
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	CityCode    string `json:"city_code"`
}

// Info describes all attributes of a VLAN.
type Info struct {
	Identifier          string `json:"identifier"`
	Name                string `json:"name"`
	CustomerDescription string `json:"description_customer"`
	InternalDescription string `json:"description_internal"`
	Role                string `json:"role_text"`
	Status              string `json:"status"`
	VMProvisioning      bool   `json:"vm_provisioning,omitempty"`
	Locations           []Location
}

// CreateDefinition contains information required to create a VLAN.
type CreateDefinition struct {
	Location            string `json:"location"`
	VMProvisioning      bool   `json:"vm_provisioning,omitempty"`
	CustomerDescription string `json:"description_customer,omitempty"`
}

// UpdateDefinition contains information required to update a VLAN.
type UpdateDefinition struct {
	CustomerDescription string `json:"description_customer,omitempty"`
}

type listResponse struct {
	Data struct {
		Data []Summary `json:"data"`
	} `json:"data"`
}

func (a api) List(ctx context.Context, page, limit int, search string) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s?page=%d&limit=%d&search=%s",
		a.client.BaseURL(),
		pathPrefix, page, limit, search,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create vlan list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute vlan list request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute vlan list request, got response %s", httpResponse.Status)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode vlan list response: %w", err)
	}

	return responsePayload.Data.Data, nil
}

func (a api) Get(ctx context.Context, identifier string) (Info, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Info{}, fmt.Errorf("could not create vlan get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("could not execute vlan get request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Info{}, fmt.Errorf("could not execute vlan get request, got response %s", httpResponse.Status)
	}

	var info Info
	err = json.NewDecoder(httpResponse.Body).Decode(&info)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Info{}, fmt.Errorf("could not decode vlan get response: %w", err)
	}

	return info, nil
}

func (a api) Create(ctx context.Context, createDefinition CreateDefinition) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(createDefinition); err != nil {
		panic(fmt.Sprintf("could not create request data for vlan creation: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create vlan post request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute vlan post request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Summary{}, fmt.Errorf("could not execute vlan post request, got response %s", httpResponse.Status)
	}

	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Summary{}, fmt.Errorf("could not decode vlan post response: %w", err)
	}

	return summary, nil
}

func (a api) Update(ctx context.Context, identifier string, updateDefinition UpdateDefinition) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix, identifier,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(updateDefinition); err != nil {
		panic(fmt.Sprintf("could not create request data for vlan update: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &requestData)
	if err != nil {
		return fmt.Errorf("could not create vlan update request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute vlan update request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute vlan update request, got response %s", httpResponse.Status)
	}

	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return fmt.Errorf("could not decode vlan update response: %w", err)
	}

	return nil
}

func (a api) Delete(ctx context.Context, identifier string) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix, identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create vlan delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute vlan delete request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute vlan delete request, got response %s", httpResponse.Status)
	}

	return httpResponse.Body.Close()
}
