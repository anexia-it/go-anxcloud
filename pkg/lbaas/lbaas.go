package lbaas

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/genericresource"
	"go.anx.io/go-anxcloud/pkg/lbaas/acl"
	"go.anx.io/go-anxcloud/pkg/lbaas/backend"
	"go.anx.io/go-anxcloud/pkg/lbaas/bind"
	"go.anx.io/go-anxcloud/pkg/lbaas/frontend"
	"go.anx.io/go-anxcloud/pkg/lbaas/loadbalancer"
	"go.anx.io/go-anxcloud/pkg/lbaas/rule"
	"go.anx.io/go-anxcloud/pkg/lbaas/server"
)

type API interface {
	LoadBalancer() genericresource.API[loadbalancer.Loadbalancer, loadbalancer.Definition]
	Frontend() genericresource.API[frontend.Frontend, frontend.Definition]
	Backend() genericresource.API[backend.Backend, backend.Definition]
	Server() genericresource.API[server.Server, server.Definition]
	Bind() genericresource.API[bind.Bind, bind.Definition]
	ACL() genericresource.API[acl.ACL, acl.Definition]
	Rule() genericresource.API[rule.Rule, rule.Definition]
}

type api struct {
	loadBalancer genericresource.API[loadbalancer.Loadbalancer, loadbalancer.Definition]
	frontend     genericresource.API[frontend.Frontend, frontend.Definition]
	backend      genericresource.API[backend.Backend, backend.Definition]
	server       genericresource.API[server.Server, server.Definition]
	bind         genericresource.API[bind.Bind, bind.Definition]
	acl          genericresource.API[acl.ACL, acl.Definition]
	rule         genericresource.API[rule.Rule, rule.Definition]
}

func (a api) Rule() genericresource.API[rule.Rule, rule.Definition] {
	return a.rule
}

func (a api) ACL() genericresource.API[acl.ACL, acl.Definition] {
	return a.acl
}

func (a api) Bind() genericresource.API[bind.Bind, bind.Definition] {
	return a.bind
}

func (a api) Backend() genericresource.API[backend.Backend, backend.Definition] {
	return a.backend
}

func (a api) Server() genericresource.API[server.Server, server.Definition] {
	return a.server
}

func (a api) LoadBalancer() genericresource.API[loadbalancer.Loadbalancer, loadbalancer.Definition] {
	return a.loadBalancer
}

func (a api) Frontend() genericresource.API[frontend.Frontend, frontend.Definition] {
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
