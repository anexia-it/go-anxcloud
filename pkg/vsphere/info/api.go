package info

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for VM info querying.
type API interface {
	Get(ctx context.Context, identifier string) (Info, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new info API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
