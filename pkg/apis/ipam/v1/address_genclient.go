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

	switch op {
	case types.OperationList:
		fh, err := filter.NewHelper(a)
		if err != nil {
			return nil, err
		}

		query := fh.BuildQuery()

		if v, ok, err := fh.Get("type"); ok && err == nil {
			query.Del("type")

			if v.(AddressType) == TypePublic {
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
	default:
		return url.Parse("/api/ipam/v1/address.json")
	}
}

func (a *Address) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch op {
	case types.OperationCreate:
		data := struct {
			Prefix              string `json:"prefix"`
			Name                string `json:"name"`
			DescriptionCustomer string `json:"description_customer,omitempty"`
			Role                string `json:"role,omitempty"`
		}{
			Prefix:              a.Prefix.Identifier,
			Name:                a.Name,
			DescriptionCustomer: a.DescriptionCustomer,
			Role:                a.RoleText,
		}

		return data, nil
	case types.OperationUpdate:
		data := struct {
			Identifier          string `json:"identifier"`
			DescriptionCustomer string `json:"description_customer,omitempty"`
			Role                string `json:"role,omitempty"`
		}{
			Identifier:          a.Identifier,
			DescriptionCustomer: a.DescriptionCustomer,
			Role:                a.RoleText,
		}
		return data, nil
	default:
		return a, nil
	}
}
