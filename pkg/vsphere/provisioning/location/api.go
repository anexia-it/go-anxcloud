package location

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for location querying.
type API interface {
	All(ctx context.Context, page, limit int) ([]Location, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new location API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
