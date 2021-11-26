// Package types contains everything needed by APIs implementing the interfaces to be compatible with the
// generic API client.
package types

// ObjectRetriever retrieves an object, parsing it to the correct go type.
type ObjectRetriever func(Object) error

// ObjectChannel streams objects.
type ObjectChannel <-chan ObjectRetriever
