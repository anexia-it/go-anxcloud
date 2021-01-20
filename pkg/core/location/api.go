package location

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for location listing.
type API interface {
	List(ctx context.Context, page, limit int, search string) ([]Location, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new VLAN API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
