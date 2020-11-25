// Package service implements API functions residing under /cofe/service.
// This path contains methods for querying services.
package service

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for listing services.
type API interface {
	List(ctx context.Context, page, limit int) ([]Service, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new service API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
