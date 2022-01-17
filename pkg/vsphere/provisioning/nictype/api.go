package nictype

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for NIC type querying.
type API interface {
	List(ctx context.Context) ([]string, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new NIC type API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
