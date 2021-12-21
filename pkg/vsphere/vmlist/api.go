package vmlist

import (
	"context"
	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for VM listing.
type API interface {
	Get(ctx context.Context, page, limit int) ([]VM, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new vmlist API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
