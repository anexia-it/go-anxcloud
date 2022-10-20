package v1

import (
	"context"
	"net/url"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	"go.anx.io/go-anxcloud/pkg/apis/internal/gs"
	"go.anx.io/go-anxcloud/pkg/utils/object/filter"
)

func endpointURL(ctx context.Context, o types.Object, apiPath string) (*url.URL, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationUpdate {
		return nil, api.ErrOperationNotSupported
	}

	// we can ignore the error since the URL is hard-coded known as valid
	u, _ := url.Parse(apiPath)

	if op == types.OperationList {
		helper, err := filter.NewHelper(o)
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

// HasState can be embedded to add the state object to a resource
type HasState struct {
	State gs.State `json:"state"`
}

// StateSuccess checks if the state is one of the successful ones
func (hs HasState) StateSuccess() bool { return hs.State.ID == "0" }

// StateProgressing checks if the state is marking any change currently being applied
func (hs HasState) StateProgressing() bool { return hs.State.ID == "2" || hs.State.ID == "3" }

// StateFailure checks if the state is marking any failure
func (hs HasState) StateFailure() bool { return hs.State.ID == "1" }

// commonRequestBody is embedded in the request body types of kubernetes resources
type commonRequestBody struct {
	// this allows us to provide the state as string
	State string `json:"state,omitempty"`
}

// requestBody removes the request body for read operations on kubernetes resources
func requestBody(ctx context.Context, br func() interface{}) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate || op == types.OperationUpdate {
		response := br()

		return response, nil
	}

	return nil, nil
}
