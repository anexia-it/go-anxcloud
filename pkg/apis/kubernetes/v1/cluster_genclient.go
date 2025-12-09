package v1

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/pointer"
)

// ErrManagedPrefixSet is returned if a prefix is both set to managed and also set
// to an existing prefix on create
var ErrManagedPrefixSet = errors.New("managed prefixes cannot be set on create")

// EndpointURL returns the common URL for operations on Cluster resource
func (c *Cluster) EndpointURL(ctx context.Context) (*url.URL, error) {
	return endpointURL(ctx, c, "cluster")
}

// explicitlyFalse returns true if the value of the provided
// bool pointer is set to false, nil and true pointer return false
func explicitlyFalse(b *bool) bool {
	return b != nil && !pointer.BoolVal(b)
}

// prefixConfigurationValidCreate validates that prefixes configured as managed
// are not set. This is only relevant for `Create` operations.
func prefixConfigurationValidCreate(c *Cluster) bool {
	var (
		internalV4 = c.InternalIPv4Prefix == nil || explicitlyFalse(c.ManageInternalIPv4Prefix)
		externalV4 = c.ExternalIPv4Prefix == nil || explicitlyFalse(c.ManageExternalIPv4Prefix)
		externalV6 = c.ExternalIPv6Prefix == nil || explicitlyFalse(c.ManageExternalIPv6Prefix)
	)

	return internalV4 && externalV4 && externalV6
}

// FilterAPIRequestBody adds the CommonRequestBody
func (c *Cluster) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate && !prefixConfigurationValidCreate(c) {
		return nil, ErrManagedPrefixSet
	}

	return requestBody(ctx, func() interface{} {
		body := &struct {
			commonRequestBody
			Cluster
			Location string `json:"location,omitempty"`
		}{
			Cluster:  *c,
			Location: c.Location.Identifier,
		}

		if op == types.OperationUpdate {
			body.commonRequestBody = commonRequestBody{
				State: strconv.Itoa(c.State.Type),
			}
		}

		return body
	})
}
