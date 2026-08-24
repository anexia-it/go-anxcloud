package v1

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// commonRequestBody is embedded in the request body types of LBaaS objects.
//
// It shadows the "state" JSON key of the embedded LBaaS object (e.g. Backend.HasState.State), because
// encoding/json resolves same-named fields at different embedding depths by picking the shallower one.
// This is the only way to control what is sent as "state", so both Create and Update need to explicitly
// route the desired state through setState.
type commonRequestBody struct {
	State string `json:"state,omitempty"`
}

func (crb *commonRequestBody) setState(state string) {
	crb.State = state
}

type commonRequestBodyInterface interface {
	setState(state string)
}

// requestBody returns the result of the given function, if a request body is needed (Create and Update operations).
//
// For Update operations the state is always forced to Updating, telling the Engine to reapply the
// configuration. For Create operations, createState is sent as-is: pass the object's own desired
// state (e.g. Backend.State.ID) to request a specific initial state (such as Managed), or an empty
// string to keep the previous behaviour of not sending "state" at all and letting the Engine default it.
func requestBody(ctx context.Context, br func() interface{}, createState string) (interface{}, error) {
	op, err := types.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == types.OperationCreate || op == types.OperationUpdate {
		response := br()

		if crbi, ok := response.(commonRequestBodyInterface); ok {
			switch op {
			case types.OperationCreate:
				crbi.setState(createState)
			case types.OperationUpdate:
				crbi.setState(Updating.ID)
			}
		}

		return response, nil
	}

	return nil, nil
}
