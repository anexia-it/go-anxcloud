package api

import (
	"context"
	"net/url"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"
)

type contextKey string

const (
	contextKeyOperation contextKey = "operation"
	contextKeyOptions   contextKey = "options"
	contextKeyURL       contextKey = "url"
)

func contextWithOperation(ctx context.Context, op types.Operation) context.Context {
	return context.WithValue(ctx, contextKeyOperation, op)
}

func contextWithOptions(ctx context.Context, opts types.Options) context.Context {
	return context.WithValue(ctx, contextKeyOptions, opts)
}

func contextWithURL(ctx context.Context, url url.URL) context.Context {
	return context.WithValue(ctx, contextKeyURL, url)
}

// OperationFromContext returns the current generic client operation from a given context.
// This is set on every context passed by the generic client to Object functions.
func OperationFromContext(ctx context.Context) (types.Operation, error) {
	if op, ok := ctx.Value(contextKeyOperation).(types.Operation); ok {
		return op, nil
	}

	return "", ErrContextKeyNotSet
}

// OptionsFromContext returns the Options for the current generic client operation.
// This is set on every context passed by the generic client to Object functions.
func OptionsFromContext(ctx context.Context) (types.Options, error) {
	if op, ok := ctx.Value(contextKeyOptions).(types.Options); ok {
		return op, nil
	}

	return nil, ErrContextKeyNotSet
}

// URLFromContext returns the url originally returned from the Objects EndpointURL method for
// the current generic client operation.
// This is set on every context passed by the generic client to Object functions _after_ the
// call to EndpointURL.
func URLFromContext(ctx context.Context) (url.URL, error) {
	if op, ok := ctx.Value(contextKeyOperation).(url.URL); ok {
		return op, nil
	}

	return url.URL{}, ErrContextKeyNotSet
}
