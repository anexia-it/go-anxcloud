// Package types contains everything needed by APIs implementing the interfaces to be compatible with the
// generic API client.
package types

import "context"

// TODO(LittleFox94): Maybe Client is a better name for this, but
// we'd then have to rename client to transport

// API is the interface to perform operations on the engine.
type API interface {
	// Get the identified object from the engine. Set the identifying attribute on the
	// object passed to this function.
	Get(context.Context, IdentifiedObject, ...GetOption) error

	// Create the given object on the engine.
	Create(context.Context, Object, ...CreateOption) error

	// Update the object on the engine.
	Update(context.Context, IdentifiedObject, ...UpdateOption) error

	// Destroy the identified object.
	Destroy(context.Context, IdentifiedObject, ...DestroyOption) error

	// List objects matching the info given in the object.
	// Beware: listing endpoints usually do not return all data for an object, sometimes
	// only the identifier is filled. This varies by specific API. If you need full objects,
	// the FullObjects option might be your friend.
	List(context.Context, FilterObject, ...ListOption) error
}

// ObjectRetriever retrieves an object, parsing it to the correct go type.
type ObjectRetriever func(Object) error

// ObjectChannel streams objects.
type ObjectChannel <-chan ObjectRetriever
