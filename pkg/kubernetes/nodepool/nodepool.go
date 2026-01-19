package nodepool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// The Nodepool resource configures settings common for all specific backend Server resources linked to it.
type Nodepool struct {
	gs.HasState

	CustomerIdentifier         string `json:"customer_identifier"`
	ResellerIdentifier         string `json:"reseller_identifier"`
	CriticalOperationPassword  string `json:"critical_operation_password"`
	CriticalOperationConfirmed bool   `json:"critical_operation_confirmed"`
	Identifier                 string `json:"identifier"`
	Name                       string `json:"name"`

	Cluster         Minimal `json:"cluster"`
	Replicas        int     `json:"replicas"`
	CPUs            int     `json:"cpus"`
	MemoryBytes     int     `json:"memory"`
	DiskSizeBytes   int     `json:"disk_size"`
	OperatingSystem string  `json:"operating_system"`

	AutomationRules []Minimal `json:"automation_rules"`
}

type Minimal struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// NodePoolInfo holds the identifier and the name of a kubernetes nodepool.
type NodePoolInfo struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func (a *api) Get(ctx context.Context, page, limit int) ([]NodePoolInfo, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path)
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
		return nil, fmt.Errorf("could not get kubernetes nodepools %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []NodePoolInfo `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse kubernetes nodepool list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a *api) GetByID(ctx context.Context, identifier string) (Nodepool, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Nodepool{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Nodepool{}, fmt.Errorf("could not execute get kubernetes nodepool request for '%s': %s", identifier,
			response.Status)
	}

	var payload Nodepool

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not parse kubernetes nodepool response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a *api) Create(ctx context.Context, definition Definition) (Nodepool, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Nodepool{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Nodepool{}, fmt.Errorf("error when creating nodepool '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Nodepool{}, fmt.Errorf("could not create kubernetes nodepool '%s': %s", definition.Name,
			response.Status)
	}

	var payload Nodepool

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not parse kubernetes nodepool creation response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a *api) Update(ctx context.Context, identifier string, definition Definition) (Nodepool, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Nodepool{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &requestBody)
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Nodepool{}, fmt.Errorf("error when updating nodepool '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Nodepool{}, fmt.Errorf("could not update kubernetes nodepool '%s': %s", definition.Name,
			response.Status)
	}

	var payload Nodepool

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Nodepool{}, fmt.Errorf("could not parse kubernetes nodepool updating response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a *api) DeleteByID(ctx context.Context, identifier string) error {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when deleting a kubernetes nodepool '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete kubernetes nodepool '%s': %s",
			identifier, response.Status)
	}
	return nil
}
