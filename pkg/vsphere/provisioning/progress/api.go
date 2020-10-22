package progress

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for progress inquiries.
type API interface {
	AwaitCompletion(ctx context.Context, progressID string) (string, error)
	Get(ctx context.Context, identifier string) (Progress, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new progress API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
