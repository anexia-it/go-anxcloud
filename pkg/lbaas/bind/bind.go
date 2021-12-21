package bind

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

const (
	path = "api/LBaaS/v1/bind.json"
)

// BindInfo holds the identifier and the name of a load balancer frontend bind.
type BindInfo struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type Bind = v1.Bind

func (a api) Get(ctx context.Context, page, limit int) ([]BindInfo, error) {
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
		return nil, fmt.Errorf("could not get frontend binds %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []BindInfo `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse frontend binds list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a api) GetByID(ctx context.Context, identifier string) (Bind, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Bind{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Bind{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Bind{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Bind{}, fmt.Errorf("could not execute get frontend binds request for '%s': %s", identifier,
			response.Status)
	}

	var payload Bind

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Bind{}, fmt.Errorf("could not parse frontend binds response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a api) Create(ctx context.Context, definition Definition) (Bind, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Bind{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Bind{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &buf)
	if err != nil {
		return Bind{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Bind{}, fmt.Errorf("error when creating a frontend bind for frontend '%s': %w",
			definition.Frontend, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Bind{}, fmt.Errorf("could not create frontend bind for frontend '%s': %s",
			definition.Frontend, response.Status)
	}

	var payload Bind
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Bind{}, fmt.Errorf("could not parse frontend bind creation response: %w", err)
	}

	return payload, nil
}

func (a api) Update(ctx context.Context, identifier string, definition Definition) (Bind, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Bind{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, path, identifier)

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(definition); err != nil {
		return Bind{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &buf)
	if err != nil {
		return Bind{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Bind{}, fmt.Errorf("error when updating a frontend bind for frontend '%s': %w",
			definition.Frontend, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Bind{}, fmt.Errorf("could not update frontend bind for frontend '%s': %s",
			definition.Frontend, response.Status)
	}

	var payload Bind
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Bind{}, fmt.Errorf("could not parse frontend bind updating response: %w", err)
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
		return fmt.Errorf("error when deleting a LBaaS frontend bind '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS frontend bind '%s': %s",
			identifier, response.Status)
	}

	return nil
}
