package v1

import (
	"context"
	"fmt"
	"net/url"

	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
)

func (v *VLAN) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == apiTypes.OperationList {
		query := url.Values{}

		if v.Status != StatusInvalid {
			query.Add("status", string(v.Status))
		}

		if l := len(v.Locations); l == 1 {
			query.Add("location", v.Locations[0].Identifier)
		} else if l > 1 {
			return nil, ErrFilterMultipleLocations
		}

		// we don't catch the error from url.Parse because this URL is hardcoded-valid.
		u, _ := url.Parse("/api/vlan/v1/vlan.json/filtered")
		u.RawQuery = query.Encode()

		return u, nil
	}

	return url.Parse("/api/vlan/v1/vlan.json")
}

func (v *VLAN) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Creating a VLAN is done with a single location only, despite the API returning an array.
	if op == apiTypes.OperationCreate {
		if len(v.Locations) != 1 {
			return nil, fmt.Errorf("%w: %v locations given", ErrLocationCount, len(v.Locations))
		}

		data := struct {
			VLAN
			Location string `json:"location"`
		}{
			VLAN:     *v,
			Location: v.Locations[0].Identifier,
		}

		data.VLAN.Locations = nil

		return data, nil
	}

	return v, nil
}
