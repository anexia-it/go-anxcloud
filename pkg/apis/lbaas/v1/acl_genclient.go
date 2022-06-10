package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/object/filter"
)

// EndpointURL returns the URL where to retrieve objects of type ACL and the identifier of the given ACL.
// It implements the api.Object interface on *ACL, making it usable with the generic API client.
func (a *ACL) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("/api/LBaaS/v1/ACL.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		helper, err := filter.NewHelper(a)
		if err != nil {
			return nil, err
		}

		filters := helper.BuildQuery().Encode()

		if filters != "" {
			query := u.Query()
			query.Set("filters", filters)
			u.RawQuery = query.Encode()
		}
	}

	return u, nil
}

// FilterAPIRequestBody generates the request body for ACLs, replacing linked Objects with just their identifier.
func (a *ACL) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	return requestBody(ctx, func() interface{} {
		return &struct {
			commonRequestBody
			ACL
			Backend  string `json:"backend,omitempty"`
			Frontend string `json:"frontend,omitempty"`
		}{
			ACL:      *a,
			Backend:  a.Backend.Identifier,
			Frontend: a.Frontend.Identifier,
		}
	})
}
