package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (z *Zone) EndpointURL(ctx context.Context) (*url.URL, error) {
	u, err := url.ParseRequestURI("/api/clouddns/v1/zone.json/")
	return u, err
}

func (z *Zone) FilterAPIRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// The Update endpoint is NOT at ".../zone.json/{zoneName}", but simply ".../zone.json"
	if op == types.OperationUpdate {
		// Strip the appended zoneName from the URL
		req.URL.Path = path.Dir(req.URL.Path)
	}

	return req, nil
}

func (z *Zone) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// The Create and Update endpoints expect the Zone's name to be in the request body under the key "zone_name"
	if op == types.OperationCreate || op == types.OperationUpdate {
		zWithZoneName := struct {
			Zone
			ZoneName string `json:"zone_name"`
		}{*z, z.Name}

		// `name` does not exist as a field on the Engine API for these requests,
		// so we strip it from the request body.
		zWithZoneName.Name = ""

		return zWithZoneName, nil
	}

	return z, nil
}

func (z *Zone) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// CloudDNS API's List response contains some non-functional pagination remnants, which are stripped here
	// Actual array of Zones is in the key 'results'
	if op == types.OperationList {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var m map[string]json.RawMessage
		err = json.Unmarshal(data, &m)
		if err != nil {
			return nil, err
		}

		data = m["results"]
		res.Body = io.NopCloser(bytes.NewReader(data))
		res.ContentLength = int64(len(data))
	}
	return res, nil
}

func (z *Zone) HasPagination(ctx context.Context) (bool, error) {
	return false, nil
}
