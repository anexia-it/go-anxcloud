package rule

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for load balancer backend management.
type API interface {
	Get(ctx context.Context, page, limit int) ([]RuleInfo, error)
	GetByID(ctx context.Context, identifier string) (Rule, error)
	Create(ctx context.Context, definition Definition) (Rule, error)
	Update(ctx context.Context, identifier string, definition Definition) (Rule, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
