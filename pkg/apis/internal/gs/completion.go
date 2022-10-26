package gs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.anx.io/go-anxcloud/pkg/api/types"
)

var (
	// ErrStateError is returned if a resource could not be provisioned (state "Error")
	ErrStateError = errors.New(`resource has state "Error"`)

	// ErrStateUnknown is returned if a resource has an unknown state
	ErrStateUnknown = errors.New("resource has an unknown state")
)

const awaitCompletionPollInterval = 30 * time.Second

// AwaitCompletion blocks until an object is no longer pending
func AwaitCompletion(ctx context.Context, a types.API, o objectWithStateRetriever) error {
	ticker := time.NewTicker(awaitCompletionPollInterval)
	defer ticker.Stop()

	for {
		if err := a.Get(ctx, o); err != nil {
			return fmt.Errorf("failed to get cluster: %w", err)
		}

		if o.StateSuccess() {
			return nil
		} else if o.StateFailure() {
			return ErrStateError
		} else if !o.StateProgressing() {
			return ErrStateUnknown
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			continue
		}
	}
}
