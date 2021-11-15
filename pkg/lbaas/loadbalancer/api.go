package loadbalancer

import (
	"context"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/pagination"
)

// API contains load balancer actions.
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]LoadBalancerInfo, error)
	GetByID(ctx context.Context, identifier string) (Loadbalancer, error)

}

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}

// EndpointURL returns the URL where to retrieve objects of type Loadbalancer and the identifier of the given Loadbalancer.
// It implements the api.Object interface on *Loadbalancer, making it usable with the generic API client.
func (lb *Loadbalancer) EndpointURL(ctx context.Context, op types.Operation, options types.Options) (*url.URL, error) {
	url, err := url.ParseRequestURI("/api/LBaaS/v1/loadbalancer.json")
	return url, err
}

// FilterAPIRequestBody generates the request body for creating a new Loadbalancer, which differs from the Loadbalancer object.
func (lb *Loadbalancer) FilterAPIRequestBody(op types.Operation, options types.Options) (interface{}, error) {
	if op == types.OperationCreate {
		return map[string]string{
			"name":       lb.Name,
			"ip_address": lb.IpAddress,
			"state":      "2",
		}, nil
	}

	return lb, nil
}
