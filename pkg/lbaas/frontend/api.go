package frontend

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for load balancer frontend management.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]FrontendInfo, error)
	GetByID(ctx context.Context, identifier string) (Frontend, error)
	Create(ctx context.Context, definition Definition) (Frontend, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new frontend API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
