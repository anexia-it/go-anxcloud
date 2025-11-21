// Package kubernetes implements API functions residing under /kubernetes.
// This path contains methods for managing kubernetes.
package kubernetes

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/kubernetes/cluster"
)

// API contains methods for IP manipulation.
type API interface {
	Cluster() cluster.API
}

type api struct {
	cluster cluster.API
}

func (a api) Cluster() cluster.API {
	return a.cluster
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		cluster.NewAPI(c),
	}
}
