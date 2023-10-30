package v2

import (
	"context"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"net/url"
)

// FilterAPIRequestBody adds the CommonRequestBody
func (c *Cluster) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return gs.RequestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			Cluster
		}{
			Cluster: *c,
		}
	})
}

// EndpointURL returns the common URL for operations on the Cluster resource
func (c *Cluster) EndpointURL(ctx context.Context) (*url.URL, error) {
	return gs.EndpointURL(ctx, c, "/api/LBaaSv2/v1/clusters.json")
}
