// Package core contains API functionality for /core.
package core

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/location"
)

// API contains methods for accessing features under /core.
type API interface {
	Location() location.API
}

type api struct {
	location location.API
}

func (a api) Location() location.API {
	return a.location
}

// NewAPI creates a new API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		location.NewAPI(c),
	}
}
