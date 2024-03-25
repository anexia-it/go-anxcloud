package v1

import (
	"context"
	"fmt"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/utils/object/filter"
)

func (p *Prefix) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		fh, err := filter.NewHelper(p)
		if err != nil {
			return nil, err
		}

		u, _ := url.Parse("/api/ipam/v1/prefix/filtered.json")
		u.RawQuery = fh.BuildQuery().Encode()
		return u, nil
	}

	return url.Parse("/api/ipam/v1/prefix.json")
}

func (p *Prefix) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Creating a Prefix is done with a single location and VLAN only, despite the API returning arrays.
	if op == types.OperationCreate {
		opts, err := types.OptionsFromContext(ctx)
		if err != nil {
			return nil, err
		}

		if len(p.Locations) != 1 {
			return nil, fmt.Errorf("%w: %v locations given", ErrLocationCount, len(p.Locations))
		}

		if len(p.VLANs) > 1 {
			return nil, fmt.Errorf("%w: %v VLANs given", ErrVLANCount, len(p.VLANs))
		}

		data := struct {
			Prefix
			CreateEmpty *bool  `json:"create_empty,omitempty"`
			Location    string `json:"location"`
			VLAN        string `json:"vlan,omitempty"`
		}{
			Prefix:   *p,
			Location: p.Locations[0].Identifier,
		}

		if len(p.VLANs) == 1 {
			data.VLAN = p.VLANs[0].Identifier
		}

		if ce, err := opts.Get(optionKeyCreateEmpty); err == nil {
			v := ce.(bool)
			data.CreateEmpty = &v
		}

		data.Prefix.Locations = nil
		data.Prefix.VLANs = nil

		return data, nil
	}

	return p, nil
}
