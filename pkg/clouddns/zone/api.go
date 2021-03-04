// Package zone implements API functions residing under /zone.
// This path contains methods for managing IPs.
package zone

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for zone and record management
type API interface {
	List(ctx context.Context) ([]Response, error)
	Get(ctx context.Context, name string) (Response, error)
	Create(ctx context.Context, create Definition) (Response, error)
	Update(ctx context.Context, name string, update Definition) (Response, error)
	Delete(ctx context.Context, name string) error
	Apply(ctx context.Context, name string, changeset ChangeSet) (Response, error)
	Import(ctx context.Context, name string, zoneData Import) error
	ListRecords(ctx context.Context, name string) ([]Record, error)
	NewRecord(ctx context.Context, zone string, record RecordRequest) (Response, error)
	UpdateRecord(ctx context.Context, zone string, id string, record RecordRequest) (Response, error)
	DeleteRecord(ctx context.Context, zone string, id string) error
}

type api struct {
	client client.Client
}

func NewAPI(c client.Client) API {
	return &api{c}
}
