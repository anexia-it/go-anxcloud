package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/apis/common"
)

// NodepoolNetworkDefinition represents the networks of a [Nodepool].
type NodepoolNetworkDefinition struct {
	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`
	Name               string `json:"name,omitempty"`

	BandwidthLimit string `json:"bandwidth_limit,omitempty"`
	VLANID         string `json:"vlan,omitempty"`
}

// NodepoolNetwork represents the networks of a [Nodepool].
type NodepoolNetwork struct {
	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`
	Identifier         string `json:"identifier,omitempty"`
	Name               string `json:"name,omitempty"`

	BandwidthLimit common.IDTitleTuple    `json:"bandwidth_limit,omitempty"`
	VLAN           common.PartialResource `json:"vlan,omitempty"`
}

func (a *api) GetByID(ctx context.Context, identifier string) (NodepoolNetwork, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("error when executing request for '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return NodepoolNetwork{},
			fmt.Errorf("could not execute get kubernetes nodepool network request for '%s"+
				"': %s",
				identifier,
				response.Status)
	}

	var payload NodepoolNetwork

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return NodepoolNetwork{},
			fmt.Errorf("could not parse kubernetes nodepool network response for '%s' : %w",
				identifier, err)
	}

	return payload, nil
}

func (a *api) Get(ctx context.Context, page, limit int) ([]common.PartialResource, error) {
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
		return nil, fmt.Errorf("could not get kubernetes nodepool network %s", response.Status)
	}

	payload := struct {
		Data []common.PartialResource `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse kubernetes nodepool network list response: %w",
			err)
	}

	return payload.Data, nil
}
func (a *api) Create(ctx context.Context, definition NodepoolNetworkDefinition) (
	NodepoolNetwork, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return NodepoolNetwork{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("error when creating nodepool network '%s': %w",
			definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return NodepoolNetwork{},
			fmt.Errorf("could not create kubernetes nodepool network '%s': %s",
				definition.Name,
				response.Status)
	}
	var payload NodepoolNetwork

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return NodepoolNetwork{},
			fmt.Errorf("could not parse kubernetes nodepool creation response for '%s' : %w",
				definition.Name, err)
	}

	return payload, nil
}
func (a *api) Update(ctx context.Context, identifier string,
	definition NodepoolNetworkDefinition) (NodepoolNetwork, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return NodepoolNetwork{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint.String(), &requestBody)
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return NodepoolNetwork{}, fmt.Errorf("error when updating nodepool network '%s': %w",
			definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return NodepoolNetwork{},
			fmt.Errorf("could not update kubernetes nodepool network '%s': %s",
				definition.Name,
				response.Status)
	}

	var payload NodepoolNetwork

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return NodepoolNetwork{},
			fmt.Errorf("could not parse kubernetes nodepool network updating response for"+
				" '%s' : %w",
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
		return fmt.Errorf("error when deleting a kubernetes nodepool network '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete kubernetes nodepool network '%s': %s",
			identifier, response.Status)
	}
	return nil
}
