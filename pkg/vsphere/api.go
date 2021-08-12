// Package vsphere contains API functionality for vsphere.
package vsphere

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/info"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/powercontrol"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/search"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/vmlist"
)

// API contains methods for VMs.
type API interface {
	Info() info.API
	PowerControl() powercontrol.API
	Provisioning() provisioning.API
	Search() search.API
	VMList() vmlist.API
}

type api struct {
	info         info.API
	powercontrol powercontrol.API
	provisioning provisioning.API
	search       search.API
	vmlist       vmlist.API
}

func (a api) Info() info.API {
	return a.info
}

func (a api) PowerControl() powercontrol.API {
	return a.powercontrol
}

func (a api) Provisioning() provisioning.API {
	return a.provisioning
}

func (a api) Search() search.API {
	return a.search
}

func (a api) VMList() vmlist.API {
	return a.vmlist
}

// NewAPI creates a new vsphere API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		info.NewAPI(c),
		powercontrol.NewAPI(c),
		provisioning.NewAPI(c),
		search.NewAPI(c),
		vmlist.NewAPI(c),
	}
}
