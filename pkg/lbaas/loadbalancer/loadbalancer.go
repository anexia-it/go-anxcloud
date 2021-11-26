package loadbalancer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"
)

const (
	path = "api/LBaaS/v1/loadbalancer.json"
)

// LoadBalancerInfo holds the identifier and the name of a load balancer
type LoadBalancerInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}

// RuleInfo holds the name and identifier of a rule.
type RuleInfo struct {
	Identifier string `json:"identifier" anxcloud:"identifier"`
	Name       string `json:"name"`
}

// anxcloud:object:hooks=RequestBodyHook

// Loadbalancer holds the information of a load balancer instance.
type Loadbalancer struct {
	CustomerIdentifier string     `json:"customer_identifier"`
	ResellerIdentifier string     `json:"reseller_identifier"`
	Identifier         string     `json:"identifier" anxcloud:"identifier"`
	Name               string     `json:"name"`
	IpAddress          string     `json:"ip_address"`
	AutomationRules    []RuleInfo `json:"automation_rules"`
}

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
