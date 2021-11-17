package backend

import (
	"context"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"
)

// API contains methods for load balancer backend management.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]BackendInfo, error)
	GetByID(ctx context.Context, identifier string) (Backend, error)
	Create(ctx context.Context, definition Definition) (Backend, error)
	Update(ctx context.Context, identifier string, definition Definition) (Backend, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}

// EndpointURL returns the URL where to retrieve objects of type Backend and the identifier of the given Backend.
// It implements the api.Object interface on *Backend, making it usable with the generic API client.
func (b *Backend) EndpointURL(ctx context.Context) (*url.URL, error) {
	u, err := url.ParseRequestURI("/api/LBaaS/v1/backend.json")

	if op, err := types.OperationFromContext(ctx); err == nil && op == types.OperationList {
		filters := make(url.Values)

		if b.LoadBalancer.Identifier != "" {
			filters.Add("load_balancer", b.LoadBalancer.Identifier)
		}

		if b.Mode != "" {
			filters.Add("mode", string(b.Mode))
		}

		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	} else if err != nil {
		return nil, err
	}

	return u, err
}

// FilterAPIRequestBody generates the request body for creating a new Backend, which differs from the Backend object.
func (b *Backend) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	if op, err := types.OperationFromContext(ctx); err == nil && op == types.OperationCreate {
		return map[string]string{
			"name":          b.Name,
			"load_balancer": b.LoadBalancer.Identifier,
			"mode":          string(b.Mode),
			"state":         "4", // "newly created"
		}, nil
	} else if err != nil {
		return nil, err
	}

	return b, nil
}
