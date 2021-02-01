// Package prefix implements API functions residing under /ipam/prefix.
// This path contains methods for querying and setting IP prefixes.
package prefix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Type defines the type of a prefix.

const (
	pathPrefix string = "/api/ipam/v1/prefix.json"
	// TypePublic means the prefix is globally routable.
	TypePublic int = 0
	// TypePrivate means the prefix is scoped to a private network.
	TypePrivate int = 1
)

// Create defines meta data of a prefix to create.
type Create struct {
	Location    string `json:"location"`
	IPVersion   int    `json:"version"`
	Type        int    `json:"type"`
	NetworkMask int    `json:"netmask"`

	CreateVLAN              bool   `json:"new_vlan,omitempty"`
	VLANID                  string `json:"vlan,omitempty"`
	EnableRedundancy        bool   `json:"router_redundancy,omitempty"`
	EnableVMProvisioning    bool   `json:"vm_provisioning,omitempty"`
	CustomerDescription     string `json:"description_customer,omitempty"`
	CustomerVLANDescription string `json:"description_vlan_customer,omitempty"`
	Organization            string `json:"organization,omitempty"`
}

// NewCreate creates a new prefix definition with required vlaues.
func NewCreate(location, vlan string, ipVersion int, prefixType int, networkMask int) Create {
	return Create{
		Location:    location,
		IPVersion:   ipVersion,
		Type:        prefixType,
		NetworkMask: networkMask,
		VLANID:      vlan,
	}
}

// Location is part of info.
type Location struct {
	ID        string `json:"identifier"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	CityCode  string `json:"city_code"`
}

// Info contains extended information about a prefix.
type Info struct {
	ID                  string     `json:"identifier"`
	Name                string     `json:"name"`
	CustomerDescription string     `json:"description_customer"`
	InternalDescription string     `json:"description_internal"`
	IPVersion           int        `json:"version"`
	NetworkMask         int        `json:"netmask"`
	Role                string     `json:"role_text"`
	Status              string     `json:"status"`
	Locations           []Location `json:"locations"`
	VLANID              string     `json:"vlan_id"`
	RouterRedundancy    bool       `json:"router_redundancy"`
}

// Update contains fields to change on a prefix.
type Update struct {
	Name                string `json:"name,omitempty"`
	CustomerDescription string `json:"description_customer,omitempty"`
}

// Summary contains a abbreviated information set about a prefix.
type Summary struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	CustomerDescription string `json:"description_customer"`
}

type listResponse struct {
	Data struct {
		Data []Summary `json:"data"`
	}
}

func (a api) List(ctx context.Context, page, limit int) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s?page=%v&limit=%v",
		a.client.BaseURL(),
		pathPrefix, page, limit,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create vlan list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute vlan list request: %w", err)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode vlan list response: %w", err)
	}

	return responsePayload.Data.Data, nil
}

func (a api) Get(ctx context.Context, id string) (Info, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Info{}, fmt.Errorf("could not create vlan get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("could not execute vlan get request: %w", err)
	}
	var info Info
	err = json.NewDecoder(httpResponse.Body).Decode(&info)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Info{}, fmt.Errorf("could not decode vlan get response: %w", err)
	}

	return info, nil
}

func (a api) Delete(ctx context.Context, id string) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix, id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create vlan delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute vlan delete request: %w", err)
	}

	return httpResponse.Body.Close()
}

func (a api) Create(ctx context.Context, create Create) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(create); err != nil {
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
	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Summary{}, fmt.Errorf("could not decode vlan post response: %w", err)
	}

	return summary, nil
}

func (a api) Update(ctx context.Context, id string, update Update) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix, id,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(update); err != nil {
		panic(fmt.Sprintf("could not create request data for vlan update: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create vlan update request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute vlan update request: %w", err)
	}
	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return summary, fmt.Errorf("could not decode vlan update response: %w", err)
	}

	return summary, err
}
