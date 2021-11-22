package api

import (
	"github.com/anexia-it/go-anxcloud/pkg/api/internal"
	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

// ObjectChannel returns a ListOption to configure the List operation to return the objects via the ObjectChannel
// given as argument.
func ObjectChannel(ch types.ObjectChannelCloser) ListOption {
	return internal.ObjectChannelOption{Channel: ch}
}

// Paged is an option valid for List operations to retrieve objects in a paged fashion (instead of all at once).
func Paged(page, limit uint, info *types.PageInfo) ListOption {
	return internal.PagedOption{
		Page:  page,
		Limit: limit,
		Info:  info,
	}
}
