package rule

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
)

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) genericresource.API[Rule, Definition] {
	return &api{c}
}
