package v2

import (
	"context"
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
	"net/url"
)

// FilterAPIRequestBody adds the CommonRequestBody
func (n *Node) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return gs.RequestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			Node
		}{
			Node: *n,
		}
	})
}

// EndpointURL returns the common URL for operations on the Node resource
func (n *Node) EndpointURL(ctx context.Context) (*url.URL, error) {
	return gs.EndpointURL(ctx, n, "/api/LBaaSv2/v1/nodes.json")
}
