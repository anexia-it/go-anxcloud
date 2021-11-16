package acl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/frontend"
	"net/http"
	"net/url"
	utils "path"
	"strconv"
)

const path = "/api/LBaaS/v1/ACL.json"

type ACLInfo struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type ACL struct {
	CustomerIdentifier string                 `json:"customer_identifier"`
	ResellerIdentifier string                 `json:"reseller_identifier"`
	Identifier         string                 `json:"identifier"`
	Name               string                 `json:"name"`
	ParentType         string                 `json:"parent_type"`
	Frontend           *frontend.FrontendInfo `json:"frontend"`
	Backend            *backend.BackendInfo   `json:"backend"`
	Criterion          string                 `json:"criterion"`
	Index              int                    `json:"index"`
	Value              string                 `json:"value"`
}

func (a api) Get(ctx context.Context, page, limit int) ([]ACLInfo, error) {

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
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get load balancer ACLS %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []ACLInfo `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer ACL list response: %w", err)
	}

	return payload.Data.Data, nil

}

func (a api) GetByID(ctx context.Context, identifier string) (ACL, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return ACL{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return ACL{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return ACL{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return ACL{}, fmt.Errorf("could not execute get load balancer ACL request for '%s': %s", identifier,
			response.Status)
	}

	var payload ACL

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return ACL{}, fmt.Errorf("could not parse load balancer ACL response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a api) Create(ctx context.Context, definition Definition) (ACL, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return ACL{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return ACL{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return ACL{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return ACL{}, fmt.Errorf("error when creating ACL '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return ACL{}, fmt.Errorf("could not create load balancer ACL '%s': %s", definition.Name,
			response.Status)
	}

	var payload ACL

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return ACL{}, fmt.Errorf("could not parse load balancer ACL creation response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (ACL, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return ACL{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(path, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return ACL{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &requestBody)
	if err != nil {
		return ACL{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return ACL{}, fmt.Errorf("error when updating ACL '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return ACL{}, fmt.Errorf("could not update load balancer ACL '%s': %s", definition.Name,
			response.Status)
	}

	var payload ACL

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return ACL{}, fmt.Errorf("could not parse load balancer ACL updating response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a api) DeleteByID(ctx context.Context, identifier string) error {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(path, identifier)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when deleting a LBaaS ACL '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS ACL '%s': %s",
			identifier, response.Status)
	}
	return nil
}
