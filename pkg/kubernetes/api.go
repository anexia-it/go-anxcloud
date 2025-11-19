// Package ipam implements API functions residing under /ipam.
// This path contains methods for managing IPs.
package ipam

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/kubernetes/cluster"
)

// API contains methods for IP manipulation.
type API interface {
	Cluster() cluster.API
}

type api struct {
	address cluster.API
}

func (a api) Cluster() cluster.API {
	return a.address
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		cluster.NewAPI(c),
	}
}
