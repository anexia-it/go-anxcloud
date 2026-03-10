package storageserverinterface

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/pagination"
)

// API contains methods for storage server interfaces
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]common.PartialResource, error)
	GetByID(ctx context.Context, identifier string) (StorageServerInterface, error)
	Create(ctx context.Context, definition Definition) (StorageServerInterface, error)
	Update(ctx context.Context, identifier string, definition Definition) (StorageServerInterface, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
}

const (
	path = "api/dynamic_volume/v2/storage_server_interfaces"
)

// NewAPI creates a new storage server interfaces API instance with the given client.
func NewAPI(c client.Client) API {
	return &api{c}
}
