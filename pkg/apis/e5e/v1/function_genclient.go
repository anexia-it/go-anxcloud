package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the base URL path of the resources API
func (f *Function) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/api/e5e/v1/function.json")
}
