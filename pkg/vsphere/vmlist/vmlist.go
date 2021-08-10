package vmlist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	path = "api/vsphere/v1/vmlist/list.json"
)

type VM struct {
	Name            string `json:"name"`
	CustomName      string `json:"custom_name"`
	Identifier      string `json:"identifier"`
	LocationCode    string `json:"location_code"`
	LocationCountry string `json:"location_country"`
	LocationName    string `json:"location_name"`
	OSName          string `json:"os_name"`
	OSFamily        string `json:"os_family"`
	Tags            string `json:"tags"`
}

func (a api) Get(ctx context.Context, page, limit int) ([]VM, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path
	query := endpoint.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(limit))
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when executing request: %w", err)
	}

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute VM list, got response %s", response.Status)
	}

	payload := struct {
		Data []VM `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse vm list response: %w", err)
	}

	return payload.Data, nil
}
