// Package prefix implements API functions residing under /ipam/prefix.
// This path contains methods for querying and setting IP prefixes.
package prefix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	CreateEmpty             bool   `json:"create_empty"`
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

// Vlan reference definition
type Vlan struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	CustomerDescription string `json:"description_customer"`
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
	RouterRedundancy    bool       `json:"router_redundancy"`
	Vlans               []Vlan     `json:"vlans"`
	PrefixType          int        `json:"type"`
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
	Data []Summary `json:"data"`
}

func (a api) List(ctx context.Context, page, limit int, search string) ([]Summary, error) {
	escapedquerry := url.QueryEscape(search)
	url := fmt.Sprintf(
		"%s%s?page=%v&limit=%v&search=%s",
		a.client.BaseURL(),
		pathPrefix, page, limit, escapedquerry,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create prefix list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute prefix list request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute prefix list request, got response %s", httpResponse.Status)
	}

	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	if err != nil {
		return nil, fmt.Errorf("could not decode prefix list response: %w", err)
	}

	return responsePayload.Data, nil
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
		return Info{}, fmt.Errorf("could not create prefix get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("could not execute prefix get request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Info{}, fmt.Errorf("could not execute prefix get request, got response %s", httpResponse.Status)
	}

	var info Info
	err = json.NewDecoder(httpResponse.Body).Decode(&info)
	if err != nil {
		return Info{}, fmt.Errorf("could not decode prefix get response: %w", err)
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
		return fmt.Errorf("could not create prefix delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute prefix delete request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute prefix delete request, got response %s", httpResponse.Status)
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
		panic(fmt.Sprintf("could not create request data for prefix creation: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create prefix post request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute prefix post request: %w", err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Summary{}, fmt.Errorf("could not execute prefix post request, got response %s", httpResponse.Status)
	}

	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	if err != nil {
		return Summary{}, fmt.Errorf("could not decode prefix post response: %w", err)
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
		panic(fmt.Sprintf("could not create request data for prefix update: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create prefix update request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute prefix update request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Summary{}, fmt.Errorf("could not execute prefix update request, got response %s", httpResponse.Status)
	}

	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	if err != nil {
		return summary, fmt.Errorf("could not decode prefix update response: %w", err)
	}

	return summary, err
}
