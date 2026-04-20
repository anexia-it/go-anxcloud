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
	Frontend() genericResource.API[frontend.Frontend, frontend.Definition]
	Backend() genericResource.API[backend.Backend, backend.Definition]
	Server() genericResource.API[server.Server, server.Definition]
	Bind() genericResource.API[bind.Bind, bind.Definition]
	ACL() genericResource.API[acl.ACL, acl.Definition]
	Rule() genericResource.API[rule.Rule, rule.Definition]
}

type api struct {
	loadBalancer genericResource.API[loadbalancer.Loadbalancer, loadbalancer.Definition]
	frontend     genericResource.API[frontend.Frontend, frontend.Definition]
	backend      genericResource.API[backend.Backend, backend.Definition]
	server       genericResource.API[server.Server, server.Definition]
	bind         genericResource.API[bind.Bind, bind.Definition]
	acl          genericResource.API[acl.ACL, acl.Definition]
	rule         genericResource.API[rule.Rule, rule.Definition]
}

func (a api) Rule() genericResource.API[rule.Rule, rule.Definition] {
	return a.rule
}

func (a api) ACL() genericResource.API[acl.ACL, acl.Definition] {
	return a.acl
}

func (a api) Bind() genericResource.API[bind.Bind, bind.Definition] {
	return a.bind
}

func (a api) Backend() genericResource.API[backend.Backend, backend.Definition] {
	return a.backend
}

func (a api) Server() genericResource.API[server.Server, server.Definition] {
	return a.server
}

func (a api) LoadBalancer() genericResource.API[loadbalancer.Loadbalancer, loadbalancer.Definition] {
	return a.loadBalancer
}

func (a api) Frontend() genericResource.API[frontend.Frontend, frontend.Definition] {
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
