package loadbalancer

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/pagination"
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
