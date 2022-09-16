package api

import (
	"go.anx.io/go-anxcloud/pkg/api/internal"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

// ObjectChannel configures the List operation to return the objects via the given channel. When listing via
// channel you either have to read until the channel is closed or pass a context you cancel explicitly - failing
// to do that will result in leaked goroutines.
func ObjectChannel(channel *types.ObjectChannel) ListOption {
	return internal.ObjectChannelOption{Channel: channel}
}

// Paged is an option valid for List operations to retrieve objects in a paged fashion (instead of all at once).
func Paged(page, limit uint, info *types.PageInfo) ListOption {
	return internal.PagedOption{
		Page:  page,
		Limit: limit,
		Info:  info,
	}
}

// FullObjects can be set to make a Get for every object before it is returned to the caller of List(). This
// is necessary since most API endpoints for listing objects only return a subset of their data.
//
// Beware: this makes one API call to retrieve the objects (ok, one call per page of objects) and an additional
// call per object. Because of this being very slow, it is an optional feature and should only be used with care.
func FullObjects(fullObjects bool) ListOption {
	return internal.FullObjectsOption(fullObjects)
}

// AutoTag can be used to automatically tag objects after creation
func AutoTag(tags ...string) CreateOption {
	return internal.AutoTagOption(tags)
}
