package mock

import (
	"time"

	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
)

// API interface for mock testing. Compatible with api.API
type API interface {
	api.API

	// All retrieves all objects from mock api (destroyed objects included)
	All() mockDataView
	// Existing retrieves existing objects
	Existing() mockDataView
	// CreatedAfter retrievs objects created after a given time (and optionally only existing)
	CreatedAfter(time.Time, bool) mockDataView
	// UpdatedAfter retrievs objects created after a given time (and optionally only existing)
	UpdatedAfter(time.Time, bool) mockDataView
	// DestroyedAfter retrievs objects destroyed after a given time
	DestroyedAfter(time.Time) mockDataView

	// FakeExisting adds an object with optional tags to the mock api without increasing the created counter
	FakeExisting(types.Object, ...string) string
	// Inspect retrieves an *APIObject or nil if not found
	Inspect(string) *APIObject
}
