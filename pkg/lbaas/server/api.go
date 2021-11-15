package server

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for load balancer backend server management.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]ServerInfo, error)
	GetByID(ctx context.Context, identifier string) (Server, error)
	Create(ctx context.Context, definition Definition) (Server, error)
	Update(ctx context.Context, identifier string, definition Definition) (Server, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend server API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
