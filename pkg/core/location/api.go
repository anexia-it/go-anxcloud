package location

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for retrieving and listing locations.
type API interface {
	pagination.Pageable
	List(ctx context.Context, page, limit int, search string) ([]Location, error)
	Get(ctx context.Context, identifier string) (Location, error)
	GetByCode(ctx context.Context, code string) (Location, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new VLAN API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
