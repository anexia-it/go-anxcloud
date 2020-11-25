package templates

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/client"
)

// API contains methods for template querying.
type API interface {
	List(ctx context.Context, locationID string, templateType TemplateType, page, limit int) ([]Template, error)
}

type api struct {
	client client.Client
}

// NewAPI creates a new template API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
