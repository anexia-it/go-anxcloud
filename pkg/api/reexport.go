package api

import "go.anx.io/go-anxcloud/pkg/api/types"

// We re-export them here to group options given by this package under their options in the docs.

// ListOption is the interface options have to implement to be usable with List operation. Re-exported from pkg/api/types.
type ListOption = types.ListOption

// GetOption is the interface options have to implement to be usable with Get operation. Re-exported from pkg/api/types.
type GetOption = types.GetOption

// CreateOption is the interface options have to implement to be usable with Create operation. Re-exported from pkg/api/types.
type CreateOption = types.CreateOption

// UpdateOption is the interface options have to implement to be usable with Update operation. Re-exported from pkg/api/types.
type UpdateOption = types.UpdateOption

// DestroyOption is the interface options have to implement to be usable with Destroy operation. Re-exported from pkg/api/types.
type DestroyOption = types.DestroyOption
