package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the common URL for operations on NodePool resource
func (np *NodePool) EndpointURL(ctx context.Context) (*url.URL, error) {
	return endpointURL(ctx, np, "node_pool")
}
