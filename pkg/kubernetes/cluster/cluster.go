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
)

// The Cluster resource configures settings common for all specific backend Server resources linked to it.
type Cluster struct {
	CustomerIdentifier         string      `json:"customer_identifier"`
	ResellerIdentifier         string      `json:"reseller_identifier"`
	CriticalOperationPassword  interface{} `json:"critical_operation_password"` // TODO
	CriticalOperationConfirmed bool        `json:"critical_operation_confirmed"`
	Identifier                 string      `json:"identifier"`
	Name                       string      `json:"name"`
	State                      struct {
		Text  string `json:"text"`
		Title string `json:"title"`
		ID    string `json:"id"`
		Type  int    `json:"type"`
	} `json:"state"`
	Location struct {
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
	} `json:"location"`
	Version                  string      `json:"version"`
	PatchVersion             string      `json:"patch_version"`
	Kubeconfig               string      `json:"kubeconfig"`
	Autoscaling              bool        `json:"autoscaling"`
	CniPlugin                string      `json:"cni_plugin"`
	APIServerAllowlist       interface{} `json:"apiserver_allowlist"`       // TODO
	MaintenanceWindowStart   interface{} `json:"maintenance_window_start"`  // TODO
	MaintenanceWindowLength  interface{} `json:"maintenance_window_length"` // TODO
	ManageInternalIpv4Prefix bool        `json:"manage_internal_ipv4_prefix"`
	InternalIpv4Prefix       struct {
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
	} `json:"internal_ipv4_prefix"`
	NeedsServiceVMs          bool   `json:"needs_service_vms"`
	EnableNATGateways        bool   `json:"enable_nat_gateways"`
	EnableLBaaS              bool   `json:"enable_lbaas"`
	ExternalIPFamilies       string `json:"external_ip_families"`
	ManageExternalIPv4Prefix bool   `json:"manage_external_ipv4_prefix"`
	ExternalIPv4Prefix       struct {
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
	} `json:"external_ipv4_prefix"`
	ManageExternalIPv6Prefix bool `json:"manage_external_ipv6_prefix"`
	ExternalIPv6Prefix       struct {
		Name string `json:"name"`
	} `json:"external_ipv6_prefix"`
	AutomationRules []struct {
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
	} `json:"automation_rules"`
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
