package backend

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/pagination"
)

// API contains methods for load balancer backend management.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]BackendInfo, error)
	GetByID(ctx context.Context, identifier string) (Backend, error)
	Create(ctx context.Context, definition Definition) (Backend, error)
	Update(ctx context.Context, identifier string, definition Definition) (Backend, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
