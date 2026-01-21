package nodepool

import (
	"context"
	"fmt"

	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/pagination"
)

// API contains methods for kubernetes nodepool
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]NodePoolInfo, error)
	GetByID(ctx context.Context, identifier string) (Nodepool, error)
	Create(ctx context.Context, definition Definition) (Nodepool, error)
	Update(ctx context.Context, identifier string, definition Definition) (Nodepool, error)
	DeleteByID(ctx context.Context, identifier string) error
}

type api struct {
	client client.Client
	path   string
}

const (
	pathFormat = "api/kubernetes%s/v1/node_pool.json"
)

// NewAPI creates a new kubernetes nodepool API instance with the given client.
func NewAPI(c client.Client, opt common.ClientOpts) API {
	envPath := ""

	if opt.Environment != common.EnvironmentProd {
		envPath = "-" + string(opt.Environment)
	}

	path := fmt.Sprintf(pathFormat, envPath)
	return &api{c, path}
}
