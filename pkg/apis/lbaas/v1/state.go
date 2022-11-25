package v1

import (
	"go.anx.io/go-anxcloud/pkg/apis/common/gs"
)

var (
	Updating        = gs.State{ID: "0", Text: "Updating", Type: gs.StateTypeOK}
	Updated         = gs.State{ID: "1", Text: "Updated", Type: gs.StateTypeOK}
	DeploymentError = gs.State{ID: "2", Text: "DeploymentError", Type: gs.StateTypeError}
	Deployed        = gs.State{ID: "3", Text: "Deployed", Type: gs.StateTypeOK}
	NewlyCreated    = gs.State{ID: "4", Text: "NewlyCreated", Type: gs.StateTypePending}
)

// Deprecated: use gs.State instead
type State gs.State
