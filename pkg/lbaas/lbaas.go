package lbaas

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/acl"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/backend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/bind"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/frontend"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/loadbalancer"
	"github.com/anexia-it/go-anxcloud/pkg/lbaas/server"
)

type API interface {
	LoadBalancer() loadbalancer.API
	Frontend() frontend.API
	Backend() backend.API
	Server() server.API
	Bind() bind.API
	ACL() acl.API
}

type api struct {
	loadBalancer loadbalancer.API
	frontend     frontend.API
	backend      backend.API
	server       server.API
	bind         bind.API
	acl          acl.API
}

func (a api) ACL() acl.API {
	return a.acl
}

func (a api) Bind() bind.API {
	return a.bind
}

func (a api) Backend() backend.API {
	return a.backend
}

func (a api) Server() server.API {
	return a.server
}

func (a api) LoadBalancer() loadbalancer.API {
	return a.loadBalancer
}

func (a api) Frontend() frontend.API {
	return a.frontend
}

// NewAPI creates a new vmlist API instance with the given client.
func NewAPI(c client.Client) API {

	return &api{
		loadBalancer: loadbalancer.NewAPI(c),
		frontend:     frontend.NewAPI(c),
		backend:      backend.NewAPI(c),
		server:       server.NewAPI(c),
		bind:         bind.NewAPI(c),
		acl:          acl.NewAPI(c),
	}
}
