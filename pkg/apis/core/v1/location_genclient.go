package v1

import (
	"context"
	"fmt"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the default URL for core location operations.
// It supports Get by-code if `Code` is set and `Identifier` is not.
func (l *Location) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Locations can only be retrieved via the public engine, nothing else
	if op != types.OperationGet && op != types.OperationList {
		return nil, api.ErrOperationNotSupported
	}

	endpointSuffix := "location.json"
	if op == types.OperationGet && l.Identifier == "" && l.Code != "" {
		endpointSuffix = "location/by-code.json"
	}

	return url.Parse(fmt.Sprintf("/api/core/v1/%s", endpointSuffix))
}

// GetIdentifier returns the objects identifier
func (l Location) GetIdentifier(ctx context.Context) (string, error) {
	op, err := types.OperationFromContext(ctx)
	if l.Identifier != "" || err != nil {
		return l.Identifier, nil
	}

	if op == types.OperationGet {
		return l.Code, nil
	}

	return "", nil
}
