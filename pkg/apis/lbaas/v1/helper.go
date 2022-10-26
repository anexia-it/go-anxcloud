package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// commonRequestBody is embedded in the request body types of LBaaS objects.
type commonRequestBody struct {
	// we want to send the correct state for Update operations, but none for Create operations
	State string `json:"state,omitempty"`
}

func (crb *commonRequestBody) setState(state string) {
	crb.State = state
}

type commonRequestBodyInterface interface {
	setState(state string)
}

// requestBody returns the result of the given function, if a request body is needed (Create and Update operations).
func requestBody(ctx context.Context, br func() interface{}) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate || op == types.OperationUpdate {
		response := br()

		if op == types.OperationUpdate {
			if crbi, ok := response.(commonRequestBodyInterface); ok {
				crbi.setState(Updating.ID)
			}
		}

		return response, nil
	}

	return nil, nil
}
