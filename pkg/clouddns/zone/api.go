// Package zone implements API functions residing under /zone.
// This path contains methods for managing IPs.
package zone

import (
	"context"
	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for zone and record management
type API interface {
	List(ctx context.Context) ([]Zone, error)
	Get(ctx context.Context, name string) (Zone, error)
	Create()
	Update()
	Delete(ctx context.Context, name string) error
	Apply(ctx context.Context, name string, changeset Changeset)
	Import(ctx context.Context, name string, zoneData string)
	NewRecord()
	UpdateRecord()
	DeleteRecord()
}

type api struct {
	client client.Client
}

func (a api) Create() {
	panic("implement me")
}

func (a api) Update() {
	panic("implement me")
}

func (a api) NewRecord() {
	panic("implement me")
}

func (a api) UpdateRecord() {
	panic("implement me")
}

func (a api) DeleteRecord() {
	panic("implement me")
}

func NewAPI(c client.Client) API {
	return &api{c}
}