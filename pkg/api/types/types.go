// Package types contains everything needed by APIs implementing the interfaces to be compatible with the
// generic API client.
package types

// ObjectReturner retrieves an object, parsing it to the correct go type.
type ObjectReturner func(Object) error

// ObjectChannel streams objects.
type ObjectChannel chan ObjectReturner
