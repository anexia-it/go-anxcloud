// Package kubernetes implements API functions residing under /kubernetes.
// This path contains methods for managing kubernetes.
package kubernetes

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/kubernetes/cluster"
	"go.anx.io/go-anxcloud/pkg/kubernetes/nodepool"
)

// API contains methods for IP manipulation.
type API interface {
	Cluster() cluster.API
	Nodepool() nodepool.API
}

type api struct {
	cluster  cluster.API
	nodepool nodepool.API
}

func (a api) Cluster() cluster.API {
	return a.cluster
}

func (a api) Nodepool() nodepool.API {
	return a.nodepool
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client, opts common.ClientOpts) API {
	return &api{
		cluster:  cluster.NewAPI(c, opts),
		nodepool: nodepool.NewAPI(c, opts),
	}
}
