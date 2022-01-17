package location

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for location querying.
type API interface {
	List(ctx context.Context, page, limit int, locationCode, organization string) ([]Location, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new location API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
