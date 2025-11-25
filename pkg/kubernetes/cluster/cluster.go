package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	utils "path"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

// The Cluster resource configures settings common for all specific backend Server resources linked to it.
type Cluster struct {
	gs.HasState

	CustomerIdentifier         string    `json:"customer_identifier"`
	ResellerIdentifier         string    `json:"reseller_identifier"`
	CriticalOperationPassword  string    `json:"critical_operation_password"`
	CriticalOperationConfirmed bool      `json:"critical_operation_confirmed"`
	Identifier                 string    `json:"identifier"`
	Name                       string    `json:"name"`
	Location                   Minimal   `json:"location"`
	Version                    string    `json:"version"`
	PatchVersion               string    `json:"patch_version"`
	Kubeconfig                 string    `json:"kubeconfig"`
	Autoscaling                bool      `json:"autoscaling"`
	CniPlugin                  string    `json:"cni_plugin"`
	APIServerAllowlist         string    `json:"apiserver_allowlist"`
	Backend                    string    `json:"backend"`
	BackendName                string    `json:"backend_name"`
	MaintenanceWindowStart     string    `json:"maintenance_window_start"`
	MaintenanceWindowLength    string    `json:"maintenance_window_length"`
	ManageInternalIPv4Prefix   bool      `json:"manage_internal_ipv4_prefix"`
	InternalIpv4Prefix         Minimal   `json:"internal_ipv4_prefix"`
	NeedsServiceVMs            bool      `json:"needs_service_vms"`
	EnableNATGateways          bool      `json:"enable_nat_gateways"`
	EnableLBaaS                bool      `json:"enable_lbaas"`
	ExternalIPFamilies         string    `json:"external_ip_families"`
	ManageExternalIPv4Prefix   bool      `json:"manage_external_ipv4_prefix"`
	ExternalIPv4Prefix         Minimal   `json:"external_ipv4_prefix"`
	ManageExternalIPv6Prefix   bool      `json:"manage_external_ipv6_prefix"`
	ExternalIPv6Prefix         Minimal   `json:"external_ipv6_prefix"`
	AutomationRules            []Minimal `json:"automation_rules"`

	ServiceVM01InternalIPv4Address Minimal `json:"service_vm_01_internal_ipv4_address"`
	ServiceVM02InternalIPv4Address Minimal `json:"service_vm_02_internal_ipv4_address"`
	ServiceVM01ExternalIPv4Address Minimal `json:"service_vm_01_external_ipv4_address"`
	ServiceVM02ExternalIPv4Address Minimal `json:"service_vm_02_external_ipv4_address"`
	ServiceVM01ExternalIPv6Address Minimal `json:"service_vm_01_external_ipv6_address"`
	ServiceVM02ExternalIPv6Address Minimal `json:"service_vm_02_external_ipv6_address"`
	ServiceLB01                    Minimal `json:"service_lb_01"`
	ServiceLB02                    Minimal `json:"service_lb_02"`
	ExternalIPv4VIP                Minimal `json:"external_ipv4_vip"`
	ExternalIPv6VIP                Minimal `json:"external_ipv6_vip"`
	KKPAPILBaaSBackend01           Minimal `json:"kkp_api_lbaas_backend_01"`
	KKPAPILBaaSBackend02           Minimal `json:"kkp_api_lbaas_backend_02"`
	KKPVPNLBaaSBackend01           Minimal `json:"kkp_vpn_lbaas_backend_01"`
	KKPVPNLBaaSBackend02           Minimal `json:"kkp_vpn_lbaas_backend_02"`
	StorageServerInterfaceAddress  Minimal `json:"storage_server_interface_address"`
}

type Minimal struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// ClusterInfo holds the identifier and the name of a kubernetes cluster.
type ClusterInfo struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func (a *api) Get(ctx context.Context, page, limit int) ([]ClusterInfo, error) {
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
		return nil, fmt.Errorf("could not get kubernetes clusters %s", response.Status)
	}

	payload := struct {
		Data struct {
			Data []ClusterInfo `json:"data"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse kubernetes cluster list response: %w", err)
	}

	return payload.Data.Data, nil
}

func (a *api) GetByID(ctx context.Context, identifier string) (Cluster, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Cluster{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Cluster{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Cluster{}, fmt.Errorf("error when executing request for '%s': %w", identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Cluster{}, fmt.Errorf("could not execute get kubernetes cluster request for '%s': %s", identifier,
			response.Status)
	}

	var payload Cluster

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Cluster{}, fmt.Errorf("could not parse kubernetes cluster response for '%s' : %w", identifier, err)
	}

	return payload, nil
}

func (a *api) Create(ctx context.Context, definition Definition) (Cluster, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Cluster{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Cluster{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &requestBody)
	if err != nil {
		return Cluster{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Cluster{}, fmt.Errorf("error when creating cluster '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Cluster{}, fmt.Errorf("could not create kubernetes cluster '%s': %s", definition.Name,
			response.Status)
	}

	var payload Cluster

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Cluster{}, fmt.Errorf("could not parse kubernetes cluster creation response for '%s' : %w",
			definition.Name, err)
	}

	return payload, nil
}

func (a *api) Update(ctx context.Context, identifier string, definition Definition) (Cluster, error) {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return Cluster{}, fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, identifier)

	requestBody := bytes.Buffer{}
	if err := json.NewEncoder(&requestBody).Encode(definition); err != nil {
		return Cluster{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint.String(), &requestBody)
	if err != nil {
		return Cluster{}, fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return Cluster{}, fmt.Errorf("error when updating cluster '%s': %w", definition.Name, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return Cluster{}, fmt.Errorf("could not update kubernetes cluster '%s': %s", definition.Name,
			response.Status)
	}

	var payload Cluster

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return Cluster{}, fmt.Errorf("could not parse kubernetes cluster updating response for '%s' : %w",
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
		return fmt.Errorf("error when deleting a kubernetes cluster '%s': %w",
			identifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not delete kubernetes cluster '%s': %s",
			identifier, response.Status)
	}
	return nil
}

func (a *api) RequestKubeConfig(ctx context.Context, cluster *Cluster) error {
	const name = "Request kubeconfig"
	var ruleID string

	for _, i := range cluster.AutomationRules {
		if i.Name == name {
			ruleID = i.Identifier
		}
	}

	return a.triggerAutomation(ctx, ruleID, cluster.Identifier)
}

func (a *api) triggerAutomation(ctx context.Context, ruleIdentifier, clusterIdentifier string) error {
	endpoint, err := url.Parse(a.client.BaseURL())
	if err != nil {
		return fmt.Errorf("could not parse URL: %w", err)
	}

	endpoint.Path = utils.Join(endpoint.Path, a.path, clusterIdentifier, "rule", ruleIdentifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request object: %w", err)
	}

	response, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("error when firing automation rule '%s': %w", ruleIdentifier, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 && response.StatusCode < 600 {
		return fmt.Errorf("could not fire automation rule '%s': %s", ruleIdentifier,
			response.Status)
	}

	return nil
}
