// Package address implements API functions residing under /ipam/address.
// This path contains methods for managing IPs.
package address

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/vsphere/v1/address.json"
)

// Address contains all the information about a specific address.
type Address struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	DescriptionInternal string `json:"description_internal"`
	Role                string `json:"role"`
	Version             int    `json:"version"`
	Status              string `json:"status"`
	VLANID              string `json:"vlan"`
	PrefixID            string `json:"prefix"`
}

// Summary is the address information returned by a listing.
type Summary struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	Role                string `json:"role"`
}

type allResponse struct {
	Data []Summary `json:"data"`
}

func (a api) All(ctx context.Context) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathPrefix,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create address list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute address list request: %w", err)
	}
	var responsePayload allResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode address list response: %w", err)
	}

	return responsePayload.Data, err
}

func (a api) Get(ctx context.Context, id string) (Address, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Address{}, fmt.Errorf("could not create address get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Address{}, fmt.Errorf("could not execute address get request: %w", err)
	}
	var responsePayload Address
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return Address{}, fmt.Errorf("could not decode address get response: %w", err)
	}

	return responsePayload, err
}

func (a api) Delete(ctx context.Context, id string) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create address delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute address delete request: %w", err)
	}

	return httpResponse.Body.Close()
}
