package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the base URL path of the resources API
func (a *Application) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/api/e5e/v1/application.json")
}
