package cpuperformancetype

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for template querying.
type API interface {
	List(ctx context.Context) ([]CPUPerformanceType, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new template API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
