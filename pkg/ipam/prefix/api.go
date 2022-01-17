package prefix

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for IP manipulation.
type API interface {
	List(ctx context.Context, page, limit int) ([]Summary, error)
	Get(ctx context.Context, id string) (Info, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, create Create) (Summary, error)
	Update(ctx context.Context, id string, update Update) (Summary, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
