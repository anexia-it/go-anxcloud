package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
)

// Retry helper to ease with retrying tasks using the passed callback function.
func Retry(ctx context.Context, count int, sleep time.Duration, cb func() (bool, error)) error {
	var (
		err       error
		retryable bool
	)

	for i := 0; i < count; i++ {
		if i > 0 {
			logr.FromContextOrDiscard(ctx).Info(fmt.Sprintf("retry callback in a second, due to an error: %s", err))
			time.Sleep(sleep)
		}

		if retryable, err = cb(); err == nil {
			return nil
		} else if !retryable {
			break
		}
	}

	return err
}
