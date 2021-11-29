package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"

	v1 "github.com/anexia-it/go-anxcloud/pkg/apis/lbaas/v1"
)

const (
	path = "api/LBaaS/v1/backend.json"
)

// The Backend resource configures settings common for all specific backend Server resources linked to it.
type Backend = v1.Backend

// BackendInfo holds the identifier and the name of a load balancer backend.
type BackendInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}

func (a api) Get(ctx context.Context, page, limit int) ([]BackendInfo, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path)
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
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get load balancer backends %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []BackendInfo `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer backend list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a api) GetByID(ctx context.Context, identifier string) (Backend, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Backend{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Backend{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Backend{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Backend{}, fmt.Errorf("could not execute get load balancer backend request for '%s': %s", identifier,
			response.Status)
	}

	var payload Backend

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Backend{}, fmt.Errorf("could not parse load balancer backend response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a api) Create(ctx context.Context, definition Definition) (Backend, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Backend{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Backend{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return Backend{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Backend{}, fmt.Errorf("error when creating backend '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Backend{}, fmt.Errorf("could not create load balancer backend '%s': %s", definition.Name,
			response.Status)
	}

	var payload Backend

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Backend{}, fmt.Errorf("could not parse load balancer backend creation response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Backend, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Backend{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Backend{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &requestBody)
	if err != nil {
		return Backend{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Backend{}, fmt.Errorf("error when updating backend '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Backend{}, fmt.Errorf("could not update load balancer backend '%s': %s", definition.Name,
			response.Status)
	}

	var payload Backend

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Backend{}, fmt.Errorf("could not parse load balancer backend updating response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when deleting a LBaaS backend '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS backend '%s': %s",
			identifier, response.Status)
	}
	return nil
}
