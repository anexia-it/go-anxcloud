package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

var (
	// ErrTemplateOperationNotSupported is returned if operations other than Get and List are performed
	ErrTemplateOperationNotSupported = fmt.Errorf("%w: Template only supports Get and List operations", api.ErrOperationNotSupported)
)

// EndpointURL returns the URL where to retrieve objects of type Template (only Get and List operations supported)
func (t *Template) EndpointURL(ctx context.Context) (*url.URL, error) {
	_, err := t.operationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return url.ParseRequestURI(
		fmt.Sprintf("/api/vsphere/v1/provisioning/templates.json/%s/%s", t.Location.Identifier, t.Type),
	)
}

// HasPagination disables pagination for Template API (not supported by engine)
func (t *Template) HasPagination(ctx context.Context) (bool, error) {
	return false, nil
}

// FilterRequestURL removes the Identifier from URL on Get operations (template needs to be parsed from list response)
func (t *Template) FilterRequestURL(ctx context.Context, url *url.URL) (*url.URL, error) {
	op, err := t.operationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationGet {
		url.Path = path.Dir(url.Path)
	}

	q := url.Query()
	q.Set("page", "1")
	q.Set("limit", "1000")
	url.RawQuery = q.Encode()

	return url, nil
}

// DecodeAPIResponse is used to filter a single template on Get operations
func (t *Template) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	op, err := t.operationFromContext(ctx)
	if err != nil {
		return err
	}

	if op == types.OperationGet {
		tpl, err := t.extractTemplateByID(data)
		if err != nil {
			return err
		}
		*t = *tpl
		return nil
	}

	return json.NewDecoder(data).Decode(t)
}

func (t *Template) extractTemplateByID(data io.Reader) (*Template, error) {
	var templates []*Template
	err := json.NewDecoder(data).Decode(&templates)
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		if template.Identifier == t.Identifier {
			return template, nil
		}
	}

	return nil, fmt.Errorf("Template with given Identifier not found in response: %w", api.ErrNotFound)
}

func (t *Template) operationFromContext(ctx context.Context) (types.Operation, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return "", err
	}

	if op != types.OperationGet && op != types.OperationList {
		return "", ErrTemplateOperationNotSupported
	}

	return op, nil
}
