package vm

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for VM provisioning.
type API interface {
	NewDefinition(location, templateType, templateID, hostname string, cpus, memory, disk int, network []Network) Definition
	NewDefinitionWithDNS(location, templateType, templateID, hostname string, cpus, memory, disk int, network []Network, dnsServers []string) Definition
	Deprovision(ctx context.Context, identifier string, delayed bool) (DeprovisionResponse, error)
	Provision(ctx context.Context, definition Definition, base64Encoding bool) (ProvisioningResponse, error)
	Update(ctx context.Context, vmID string, change Change) (ProvisioningResponse, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new provisioning API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
