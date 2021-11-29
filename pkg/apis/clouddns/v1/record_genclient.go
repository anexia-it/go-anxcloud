package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path"

	"github.com/anexia-it/go-anxcloud/pkg/api"
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

func (r *Record) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.ParseRequestURI(fmt.Sprintf("/api/clouddns/v1/zone.json/%s/records", r.ZoneName))
	if err != nil {
		return nil, err
	}

	// There is no endpoint to get details of a single record
	if op == types.OperationGet {
		return nil, api.ErrOperationNotSupported
	}

	if op == types.OperationList {
		query := u.Query()

		if r.Name != "" {
			query.Add("name", r.Name)
		}

		if r.RData != "" {
			query.Add("data", r.RData)
		}

		if r.Type != "" {
			query.Add("type", r.Type)
		}
		u.RawQuery = query.Encode()
	}
	return u, err
}

func (r *Record) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	// Response to POST and PUT are the _Zone_ details, which contain some of the updated Record details, but not all
	// To work around these inconsistencies, we just leave the receiver as it is
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return err
	}
	if op == types.OperationCreate || op == types.OperationUpdate {
		return nil
	}

	d := json.NewDecoder(data)
	err = d.Decode(r)
	if err != nil {
		return err
	}

	// Get zoneName from URL and put that into r.ZoneName
	if op == types.OperationList {
		url, err := types.URLFromContext(ctx)
		if err != nil {
			return err
		}
		r.ZoneName = path.Base(path.Dir(url.Path))
	}
	return nil
}

func (r *Record) HasPagination(ctx context.Context) (bool, error) {
	return false, nil
}
