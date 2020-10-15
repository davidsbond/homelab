// Package retry contains utilities for performing retries when functions return errors.
package retry

import "context"

// Do performs an action with a maximum number of retries. If fn returns a non-nil error, it will
// be retried up to the maximum number of times and then returned. If fn returns a nil error, the retry
// loop is broken. The loop is also broken if the provided context is cancelled.
func Do(ctx context.Context, max int, fn func(ctx context.Context) error) error {
	var err error

retry:
	for i := 0; i < max; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = fn(ctx)
			if err == nil {
				break retry
			}
		}
	}

	return err
}
