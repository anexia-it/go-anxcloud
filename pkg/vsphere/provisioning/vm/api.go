package vm

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for VM provisioning.
type API interface {
	NewDefinition(location, templateType, templateID, hostname string, cpus, memory, disk int, network []Network) Definition
	Deprovision(ctx context.Context, identifier string, delayed bool) error
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
