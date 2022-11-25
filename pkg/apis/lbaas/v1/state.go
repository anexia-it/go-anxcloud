package v1

import (
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

var (
	Updating        = State{ID: "0", Text: "Updating", Type: gs.StateTypeOK}
	Updated         = State{ID: "1", Text: "Updated", Type: gs.StateTypeOK}
	DeploymentError = State{ID: "2", Text: "DeploymentError", Type: gs.StateTypeError}
	Deployed        = State{ID: "3", Text: "Deployed", Type: gs.StateTypeOK}
	NewlyCreated    = State{ID: "4", Text: "NewlyCreated", Type: gs.StateTypePending}
)

type State gs.State
