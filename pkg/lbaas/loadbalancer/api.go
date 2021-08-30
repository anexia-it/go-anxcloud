package loadbalancer

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains load balancer actions.
type API interface {
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
