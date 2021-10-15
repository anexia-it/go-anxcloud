package frontend

import (
	"context"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"
)

// API contains methods for load balancer frontend management.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]FrontendInfo, error)
	GetByID(ctx context.Context, identifier string) (Frontend, error)
	Create(ctx context.Context, definition Definition) (Frontend, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new frontend API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}

// EndpointURL returns the URL where to retrieve objects of type Frontend and the identifier of the given Frontend.
// It implements the api.Object interface on *Frontend, making it usable with the generic API client.
func (f *Frontend) EndpointURL(ctx context.Context, op types.Operation, options types.Options) (*url.URL, string, error) {
	url, err := url.ParseRequestURI("/api/LBaaS/v1/frontend.json")
	return url, f.Identifier, err
}

// FilterAPIRequestBody generates the request body for creating a new Frontend, which differs from the Frontend object.
func (f *Frontend) FilterAPIRequestBody(op types.Operation, options types.Options) (interface{}, error) {
	if op == types.OperationCreate {
		return map[string]string{
			"name":            f.Name,
			"load_balancer":   f.LoadBalancer.Identifier,
			"default_backend": f.DefaultBackend.Identifier,
			"state":           "4", // "newly created"
		}, nil
	}

	return f, nil
}
