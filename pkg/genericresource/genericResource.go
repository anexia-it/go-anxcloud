package genericresource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	path2 "path"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/pagination"
)

type IGenericResource interface {
	GetIdentifier() string
	GetName() string
}

type API[R any, D any] interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]Identity, error)
	GetByID(ctx context.Context, identifier string) (R, error)
	Create(ctx context.Context, definition D) (R, error)
	Update(ctx context.Context, identifier string, definition D) (R, error)
	DeleteByID(ctx context.Context, identifier string) error
	//GetPath() string
}

type Identity struct {
	IGenericResource
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func (g Identity) GetIdentifier() string {
	return g.Identifier
}

func (g Identity) GetName() string {
	return g.Name
}

func GetPagedGeneric(ctx context.Context, page int, limit int, client client.Client, name string, path string) ([]Identity, error) {
	endpoint, err := url.Parse(client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path2.Join(endpoint.Path, path)
	query := endpoint.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(limit))
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when executing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not get load balancer %s %s", name, response.Status)
	}

	payload := struct {
		Data struct {
			Data []Identity `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer %s list response: %w", name, err)
	}

	return payload.Data.Data, nil
}

func GenericGetByID[R any](ctx context.Context, identifier string, client client.Client, name string, apipath string) (*R, error) {
	endpoint, err := url.Parse(client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path2.Join(endpoint.Path, apipath, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute get load balancer %s request for '%s': %s", name, identifier,
			response.Status)
	}

	var payload R

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer %s response for '%s' : %w", name, identifier, err)
	}

	return &payload, nil
}

func GenericCreate[R any, D any](ctx context.Context, definition D, client client.Client, name string, apiPath string) (*R, error) {
	endpoint, err := url.Parse(client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path2.Join(endpoint.Path, apiPath)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when creating %s': %w", name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not create load balancer %s: %s", name,
			response.Status)
	}

	var payload R

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer %s creation response: %w", name, err)
	}

	return &payload, nil
}

func GenericUpdate[R any, D any](ctx context.Context, identifier string, definition D, client client.Client, name string, apiPath string) (*R, error) {
	endpoint, err := url.Parse(client.BaseURL())
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path2.Join(endpoint.Path, apiPath, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when updating %s: %w", name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return nil, fmt.Errorf("could not update load balancer %s: %s", name, response.Status)
	}

	var payload R

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse load balancer %s updating response: %w", name, err)
	}

	return &payload, nil
}

func GenericDelete(ctx context.Context, identifier string, client client.Client, name string, apiPath string) error {
	endpoint, err := url.Parse(client.BaseURL())
	if err != nil {
		return fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = path2.Join(endpoint.Path, apiPath, identifier)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request object: %w", err)
	}

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error when deleting a LBaaS %s '%s': %w", name,
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete LBaaS %s '%s': %s", name,
			identifier, response.Status)
	}
	return nil
}
