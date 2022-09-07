package api

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

// TODO(LittleFox94): Maybe Client is a better name for this, but
// we'd then have to rename client to transport

// API is the interface to perform operations on the engine.
type API types.API
