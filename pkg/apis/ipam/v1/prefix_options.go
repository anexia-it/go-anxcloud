package v1

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

const (
	optionKeyCreateEmpty = "ipamv1/prefix/createEmpty"
)

// CreateEmpty can be used to define if a Prefix is to be created with all Address objects
// in it created (false) or only Address objects created that are actually in use (true).
func CreateEmpty(empty bool) types.CreateOption {
	return createEmptyOption(empty)
}

type createEmptyOption bool

func (ceo createEmptyOption) ApplyToCreate(opts *types.CreateOptions) {
	// It can return an error when the requested key is already set, but overwriting is disabled.
	// Since we have overwriting enabled, we can ignore the error.
	_ = opts.Set(optionKeyCreateEmpty, bool(ceo), true)
}
