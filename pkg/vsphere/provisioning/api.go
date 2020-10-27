// Package provisioning contains  APi funcationality for the provisioning of VMs.
package provisioning

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/ips"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/location"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/progress"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/vm"
)

// API contains methods for VM provisioning.
type API interface {
	IPs() ips.API
	Location() location.API
	Progress() progress.API
	VM() vm.API
}

type api struct {
	ips      ips.API
	location location.API
	progress progress.API
	vm       vm.API
}

func (a api) IPs() ips.API {
	return a.ips
}

func (a api) Location() location.API {
	return a.location
}

func (a api) Progress() progress.API {
	return a.progress
}

func (a api) VM() vm.API {
	return a.vm
}

// NewAPI creates a new provisioning API instance with the given client.
func NewAPI(c client.Client) API {
	return api{
		ips.NewAPI(c),
		location.NewAPI(c),
		progress.NewAPI(c),
		vm.NewAPI(c),
	}
}
