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
	"go.anx.io/go-anxcloud/pkg/kubernetes/disk"
	"go.anx.io/go-anxcloud/pkg/kubernetes/network"
)

type SyncSource string

const (
	SyncSourceEngine  SyncSource = "engine"
	SyncSourceCluster SyncSource = "cluster"
)

const (
	OSFlatcar = "Flatcar Linux"
	MebiByte  = 1024 * 1024
	GibiByte  = MebiByte * 1024
)

var (
	StateOK    = gs.State{ID: "0", Text: "Deployed", Type: gs.StateTypeOK}
	StateNoGA  = gs.State{ID: "1", Text: "noGA", Type: gs.StateTypeOK}
	StateError = gs.State{ID: "2", Text: "Error", Type: gs.StateTypeError}
)

type CPUPerformanceType string

const (
	CPUPerformanceTypeBestEffort      CPUPerformanceType = "best-effort"
	CPUPerformanceTypeStandard        CPUPerformanceType = "standard"
	CPUPerformanceTypeEnterprise      CPUPerformanceType = "enterprise"
	CPUPerformanceTypePerformance     CPUPerformanceType = "performance"
	CPUPerformanceTypePerformancePlus CPUPerformanceType = "performance-plus"

	CPUPerformanceTypeDefault = CPUPerformanceTypePerformance
)

// The Nodepool resource represents the main resource to map to the MachineDeployment in the customer cluster.
type Nodepool struct {
	gs.HasState

	State              gs.State `json:"state,omitempty"`
	CustomerIdentifier string   `json:"customer_identifier"`
	ResellerIdentifier string   `json:"reseller_identifier"`
	Identifier         string   `json:"identifier"`
	Name               string   `json:"name"`

	Cluster            common.PartialResource `json:"cluster"`
	SyncSource         common.IDTitleTuple    `json:"syncsource"`
	Replicas           uint                   `json:"replicas"`
	CPUs               uint                   `json:"cpus"`
	CPUType            common.IDTitleTuple    `json:"cpu_performance_type"`
	MemoryBytes        uint64                 `json:"memory"`
	OperatingSystem    common.IDTitleTuple    `json:"operating_system"`
	AutoscalerEnabled  bool                   `json:"autoscaler_enabled"`
	AutoscalerMinNodes uint                   `json:"autoscaler_min_nodes"`
	AutoscalerMaxNodes uint                   `json:"autoscaler_max_nodes"`

	DiskSize            uint64                    `json:"disk_size"`
	DiskPerformanceType common.IDTitleTuple       `json:"disk_performance_type"`
	AdditionalDisks     []disk.NodepoolDisks      `json:"additional_disks"`
	Networks            []network.NodepoolNetwork `json:"networks"`

	DNSOverrideIPv4 bool   `json:"dns_override_ipv4"`
	DNSv4Entry1     string `json:"dns_v4_1"`
	DNSv4Entry2     string `json:"dns_v4_2"`

	DNSOverrideIPv6 bool   `json:"dns_override_ipv6"`
	DNSv6Entry1     string `json:"dns_v6_1"`
	DNSv6Entry2     string `json:"dns_v6_2"`

	Taints      string `json:"taints"`
	Labels      string `json:"labels"`
	Annotations string `json:"annotations"`
	SSHPubKeys  string `json:"sshpubkeys"`

	AutomationRules []common.PartialResource `json:"automation_rules"`
}

// The Definition resource represents the main resource to map to the MachineDeployment in the customer cluster.
type Definition struct {
	CustomerIdentifier string `json:"customer_identifier,omitempty"`
	ResellerIdentifier string `json:"reseller_identifier,omitempty"`
	Name               string `json:"name,omitempty"`

	ClusterID          string             `json:"cluster,omitempty"`
	SyncSource         SyncSource         `json:"syncsource,omitempty"`
	Replicas           uint               `json:"replicas,omitempty"`
	CPUs               uint               `json:"cpus,omitempty"`
	CPUType            CPUPerformanceType `json:"cpu_performance_type,omitempty"`
	MemoryBytes        uint64             `json:"memory,omitempty"`
	OperatingSystem    string             `json:"operating_system,omitempty"`
	AutoscalerEnabled  bool               `json:"autoscaler_enabled,omitempty"`
	AutoscalerMinNodes uint               `json:"autoscaler_min_nodes,omitempty"`
	AutoscalerMaxNodes uint               `json:"autoscaler_max_nodes,omitempty"`

	DiskSize            uint64                              `json:"disk_size"`
	DiskPerformanceType disk.DiskPerformanceType            `json:"disk_performance_type"`
	AdditionalDisks     []disk.NodepoolDisksDefinition      `json:"additional_disks,omitempty"`
	Networks            []network.NodepoolNetworkDefinition `json:"networks,omitempty"`

	DNSOverrideIPv4 bool   `json:"dns_override_ipv4,omitempty"`
	DNSv4Entry1     string `json:"dns_v4_1,omitemptyy"`
	DNSv4Entry2     string `json:"dns_v4_2,omitemptyy"`

	DNSOverrideIPv6 bool   `json:"dns_override_ipv6,omitempty"`
	DNSv6Entry1     string `json:"dns_v6_1,omitemptyy"`
	DNSv6Entry2     string `json:"dns_v6_2,omitempty"`

	Taints      string `json:"taints,omitemptyy"`
	Labels      string `json:"labels,omitemptyy"`
	Annotations string `json:"annotations,omitemptyy"`
	SSHPubKeys  string `json:"sshpubkeys,omitemptyy"`
}

func NewIDTitleTuple(id, title string) common.IDTitleTuple {
	return common.IDTitleTuple{
		ID:    id,
		Title: title,
	}
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
		Data []common.PartialResource `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse kubernetes nodepool list response: %w", err)
	}

	return payload.Data, nil
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint.String(), &requestBody)
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
