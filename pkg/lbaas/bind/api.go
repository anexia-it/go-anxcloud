package bind

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/genericResource"
	"go.anx.io/go-anxcloud/pkg/pagination"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for frontend bind management.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]genericResource.Identity, error)
	GetByID(ctx context.Context, identifier string) (Bind, error)
	Create(ctx context.Context, definition Definition) (Bind, error)
	Update(ctx context.Context, identifier string, definition Definition) (Bind, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new bind API instance with the given client.
func NewAPI(c client.Client) genericResource.API[Bind, Definition] {
	return &api{c}
}
