// Package cron contains methods for running actions as cron jobs.
package cron

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"

	"pkg.dsb.dev/retry"
	"pkg.dsb.dev/tracing"
)

type (
	// The Action type is a function that is invoked by a cron.
	Action func(context.Context) error
)

const retryMax = 3

// Every invokes 'fn' every time the 'freq' duration passes. This will continue until the provided context
// is cancelled, or 'fn' returns an error.
func Every(ctx context.Context, freq time.Duration, fn Action) error {
	ticker := time.NewTicker(freq)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			span, ctx := opentracing.StartSpanFromContext(ctx, "cron-run")
			span.SetTag("cron.frequency", freq)
			if err := retry.Do(ctx, retryMax, fn); err != nil {
				span.Finish()
				return tracing.WithError(span, err)
			}

			span.Finish()
		}
	}
}

// At executes fn at the desired time. Aims to have an accuracy of within a second.
func At(ctx context.Context, at time.Time, fn Action) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ts := <-ticker.C:
			dur := time.Date(
				at.Year(), at.Month(), at.Day(), ts.Hour(),
				ts.Minute(), ts.Second(), ts.Nanosecond(), ts.Location(),
			).Sub(at)

			if dur >= time.Second || dur <= -time.Second {
				continue
			}

			span, ctx := opentracing.StartSpanFromContext(ctx, "cron-run")
			span.SetTag("cron.at", at)
			if err := retry.Do(ctx, retryMax, fn); err != nil {
				span.Finish()
				return tracing.WithError(span, err)
			}

			span.Finish()
		}
	}
}
