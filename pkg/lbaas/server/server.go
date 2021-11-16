package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"

	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
)

const (
	path = "/api/LBaaS/v1/server.json"
)

// ServerInfo holds the identifier and the name of a load balancer backend.
type ServerInfo struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// Server holds the information of a load balancers backend server
type Server struct {
	CustomerIdentifier string              `json:"customer_identifier"`
	ResellerIdentifier string              `json:"reseller_identifier"`
	Identifier         string              `json:"identifier"`
	Name               string              `json:"name"`
	IP                 string              `json:"ip"`
	Port               int                 `json:"port"`
	Backend            backend.BackendInfo `json:"backend"`
	Check              string              `json:"check"`
}

func (a api) Get(ctx context.Context, page, limit int) ([]ServerInfo, error) {
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

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get load balancer backend servers %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []ServerInfo `json:"data"`
		}
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer backend server list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a api) GetByID(ctx context.Context, identifier string) (Server, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Server{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Server{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Server{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Server{}, fmt.Errorf("could not execute get load balancer backend server request for '%s': %s", identifier,
			response.Status)
	}

	var payload Server

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Server{}, fmt.Errorf("could not parse load balancer backend server response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a api) Create(ctx context.Context, definition Definition) (Server, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Server{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Server{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &buf)
	if err != nil {
		return Server{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Server{}, fmt.Errorf("error when creating a LBaaS server for backend '%s': %w",
			definition.Backend, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Server{}, fmt.Errorf("could not create LBaaS server for backend '%s': %s",
			definition.Backend, response.Status)
	}

	var payload Server
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Server{}, fmt.Errorf("could not parse loadbalancer server creation response: %w", err)
	}

	return payload, nil
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Server, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Server{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Server{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &buf)
	if err != nil {
		return Server{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Server{}, fmt.Errorf("error when updating a LBaaS server for backend '%s': %w",
			definition.Backend, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Server{}, fmt.Errorf("could not update LBaaS server for backend '%s': %s",
			definition.Backend, response.Status)
	}

	var payload Server
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Server{}, fmt.Errorf("could not parse loadbalancer server updating response: %w", err)
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
		return fmt.Errorf("error when deleting a LBaaS server '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS server '%s': %s",
			identifier, response.Status)
	}

	return nil
}
