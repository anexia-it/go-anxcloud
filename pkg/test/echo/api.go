// Package echo contains API functionality for issuing echo requests to the API.
package echo

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for echo calls.
type API interface {
	Echo(ctx context.Context) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new echo API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
