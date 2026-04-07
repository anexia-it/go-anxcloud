package network

import (
	"context"
	"fmt"

	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/pagination"
)

// API contains methods for kubernetes networks
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]common.PartialResource, error)
	GetByID(ctx context.Context, identifier string) (NodepoolNetwork, error)
	Create(ctx context.Context, definition NodepoolNetworkDefinition) (NodepoolNetwork, error)
	Update(ctx context.Context, identifier string, definition NodepoolNetworkDefinition) (
		NodepoolNetwork, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
	path   string
}

const (
	pathFormat = "api/kubernetes%s/v2/node_pool_network"
)

// NewAPI creates a new kubernetes nodepool network API instance with the given client.
func NewAPI(c client.Client, opt common.ClientOpts) API {
	envPath := ""

	if opt.Environment != common.EnvironmentProd {
		envPath = "-" + string(opt.Environment)
	}

	path := fmt.Sprintf(pathFormat, envPath)
	return &api{c, path}
}
