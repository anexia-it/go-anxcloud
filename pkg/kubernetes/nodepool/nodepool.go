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

	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

type SyncSource string

const (
	SyncSourceEngine  SyncSource = "engine"
	SyncSourceCluster SyncSource = "cluster"
)

type GSBase struct {
	CustomerIdentifier string `json:"customer_identifier"`
	ResellerIdentifier string `json:"reseller_identifier"`
	Identifier         string `json:"identifier"`
	Name               string `json:"name"`
}

// The Nodepool resource represents the main resource to map to the MachineDeployment in the customer cluster.
type Nodepool struct {
	gs.HasState
	GSBase

	CriticalOperationPassword  string `json:"critical_operation_password"`
	CriticalOperationConfirmed bool   `json:"critical_operation_confirmed"`

	Cluster            common.PartialResource `json:"cluster"`
	SyncSource         SyncSource             `json:"syncsource"`
	Replicas           uint                   `json:"replicas"`
	CPUs               uint                   `json:"cpus"`
	CPUType            string                 `json:"cputype"`
	MemoryBytes        uint64                 `json:"memory"`
	DiskSizeBytes      uint64                 `json:"disk_size"`
	OperatingSystem    string                 `json:"operating_system"`
	AutoscalerEnabled  bool                   `json:"autoscaler_enabled"`
	AutoscalerMinNodes bool                   `json:"autoscaler_min_nodes"`
	AutoscalerMaxNodes bool                   `json:"autoscaler_max_nodes"`

	Disks    []NodepoolDisks `json:"disks"`
	Networks []NodepoolDisks `json:"networks"`

	CustomDNSEnabled bool   `json:"customdns_enabled"`
	DNSOverrideIPv4  bool   `json:"dns_override_ipv4"`
	DNSv4Entry1      string `json:"dns_v4_1"`
	DNSv4Entry2      string `json:"dns_v4_2"`

	DNSOverrideIPv6 bool   `json:"dns_override_ipv6"`
	DNSv6Entry1     string `json:"dns_v6_1"`
	DNSv6Entry2     string `json:"dns_v6_2"`

	Taints      string `json:"taints"`
	Labels      string `json:"labels"`
	Annotations string `json:"annotations"`
	SSHPubKeys  string `json:"sshpubkeys"`

	AutomationRules []common.PartialResource `json:"automation_rules"`
}

type NodepoolDisks struct {
	GSBase
	SizeBytes       uint64 `json:"size_bytes"`
	PerformanceType string `json:"performance_type"`
}

type NodepoolNetwork struct {
	GSBase
	BandwidthLimit uint                   `json:"bandwidth_limit"`
	VLAN           common.PartialResource `json:"vlan"`
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
		return nil, fmt.Errorf("could not get kubernetes nodepools %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []common.PartialResource `json:"data"`
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
