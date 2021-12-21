package tags

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for tag control.
type API interface {
	List(ctx context.Context, page, limit int, query, serviceIdentifier, organizationIdentifier, order string, sortAscending bool) ([]Summary, error)
	Get(ctx context.Context, identifier string) (Info, error)
	Create(ctx context.Context, create Create) (Summary, error)
	Delete(ctx context.Context, tagID, serviceID string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new tags API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
