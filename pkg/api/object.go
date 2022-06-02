package api

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

var (
	// ErrUnidentifiedObject is returned when an IdentifiedObject was required, but the passed object didn't have the identifying attribute set.
	//
	// Deprecated: moved to pkg/api/types
	ErrUnidentifiedObject = types.ErrUnidentifiedObject

	// ErrTypeNotSupported is returned when an argument is of type interface{}, manual type checking via reflection is done and the given arguments type cannot be used.
	//
	// Deprecated: moved to pkg/api/types
	ErrTypeNotSupported = types.ErrTypeNotSupported

	// ErrObjectWithoutIdentifier is a specialized ErrTypeNotSupport for Objects not having a fields tagged with `anxcloud:"identifier"`.
	//
	// Deprecated: moved to pkg/api/types
	ErrObjectWithoutIdentifier = types.ErrObjectWithoutIdentifier

	// ErrObjectWithMultipleIdentifier is a specialized ErrTypeNotSupport for Objects having multiple fields tagged with `anxcloud:"identifier"`.
	//
	// Deprecated: moved to pkg/api/types
	ErrObjectWithMultipleIdentifier = types.ErrObjectWithMultipleIdentifier

	// ErrObjectIdentifierTypeNotSupported is a specialized ErrTypeNotSupport for Objects having a field tagged with `anxcloud:"identifier"` with an unsupported type.
	//
	// Deprecated: moved to pkg/api/types
	ErrObjectIdentifierTypeNotSupported = types.ErrObjectIdentifierTypeNotSupported
)

// GetObjectIdentifier extracts the identifier of the given object, returning an error if no identifier field
// is found or singleObjectOperation is true and an identifier field is found, but empty.
//
// Deprecated: moved to pkg/api/types
func GetObjectIdentifier(obj types.Object, singleObjectOperation bool) (string, error) {
	return types.GetObjectIdentifier(obj, singleObjectOperation)
}
