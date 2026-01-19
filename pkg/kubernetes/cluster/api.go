package cluster

import (
	"context"
	"fmt"

	"go.anx.io/go-anxcloud/pkg/apis/common"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/pagination"
)

// API contains methods for kubernetes cluster
type API interface {
	pagination.Pageable
	Get(ctx context.Context, page, limit int) ([]ClusterInfo, error)
	GetByID(ctx context.Context, identifier string) (Cluster, error)
	Create(ctx context.Context, definition Definition) (Cluster, error)
	Update(ctx context.Context, identifier string, definition Definition) (Cluster, error)
	DeleteByID(ctx context.Context, identifier string) error

	RequestKubeConfig(ctx context.Context, cluster *Cluster) error
}

type api struct {
	client client.Client
	path   string
}

const (
	pathFormat = "api/kubernetes%s/v1/cluster.json"
)

// NewAPI creates a new kubernetes cluster API instance with the given client.
func NewAPI(c client.Client, opt common.ClientOpts) API {
	envPath := ""

	if opt.Environment != common.EnvironmentProd {
		envPath = "-" + string(opt.Environment)
	}

	path := fmt.Sprintf(pathFormat, envPath)
	return &api{c, path}
}
