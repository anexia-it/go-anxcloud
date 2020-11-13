package powercontrol

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for VM power control.
type API interface {
	Get(ctx context.Context, vmIdentifier string) (State, error)
	Set(ctx context.Context, vmIdentifier string, request Request) (Task, error)
	AwaitCompletion(ctx context.Context, vmID, taskID string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new powercontrol API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
