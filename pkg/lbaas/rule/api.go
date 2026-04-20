package rule

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericResource"
)

type api struct {
	client client.Client
}

// NewAPI creates a new load balancer backend API instance with the given client.
func NewAPI(c client.Client) genericResource.API[Rule, Definition] {
	return &api{c}
}
