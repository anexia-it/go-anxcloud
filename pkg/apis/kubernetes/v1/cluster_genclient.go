package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the common URL for operations on Cluster resource
func (c *Cluster) EndpointURL(ctx context.Context) (*url.URL, error) {
	return endpointURL(ctx, c, "/api/kubernetes/v1/cluster.json")
}
