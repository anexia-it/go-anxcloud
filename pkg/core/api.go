// Package core contains API functionality for /core.
package core

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/core/resource"
	"go.anx.io/go-anxcloud/pkg/core/service"
	"go.anx.io/go-anxcloud/pkg/core/tags"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/location"
)

// API contains methods for accessing features under /core.
type API interface {
	Resource() resource.API
	Service() service.API
	Tags() tags.API
	Location() location.API
}

type api struct {
	resource resource.API
	service  service.API
	tags     tags.API
	location location.API
}

func (a api) Resource() resource.API {
	return a.resource
}

func (a api) Service() service.API {
	return a.service
}

func (a api) Tags() tags.API {
	return a.tags
}

func (a api) Location() location.API {
	return a.location
}

// NewAPI creates a new API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		resource.NewAPI(c),
		service.NewAPI(c),
		tags.NewAPI(c),
		location.NewAPI(c),
	}
}
