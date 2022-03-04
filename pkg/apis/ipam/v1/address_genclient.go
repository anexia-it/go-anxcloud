package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/object/filter"
)

func (a *Address) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		fh, err := filter.NewHelper(a)
		if err != nil {
			return nil, err
		}

		query := fh.BuildQuery()

		if v, ok, err := fh.Get("type"); ok && err == nil {
			delete(query, "type")

			if v.(AddressSpace) == AddressSpacePublic {
				query.Set("private", "false")
			} else {
				query.Set("private", "true")
			}
		} else if err != nil {
			return nil, err
		}

		u, _ := url.Parse("/api/ipam/v1/address/filtered.json")
		u.RawQuery = query.Encode()
		return u, nil
	}

	return url.Parse("/api/ipam/v1/address.json")
}

func (a *Address) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Creating a Prefix is done with a single location and VLAN only, despite the API returning arrays.
	if op == types.OperationCreate {
		data := struct {
			Address
			Prefix string `json:"prefix"`
			VLAN   string `json:"vlan,omitempty"`
		}{
			Address: *a,
			Prefix:  a.Prefix.Identifier,
			VLAN:    a.VLAN.Identifier,
		}

		return data, nil
	}

	return a, nil
}
