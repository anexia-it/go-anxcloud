package mock

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// Hook defines actions that can be performed on a given object.
// They can run pre-creation or pre-update to mock certain domain-specific behaviour
// of the Anexia Engine.
type Hook func(ctx context.Context, a API, o types.IdentifiedObject)

type hookName string

const (
	preCreateHook hookName = "pre-create"
	preUpdateHook hookName = "pre-update"
)

// WithPreCreateHook adds the given hook function to the API.
//
// Pre-create hooks are invoked in the order they were added *immediately* upon invocation of Create.
func WithPreCreateHook(h Hook) APIOption {
	return func(a *mockAPI) {
		a.hooks[preCreateHook] = append(a.hooks[preCreateHook], h)
	}
}

// WithPreUpdateHook adds the given hook function to the API.
//
// Pre-update hooks are invoked in the order they were added *immediately* upon invocation of Update.
func WithPreUpdateHook(h Hook) APIOption {
	return func(a *mockAPI) {
		a.hooks[preUpdateHook] = append(a.hooks[preUpdateHook], h)
	}
}
