package disk

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

// NodepoolDisks represents the disks of a [Nodepool].
type NodepoolDisks struct {
	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`
	Identifier         string `json:"identifier,omitempty"`
	Name               string `json:"name,omitempty"`

	SizeBytes       uint64              `json:"size_bytes,omitempty"`
	PerformanceType common.IDTitleTuple `json:"performance_type,omitempty"`
}

// NodepoolDisksDefinition represents the disks of a [Nodepool].
type NodepoolDisksDefinition struct {
	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`
	Name               string `json:"name,omitempty"`

	Identifier      string              `json:"identifier,omitempty"`
	SizeBytes       uint64              `json:"size_bytes,omitempty"`
	PerformanceType DiskPerformanceType `json:"performance_type,omitempty"`
}

type DiskPerformanceType string

const (
	DiskPerformanceTypeSTD1 DiskPerformanceType = "STD1"
	DiskPerformanceTypeSTD2 DiskPerformanceType = "STD2"
	DiskPerformanceTypeSTD3 DiskPerformanceType = "STD3"
	DiskPerformanceTypeSTD4 DiskPerformanceType = "STD4"
	DiskPerformanceTypeSTD5 DiskPerformanceType = "STD5"
	DiskPerformanceTypeENT1 DiskPerformanceType = "ENT1"
	DiskPerformanceTypeENT2 DiskPerformanceType = "ENT2"
	DiskPerformanceTypeENT3 DiskPerformanceType = "ENT3"
	DiskPerformanceTypeENT4 DiskPerformanceType = "ENT4"
	DiskPerformanceTypeENT5 DiskPerformanceType = "ENT5"
	DiskPerformanceTypeENT6 DiskPerformanceType = "ENT6"

	DiskPerformanceTypeDefault = DiskPerformanceTypeENT6
)

func (a *api) GetByID(ctx context.Context, identifier string) (NodepoolDisks, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return NodepoolDisks{},
			fmt.Errorf("could not execute get kubernetes nodepool disks request for '%s': %s",
				identifier,
				response.Status)
	}

	var payload NodepoolDisks

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return NodepoolDisks{},
			fmt.Errorf("could not parse kubernetes nodepool disks response for '%s' : %w",
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
		return nil, fmt.Errorf("could not get kubernetes nodepool disks %s", response.Status)
	}

	payload := struct {
		Data []common.PartialResource `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse kubernetes nodepool disks list response: %w",
			err)
	}

	return payload.Data, nil
}
func (a *api) Create(ctx context.Context, definition NodepoolDisksDefinition) (NodepoolDisks, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return NodepoolDisks{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("error when creating nodepool disk '%s': %w",
			definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return NodepoolDisks{},
			fmt.Errorf("could not create kubernetes nodepool disk '%s': %s",
				definition.Name,
				response.Status)
	}
	var payload NodepoolDisks

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return NodepoolDisks{},
			fmt.Errorf("could not parse kubernetes nodepool creation response for '%s' : %w",
				definition.Name, err)
	}

	return payload, nil
}
func (a *api) Update(ctx context.Context, identifier string,
	definition NodepoolDisksDefinition) (NodepoolDisks, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return NodepoolDisks{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint.String(), &requestBody)
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return NodepoolDisks{}, fmt.Errorf("error when updating nodepool disk '%s': %w",
			definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return NodepoolDisks{},
			fmt.Errorf("could not update kubernetes nodepool disk '%s': %s",
				definition.Name,
				response.Status)
	}

	var payload NodepoolDisks

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return NodepoolDisks{},
			fmt.Errorf("could not parse kubernetes nodepool disk updating response for"+
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
		return fmt.Errorf("error when deleting a kubernetes nodepool disk '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete kubernetes nodepool disk '%s': %s",
			identifier, response.Status)
	}
	return nil
}
