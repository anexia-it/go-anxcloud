package resource

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for tag control.
type API interface {
	List(ctx context.Context, page, limit int) ([]Summary, error)
	Get(ctx context.Context, id string) (Info, error)
	AttachTag(ctx context.Context, resourceID, tagName string) ([]Summary, error)
	DetachTag(ctx context.Context, resourceID, tagName string) error
}

type api struct {
	client client.Client
}

// NewAPI creates a new tags API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
