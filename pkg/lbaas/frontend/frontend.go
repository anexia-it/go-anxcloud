package frontend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/common"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/loadbalancer"
	"net/http"
	"net/url"
	utils "path"
	"strconv"
)

const path = "api/LBaaS/v1/frontend.json"

// FrontendInfo holds the name and the identifier of a frontend
type FrontendInfo struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// Frontend represents a LBaaS Frontend.
type Frontend struct {
	CustomerIdentifier string                         `json:"customer_identifier"`
	ResellerIdentifier string                         `json:"reseller_identifier"`
	Identifier         string                         `json:"identifier"`
	Name               string                         `json:"name"`
	LoadBalancer       *loadbalancer.LoadBalancerInfo `json:"load_balancer,omitempty"`
	DefaultBackend     *backend.BackendInfo           `json:"default_backend,omitempty"`
	Mode               common.Mode                    `json:"mode"`
	ClientTimeout      string                         `json:"client_timeout"`
}

func (a api) Get(ctx context.Context, page, limit int) ([]FrontendInfo, error) {
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
		return nil, fmt.Errorf("could not get load balancer frontends %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []FrontendInfo `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer frontend list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a api) GetByID(ctx context.Context, identifier string) (Frontend, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Frontend{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Frontend{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Frontend{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Frontend{}, fmt.Errorf("could not execute get load balancer frontend request for '%s': %s", identifier,
			response.Status)
	}

	var payload Frontend

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Frontend{}, fmt.Errorf("could not parse load balancer frontend response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a api) Create(ctx context.Context, definition Definition) (Frontend, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Frontend{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Frontend{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &buf)
	if err != nil {
		return Frontend{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Frontend{}, fmt.Errorf("error when creating a LBaaS frontend for load balancer '%s': %w",
			definition.LoadBalancer, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Frontend{}, fmt.Errorf("could not create LBaaS frontend for load balancer '%s': %s",
			definition.LoadBalancer, response.Status)
	}

	var payload Frontend
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Frontend{}, fmt.Errorf("could not parse loadbalancer frontend creation response: %w", err)
	}

	return payload, nil
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Frontend, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Frontend{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Frontend{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &buf)
	if err != nil {
		return Frontend{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Frontend{}, fmt.Errorf("error when updating a LBaaS frontend for load balancer '%s': %w",
			definition.LoadBalancer, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Frontend{}, fmt.Errorf("could not update LBaaS frontend for load balancer '%s': %s",
			definition.LoadBalancer, response.Status)
	}

	var payload Frontend
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Frontend{}, fmt.Errorf("could not parse loadbalancer frontend updating response: %w", err)
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
		return fmt.Errorf("error when deleting a LBaaS frontend '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS frontend '%s': %s",
			identifier, response.Status)
	}

	return nil
}
