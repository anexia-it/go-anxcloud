// Package clouddns contains API functionality for clouddns.
package clouddns

import (
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/clouddns/zone"
)

// API contains methods managing zones and records.
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
	return &api{zone.NewAPI(c)}
}
