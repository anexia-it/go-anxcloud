package filter

import (
	"errors"
	"net/url"
)

// ErrUnknownField is returned when trying to retrieve a filter value that is not configured as filterable.
var ErrUnknownField = errors.New("unknown field for filtering")

// Helper allows easy access to filter parameters and building the query parameters.
type Helper interface {
	// BuildQuery returns the query parameters to set for filtering.
	BuildQuery() url.Values

	// Get returns the value and if it was set for a given named field.
	Get(field string) (value interface{}, ok bool, err error)
}
