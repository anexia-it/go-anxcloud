package types

import (
	"errors"
	"fmt"
)

var (
	// ErrUnidentifiedObject is returned when an IdentifiedObject was required, but the passed object didn't have the identifying attribute set.
	ErrUnidentifiedObject = errors.New("passed object does not have its identifying attribute set")

	// ErrInvalidFilter is returned when the configured filters cannot be applied to a List operation.
	ErrInvalidFilter = errors.New("invalid filters configured")

	// ErrInvalidFilterCombination is returned when the configured filters cannot be combined in a single
	// List operation. It wraps ErrInvalidFilter.
	ErrInvalidFilterCombination = fmt.Errorf("%w: combination of filters", ErrInvalidFilter)

	// ErrTypeNotSupported is returned when an argument is of type interface{}, manual type checking via reflection is done and the given arguments type cannot be used.
	ErrTypeNotSupported = errors.New("the given type cannot be used for the requested operation")

	// ErrObjectWithoutIdentifier is a specialized ErrTypeNotSupport for Objects not having a fields tagged with `anxcloud:"identifier"`.
	ErrObjectWithoutIdentifier = fmt.Errorf("%w: Object lacks identifier field", ErrTypeNotSupported)

	// ErrObjectWithMultipleIdentifier is a specialized ErrTypeNotSupport for Objects having multiple fields tagged with `anxcloud:"identifier"`.
	ErrObjectWithMultipleIdentifier = fmt.Errorf("%w: Object has multiple fields tagged as identifier", ErrTypeNotSupported)

	// ErrObjectIdentifierTypeNotSupported is a specialized ErrTypeNotSupport for Objects having a field tagged with `anxcloud:"identifier"` with an unsupported type.
	ErrObjectIdentifierTypeNotSupported = fmt.Errorf("%w: Objects identifier field has an unsupported type", ErrTypeNotSupported)
)
