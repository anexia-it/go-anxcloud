// Package automation implements API functions residing under /automation.
// This path contains methods for managing automations.
package automation

import (
	"go.anx.io/go-anxcloud/pkg/automation/rules"
	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for IP manipulation.
type API interface {
	Rules() rules.API
}

type api struct {
	rules rules.API
}

func (a api) Rules() rules.API {
	return a.rules
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{
		rules.NewAPI(c),
	}
}
