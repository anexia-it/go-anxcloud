package address

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for IP manipulation.
type API interface {
	List(ctx context.Context, page, limit int, search string) ([]Summary, error)
	Get(ctx context.Context, id string) (Address, error)
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
