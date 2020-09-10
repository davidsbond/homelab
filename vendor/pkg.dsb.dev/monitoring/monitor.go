// Package monitoring contains helpers for application monitoring. Currently configured
// to use sentry.
package monitoring

import (
	"context"
	"io"
	"time"

	"github.com/getsentry/sentry-go"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/environment"
)

type (
	// The Monitor type is responsible for maintaining the lifetime of the connection to the
	// application monitoring provider. It implements io.Closer and should be closed on application
	// exit.
	Monitor struct {
		flushTimeout time.Duration
	}
)

// New creates a new instance of the Monitor type and sets up the connection to the monitoring
// provider.
func New() (io.Closer, error) {
	if config.disabled {
		return closers.Noop, nil
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         config.dsn,
		ServerName:  environment.ApplicationName,
		Release:     environment.Version,
		Environment: config.environment,
	}); err != nil {
		return nil, err
	}

	return &Monitor{
		flushTimeout: config.flushTimeout,
	}, nil
}

// Close the connection to the monitoring provider.
func (m *Monitor) Close() error {
	if sentry.Flush(m.flushTimeout) {
		return nil
	}

	return context.DeadlineExceeded
}

// WithError reports an error to the provider.
func WithError(err error) error {
	if err == nil {
		return nil
	}

	sentry.CaptureException(err)
	return err
}
