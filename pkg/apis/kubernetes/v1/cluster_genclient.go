package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the common URL for operations on Cluster resource
func (c *Cluster) EndpointURL(ctx context.Context) (*url.URL, error) {
	return endpointURL(ctx, "/api/kubernetes/v1/cluster.json")
}

// FilterAPIRequestBody adds the CommonRequestBody
func (c *Cluster) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return requestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			Cluster
			Location string `json:"location,omitempty"`
		}{
			Cluster:  *c,
			Location: c.Location.Identifier,
		}
	})
}
