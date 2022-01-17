package vlan

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for VLAN control.
type API interface {
	List(ctx context.Context, page, limit int, search string) ([]Summary, error)
	Get(ctx context.Context, identifier string) (Info, error)
	Create(ctx context.Context, createDefinition CreateDefinition) (Summary, error)
	Delete(ctx context.Context, identifier string) error
	Update(ctx context.Context, identifier string, updateDefinition UpdateDefinition) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new VLAN API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
