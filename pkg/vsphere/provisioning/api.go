// Package provisioning contains  APi funcationality for the provisioning of VMs.
package provisioning

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/disktype"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/ips"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/location"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/progress"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/templates"
	"go.anx.io/go-anxcloud/pkg/vsphere/provisioning/vm"
)

// API contains methods for VM provisioning.
type API interface {
	DiskType() disktype.API
	IPs() ips.API
	Location() location.API
	Progress() progress.API
	Templates() templates.API
	VM() vm.API
}

type api struct {
	diskType  disktype.API
	ips       ips.API
	location  location.API
	progress  progress.API
	templates templates.API
	vm        vm.API
}

func (a api) DiskType() disktype.API {
	return a.diskType
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

func (a api) Templates() templates.API {
	return a.templates
}

func (a api) VM() vm.API {
	return a.vm
}

// NewAPI creates a new provisioning API instance with the given client.
func NewAPI(c client.Client) API {
	return api{
		disktype.NewAPI(c),
		ips.NewAPI(c),
		location.NewAPI(c),
		progress.NewAPI(c),
		templates.NewAPI(c),
		vm.NewAPI(c),
	}
}
