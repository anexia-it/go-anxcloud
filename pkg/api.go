// Package pkg contains all API functionality and helpers.
package pkg

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere"
)

// API contains all API calls structured their location in the API.
type API interface {
	VSphere() vsphere.API
}

type api struct {
	vsphere vsphere.API
}

func (a api) VSphere() vsphere.API {
	return a.vsphere
}

// NewAPI creates a new API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		vsphere.NewAPI(c),
	}
}
