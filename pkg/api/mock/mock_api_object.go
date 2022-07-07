package mock

import (
	"time"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

// APIObject wraps types.Object with mock specific metadata
// Do not use this directly, but rather the custom Gomega matchers
// located in the 'matcher' sub-package
type APIObject struct {
	wrapped types.Object
	tags    map[string]interface{}

	existing bool

	createdCount   int
	updatedCount   int
	destroyedCount int

	createdTime   time.Time
	updatedTime   time.Time
	destroyedTime time.Time
}

// Unwrap returns the wrapped types.Object
func (o *APIObject) Unwrap() types.Object {
	return o.wrapped
}

// Tags returns list of tags the object was tagged with
func (o *APIObject) Tags() []string {
	out := make([]string, 0, len(o.tags))
	for tag := range o.tags {
		out = append(out, tag)
	}
	return out
}

// HasTags checks if APIObject has provided tags
func (o *APIObject) HasTags(tags ...string) bool {
	for _, tag := range tags {
		if _, ok := o.tags[tag]; !ok {
			return false
		}
	}
	return true
}

// Existing returns whether or not an object currently exists
func (o *APIObject) Existing() bool {
	return o.existing
}

// CreatedCount returns how often an object was created
func (o *APIObject) CreatedCount() int {
	return o.createdCount
}

// UpdatedCount returns how often an object was updated
func (o *APIObject) UpdatedCount() int {
	return o.updatedCount
}

// DestroyedCount returns how often an object was destroyed
func (o *APIObject) DestroyedCount() int {
	return o.destroyedCount
}
