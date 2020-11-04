// Package retry contains utilities for performing retries when functions return errors.
package retry

import (
	"context"
	"errors"

	"github.com/cenkalti/backoff"
)

// Do performs an action with a maximum number of retries. If fn returns a non-nil error, it will
// be retried up to the maximum number of times and then returned. If fn returns a nil error, the retry
// loop is broken. The loop is also broken if the provided context is cancelled. Returning an error
// wrapped with the Stop function will prevent further retries from occurring.
func Do(ctx context.Context, max int, fn func(ctx context.Context) error) error {
	var stopErr *stopper
	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(max))

	return backoff.Retry(func() error {
		if err := ctx.Err(); err != nil {
			return backoff.Permanent(err)
		}

		err := fn(ctx)
		if errors.As(err, &stopErr) {
			return backoff.Permanent(stopErr.err)
		}

		return err
	}, bo)
}

type (
	stopper struct {
		err error
	}
)

// Stop wraps err so that when it is returned to the retry mechanism no further retries
// are attempted.
func Stop(err error) error {
	return &stopper{err: err}
}

func (s *stopper) Error() string {
	return s.err.Error()
}
