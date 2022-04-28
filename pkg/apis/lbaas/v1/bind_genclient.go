package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (b Bind) EndpointURL(ctx context.Context) (*url.URL, error) {

	// EndpointURL returns the URL where to retrieve objects of type Frontend and the identifier of the given Frontend.
	// It implements the api.Object interface on *Frontend, making it usable with the generic API client.
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse("api/LBaaS/v1/bind.json")
	if err != nil {
		return nil, err
	}

	if op == types.OperationList {
		filters := make(url.Values)
		if b.Frontend.Identifier != "" {
			filters.Add("frontend", b.Frontend.Identifier)
		}

		query := u.Query()
		query.Add("filters", filters.Encode())
		u.RawQuery = query.Encode()
	}

	return u, nil
}

// FilterAPIRequestBody generates the request body for creating a new FrontendBind, which differs from the Bind object.
func (b *Bind) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate || op == types.OperationUpdate {
		return struct {
			Bind
			Frontend string `json:"frontend"`

			// we never want to send the state field, so making sure to omit it here
			State string `json:"state,omitempty"`
		}{
			Bind:     *b,
			Frontend: b.Frontend.Identifier,
		}, nil
	}

	return b, nil
}

func (b *Bind) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return res, err
	}

	if op == types.OperationDestroy {
		err = res.Body.Close()
		if err != nil {
			return res, err
		}

		res.Body = io.NopCloser(bytes.NewReader([]byte("{}")))
		return res, nil
	}
	return res, nil
}
