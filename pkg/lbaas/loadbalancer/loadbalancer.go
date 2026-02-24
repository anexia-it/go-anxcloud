package loadbalancer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"

	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
)

type (
	// RuleInfo holds the name and identifier of a rule.
	RuleInfo = v1.RuleInfo

	// Loadbalancer holds the information of a load balancer instance.
	Loadbalancer = v1.LoadBalancer
)

// LoadBalancerInfo holds the identifier and the name of a load balancer
type LoadBalancerInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}

const (
	path = "api/LBaaS/v1/loadbalancer.json"
)

func (a api) Get(ctx context.Context, page, limit int) ([]LoadBalancerInfo, error) {
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
		return nil, fmt.Errorf("could not get load balancers %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []LoadBalancerInfo `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a api) GetByID(ctx context.Context, identifier string) (Loadbalancer, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Loadbalancer{}, fmt.Errorf("could not execute get load balancer request for '%s': %s", identifier,
			response.Status)
	}

	var payload Loadbalancer

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not parse load balancer response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a api) Create(ctx context.Context, definition Definition) (Loadbalancer, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Loadbalancer{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &buf)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("error when creating a LBaaS Loadbalancer '%s': %w",
			definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Loadbalancer{}, fmt.Errorf("could not create LBaaS Loadbalancer '%s': %s",
			definition.Name, response.Status)
	}

	var payload Loadbalancer
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not parse Loadbalancer creation response: %w", err)
	}

	return payload, nil
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Loadbalancer, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Loadbalancer{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &buf)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("error when updating a LBaaS Loadbalancer '%s': %w",
			definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Loadbalancer{}, fmt.Errorf("could not update LBaaS Loadbalancer '%s': %s",
			definition.Name, response.Status)
	}

	var payload Loadbalancer
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Loadbalancer{}, fmt.Errorf("could not parse Loadbalancer updating response: %w", err)
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
		return fmt.Errorf("error when deleting a LBaaS Loadbalancer '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS Loadbalancer '%s': %s",
			identifier, response.Status)
	}

	return nil
}
