// Package cron contains methods for running actions as cron jobs.
package cron

import (
	"context"
	"time"
)

type (
	// The Action type is a function that is invoked by a cron.
	Action func(context.Context) error
)

// Every invokes 'fn' every time the 'freq' duration passes. This will continue until the provided context
// is cancelled, or 'fn' returns an error.
func Every(ctx context.Context, freq time.Duration, fn Action) error {
	ticker := time.NewTicker(freq)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := fn(ctx); err != nil {
				return err
			}
		}
	}
}
