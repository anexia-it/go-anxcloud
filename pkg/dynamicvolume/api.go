// Package dynamicvolume implements API functions residing under /dynamic_volume.
// This path contains methods for managing storageserverinterfaces.
package dynamicvolume

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/dynamicvolume/storageserverinterface"
)

// API contains methods for IP manipulation.
type API interface {
	StorageServerInterface() storageserverinterface.API
}

type api struct {
	storageserverinterface storageserverinterface.API
}

func (a api) StorageServerInterface() storageserverinterface.API {
	return a.storageserverinterface
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		storageserverinterface: storageserverinterface.NewAPI(c),
	}
}
