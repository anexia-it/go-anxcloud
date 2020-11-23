// Package disktype implements API functions residing under /provisioning/disk_type.
// This path contains methods for querying ips available disk types.
package disktype

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for disk_type querying.
type API interface {
	List(ctx context.Context, locationID string) ([]DiskType, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new disk_type API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
