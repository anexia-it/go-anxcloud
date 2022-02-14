package types

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidFilter is returned when the configured filters cannot be applied to a List operation.
	ErrInvalidFilter = errors.New("invalid filters configured")

	// ErrInvalidFilterCombination is returned when the configured filters cannot be combined in a single
	// List operation. It wraps ErrInvalidFilter.
	ErrInvalidFilterCombination = fmt.Errorf("%w: combination of filters", ErrInvalidFilter)
)
