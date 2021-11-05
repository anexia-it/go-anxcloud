package resource

import (
	"context"
	genericAPI "github.com/anexia-it/go-anxcloud/pkg/api"
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
	"github.com/go-logr/logr"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/client"
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

func (i Info) EndpointURL(ctx context.Context, op types.Operation, options types.Options) (*url.URL, error) {
	u, err := url.ParseRequestURI(pathPrefix)

	switch op {
	// OperationCreate is not supported because the API does not exist in the engine.
	// OperationDestroy and OperationUpdate is not yet implemented
	case types.OperationCreate, types.OperationDestroy, types.OperationUpdate:
		return nil, genericAPI.ErrOperationNotSupported
	}

	if op == types.OperationList {
		query := u.Query()

		if len(i.Tags) > 1 {
			logr.FromContextOrDiscard(ctx).Info("Listing with multiple tags isn't supported. Only first one used")
		}

		if len(i.Tags) > 0 {
			query.Add("tag_name", i.Tags[0])
		}
		u.RawQuery = query.Encode()
	}
	return u, err
}

// NewAPI creates a new tags API instance with the given client.
func NewAPI(c client.Client) API {
	return api{c}
}
