package matcher

import "errors"

var (
	// ErrMatcherExpectsAPIObjectPointer is returned when matcher is called with an actual other than *APIObject
	ErrMatcherExpectsAPIObjectPointer = errors.New("matcher expects an *APIObject")
)
