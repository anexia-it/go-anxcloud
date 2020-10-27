package search

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for VM searching.
type API interface {
	ByName(ctx context.Context, name string) ([]VM, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new search API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
