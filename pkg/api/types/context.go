package types

import (
	"context"
	"errors"
	"net/url"
)

type contextKey string

const (
	contextKeyOperation contextKey = "operation"
	contextKeyOptions   contextKey = "options"
	contextKeyURL       contextKey = "url"
)

// ErrContextKeyNotSet is returned when trying to retrieve an unset value from a context.
var ErrContextKeyNotSet = errors.New("requested context key is not set")

// ContextWithOperation returns a new context based on the given one with the operation added to it.
func ContextWithOperation(ctx context.Context, op Operation) context.Context {
	return context.WithValue(ctx, contextKeyOperation, op)
}

// ContextWithOptions returns a new context based on the given one with the options added to it.
func ContextWithOptions(ctx context.Context, opts Options) context.Context {
	return context.WithValue(ctx, contextKeyOptions, opts)
}

// ContextWithURL returns a new context based on the given one with the URL added to it.
func ContextWithURL(ctx context.Context, url url.URL) context.Context {
	return context.WithValue(ctx, contextKeyURL, url)
}

// OperationFromContext returns the current generic client operation from a given context.
// This is set on every context passed by the generic client to Object functions.
func OperationFromContext(ctx context.Context) (Operation, error) {
	if op, ok := ctx.Value(contextKeyOperation).(Operation); ok {
		return op, nil
	}

	return "", ErrContextKeyNotSet
}

// OptionsFromContext returns the Options for the current generic client operation.
// This is set on every context passed by the generic client to Object functions.
func OptionsFromContext(ctx context.Context) (Options, error) {
	if op, ok := ctx.Value(contextKeyOptions).(Options); ok {
		return op, nil
	}

	return nil, ErrContextKeyNotSet
}

// URLFromContext returns the url originally returned from the Objects EndpointURL method for
// the current generic client operation.
// This is set on every context passed by the generic client to Object functions _after_ the
// call to EndpointURL.
func URLFromContext(ctx context.Context) (url.URL, error) {
	if op, ok := ctx.Value(contextKeyURL).(url.URL); ok {
		return op, nil
	}

	return url.URL{}, ErrContextKeyNotSet
}
