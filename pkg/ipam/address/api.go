package address

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"
	"github.com/anexia-it/go-anxcloud/pkg/utils/param"
)

// API contains methods for IP manipulation.
type API interface {
	pagination.Pageable
	List(ctx context.Context, page, limit int, search string) ([]Summary, error)
	Get(ctx context.Context, id string) (Address, error)
	GetFiltered(ctx context.Context, page, limit int, filters ...param.Parameter) ([]Summary, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, create Create) (Summary, error)
	Update(ctx context.Context, id string, update Update) (Summary, error)
	ReserveRandom(ctx context.Context, reserve ReserveRandom) (ReserveRandomSummary, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
