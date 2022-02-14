package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (l *Location) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Locations can only be retrieved via the public engine, nothing else
	if op != types.OperationGet && op != types.OperationList {
		return nil, api.ErrOperationNotSupported
	}

	return url.Parse("/api/core/v1/location.json")
}
