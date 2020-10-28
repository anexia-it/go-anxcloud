package address

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for IP manipulation.
type API interface {
	All(ctx context.Context) ([]Summary, error)
	Get(ctx context.Context, id string) (Address, error)
	Delete(ctx context.Context, id string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
