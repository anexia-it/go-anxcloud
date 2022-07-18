package v1

import (
	"context"
	"net/url"
)

// EndpointURL returns the common URL for operations on NodePool resource
func (np *NodePool) EndpointURL(ctx context.Context) (*url.URL, error) {
	return endpointURL(ctx, "/api/kubernetes/v1/node_pool.json")
}

// FilterAPIRequestBody adds the CommonRequestBody
// and unwraps the identifiers of related Objects
func (np *NodePool) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return requestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			NodePool
			Cluster string `json:"cluster,omitempty"`
		}{
			NodePool: *np,
			Cluster:  np.Cluster.Identifier,
		}
	})
}
