package rule

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
	"go.anx.io/go-anxcloud/pkg/genericResource"
)

const path = "/api/LBaaS/v1/rule.json"

type Rule = v1.Rule

func (a api) Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error) {

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
		return nil, fmt.Errorf("could not get load balancer Rules %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []genericResource.Identity `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer Rule list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a api) GetByID(ctx context.Context, identifier string) (Rule, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Rule{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Rule{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Rule{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Rule{}, fmt.Errorf("could not execute get load balancer Rule request for '%s': %s", identifier,
			response.Status)
	}

	var payload Rule

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Rule{}, fmt.Errorf("could not parse load balancer Rule response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a api) Create(ctx context.Context, definition Definition) (Rule, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Rule{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Rule{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return Rule{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Rule{}, fmt.Errorf("error when creating Rule '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Rule{}, fmt.Errorf("could not create load balancer Rule '%s': %s", definition.Name,
			response.Status)
	}

	var payload Rule

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Rule{}, fmt.Errorf("could not parse load balancer Rule creation response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Rule, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Rule{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Rule{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &requestBody)
	if err != nil {
		return Rule{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Rule{}, fmt.Errorf("error when updating Rule '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Rule{}, fmt.Errorf("could not update load balancer Rule '%s': %s", definition.Name,
			response.Status)
	}

	var payload Rule

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Rule{}, fmt.Errorf("could not parse load balancer Rule updating response for '%s' : %w",
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
		return fmt.Errorf("error when deleting a LBaaS Rule '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS Rule '%s': %s",
			identifier, response.Status)
	}
	return nil
}
