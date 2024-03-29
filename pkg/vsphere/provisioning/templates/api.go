package templates

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/client"
)

// API contains methods for template querying.
type API interface {
	List(ctx context.Context, locationID string, templateType string, page, limit int) ([]Template, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new template API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
