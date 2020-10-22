// Package ipam implements API functions residing under /ipam.
// This path contains methods for managing IPs.
package ipam

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/ipam/address"
)

// API contains methods for IP manipulation.
type API interface {
	Address() address.API
}

type api struct {
	address address.API
}

func (a api) Address() address.API {
	return a.address
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		address.NewAPI(c),
	}
}
