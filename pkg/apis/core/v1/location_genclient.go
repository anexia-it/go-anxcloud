package v1

import (
	"context"
	"net/url"
	"strings"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

// EndpointURL returns the default URL for core location operations
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

// FilterRequestURL rewrites the request URL to use the /location/by-code.json endpoint
// when Get operation by Code is requested
func (l *Location) FilterRequestURL(ctx context.Context, url *url.URL) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationGet && l.Identifier == "" && l.Code != "" {
		url.Path = strings.Replace(url.Path, "location.json", "location/by-code.json", 1)
	}

	return url, nil
}
