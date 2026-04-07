// Package kubernetes implements API functions residing under /kubernetes.
// This path contains methods for managing kubernetes.
package kubernetes

import (
	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/kubernetes/cluster"
	"go.anx.io/go-anxcloud/pkg/kubernetes/disk"
	"go.anx.io/go-anxcloud/pkg/kubernetes/network"
	"go.anx.io/go-anxcloud/pkg/kubernetes/nodepool"
)

// API contains methods for IP manipulation.
type API interface {
	Cluster() cluster.API
	Nodepool() nodepool.API
	Disk() disk.API
	Network() network.API
}

type api struct {
	cluster  cluster.API
	nodepool nodepool.API
	disk     disk.API
	network  network.API
}

func (a api) Cluster() cluster.API {
	return a.cluster
}

func (a api) Nodepool() nodepool.API {
	return a.nodepool
}

func (a api) Disk() disk.API {
	return a.disk
}

func (a api) Network() network.API {
	return a.network
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client, opts common.ClientOpts) API {
	return &api{
		cluster:  cluster.NewAPI(c, opts),
		nodepool: nodepool.NewAPI(c, opts),
		disk:     disk.NewAPI(c, opts),
		network:  network.NewAPI(c, opts),
	}
}
