package mock

import (
	"context"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

type hook func(ctx context.Context, a API, o types.IdentifiedObject)

type hookName string

const (
	preCreateHook hookName = "pre-create"
)

func WithPreCreateHook(h hook) APIOption {
	return func(a *mockAPI) {
		a.hooks[preCreateHook] = append(a.hooks[preCreateHook], h)
	}
}
