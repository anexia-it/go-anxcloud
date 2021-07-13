// Package zone implements API functions residing under /zone.
// This path contains methods for managing IPs.
package zone

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/client"
	uuid "github.com/satori/go.uuid"
)

// API contains methods for zone and record management
type API interface {
	List(ctx context.Context) ([]Zone, error)
	Get(ctx context.Context, name string) (Zone, error)
	Create(ctx context.Context, create Definition) (Zone, error)
	Update(ctx context.Context, name string, update Definition) (Zone, error)
	Delete(ctx context.Context, name string) error
	Apply(ctx context.Context, name string, changeset ChangeSet) ([]Record, error)
	Import(ctx context.Context, name string, zoneData Import) (Revision, error)
	ListRecords(ctx context.Context, name string) ([]Record, error)
	NewRecord(ctx context.Context, zone string, record RecordRequest) (Zone, error)
	UpdateRecord(ctx context.Context, zone string, id uuid.UUID, record RecordRequest) (Zone, error)
	DeleteRecord(ctx context.Context, zone string, id uuid.UUID) error
	// Export zone
	// Export zone for specific region
}

type api struct {
	client client.Client
}

func NewAPI(c client.Client) API {
	return &api{c}
}
