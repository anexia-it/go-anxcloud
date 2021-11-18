package api

import (
	"context"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

// TODO(LittleFox94): Maybe Client is a better name for this, but
// we'd then have to rename client to transport

// API is the interface to perform operations on the engine.
type API interface {
	// Get the identified object from the engine. Set the identifying attribute on the
	// object passed to this function.
	Get(context.Context, types.IdentifiedObject, ...GetOption) error

	// Create the given object on the engine.
	Create(context.Context, types.Object, ...CreateOption) error

	// Update the object on the engine.
	Update(context.Context, types.IdentifiedObject, ...UpdateOption) error

	// Destroy the identified object.
	Destroy(context.Context, types.IdentifiedObject, ...DestroyOption) error

	// List objects matching the info given in the object.
	// Beware: listing endpoints usually do not return all data for an object, sometimes
	// only the identifier is filled. This varies by specific API.
	List(context.Context, types.FilterObject, ...ListOption) error
}
