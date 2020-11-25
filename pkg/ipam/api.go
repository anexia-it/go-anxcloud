// Package ipam implements API functions residing under /ipam.
// This path contains methods for managing IPs.
package ipam

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/address"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/prefix"
)

// API contains methods for IP manipulation.
type API interface {
	Address() address.API
	Prefix() prefix.API
}

type api struct {
	address address.API
	prefix  prefix.API
}

func (a api) Address() address.API {
	return a.address
}

func (a api) Prefix() prefix.API {
	return a.prefix
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		address.NewAPI(c),
		prefix.NewAPI(c),
	}
}
