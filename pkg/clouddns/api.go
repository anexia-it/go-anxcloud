// Package clouddns contains API functionality for clouddns.
package clouddns

import (
	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/clouddns/zone"
)

// API contains methods for VMs.
type API interface {
	//Countries()
	//Regions()
	Zone() zone.API
	//Pool()
	//Instance()
	//Nameserverset()
}

type api struct {
	zone zone.API
}

func (a api) Zone() zone.API {
	return a.zone
}

func NewAPI(c client.Client) API {
	return &api{}
}
