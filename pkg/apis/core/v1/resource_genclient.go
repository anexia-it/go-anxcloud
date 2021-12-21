package v1

import (
	"context"
	"net/url"

	"github.com/go-logr/logr"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (i Info) EndpointURL(ctx context.Context) (*url.URL, error) {
	u, err := url.ParseRequestURI("/api/core/v1/resource.json")
	if err != nil {
		return nil, err
	}

	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	switch op {
	// OperationCreate is not supported because the API does not exist in the engine.
	// OperationDestroy and OperationUpdate is not yet implemented
	case types.OperationCreate, types.OperationDestroy, types.OperationUpdate:
		return nil, api.ErrOperationNotSupported
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
