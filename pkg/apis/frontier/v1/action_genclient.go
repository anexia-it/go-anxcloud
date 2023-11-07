package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the base URL path of the resources API
func (a *Action) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/api/frontier/v1/action.json")
}
