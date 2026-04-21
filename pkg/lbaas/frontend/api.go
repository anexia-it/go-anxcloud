package frontend

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
	"go.anx.io/go-anxcloud/pkg/pagination"
)

// API contains methods for load balancer frontend management.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]genericresource.Identity, error)
	GetByID(ctx context.Context, identifier string) (Frontend, error)
	Create(ctx context.Context, definition Definition) (Frontend, error)
	Update(ctx context.Context, identifier string, definition Definition) (Frontend, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new frontend API instance with the given client.
func NewAPI(c client.Client) genericresource.API[Frontend, Definition] {
	return &api{c}
}
