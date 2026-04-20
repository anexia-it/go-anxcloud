package lbaas

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericResource"
	"go.anx.io/go-anxcloud/pkg/lbaas/acl"
	"go.anx.io/go-anxcloud/pkg/lbaas/backend"
	"go.anx.io/go-anxcloud/pkg/lbaas/bind"
	"go.anx.io/go-anxcloud/pkg/lbaas/frontend"
	"go.anx.io/go-anxcloud/pkg/lbaas/loadbalancer"
	"go.anx.io/go-anxcloud/pkg/lbaas/rule"
	"go.anx.io/go-anxcloud/pkg/lbaas/server"
)

type API interface {
	LoadBalancer() genericResource.API[loadbalancer.Loadbalancer, loadbalancer.Definition]
	Frontend() frontend.API
	Backend() backend.API
	Server() server.API
	Bind() bind.API
	ACL() acl.API
	Rule() genericResource.API[rule.Rule, rule.Definition]
}

type api struct {
	loadBalancer genericResource.API[loadbalancer.Loadbalancer, loadbalancer.Definition]
	frontend     frontend.API
	backend      backend.API
	server       server.API
	bind         bind.API
	acl          acl.API
	rule         genericResource.API[rule.Rule, rule.Definition]
}

func (a api) Rule() genericResource.API[rule.Rule, rule.Definition] {
	return a.rule
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

func (a api) LoadBalancer() genericResource.API[loadbalancer.Loadbalancer, loadbalancer.Definition] {
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
		rule:         rule.NewAPI(c),
	}
}
