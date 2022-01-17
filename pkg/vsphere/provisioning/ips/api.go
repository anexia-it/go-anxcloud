package ips

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for IP manipulation (but only in the context of provisioning).
type API interface {
	GetFree(ctx context.Context, location, vlan string) ([]IP, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new IP API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
