package cluster

import (
	"context"
	"fmt"

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
}

type api struct {
	client client.Client
	path   string
}

const (
	pathFormat = "api/kubernetes%s/v1/cluster.json"
)

type ClientOpts struct {
	Environment Environment
}

type Environment string

const EnvironmentDev = Environment("dev")

// NewAPI creates a new kubernetes cluster API instance with the given client.
func NewAPI(c client.Client, opts ...ClientOpts) API {
	envPath := ""

	if len(opts) > 1 {
		panic("too many options, only one supported")
	} else if len(opts) == 1 {
		opt := opts[0]
		envPath = "-" + string(opt.Environment)
	}

	path := fmt.Sprintf(pathFormat, envPath)
	return &api{c, path}
}
