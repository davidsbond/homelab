// Package closers contains utilities for dealing with io.Closer implementations
package closers

import (
	"io"

	"pkg.dsb.dev/logging"
	"pkg.dsb.dev/multierror"
)

type (
	noopCloser struct{}

	functionCloser struct {
		fn func() error
	}

	// The MultiCloser type is a collection of io.Closer implementations that can be collected
	// and closed at once.
	MultiCloser []io.Closer
)

// Noop is an io.Closer implementation that does nothing.
var Noop io.Closer = &noopCloser{}

// Close the given closer. If a non-nil error is returned, it is logged.
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		logging.WithError(err).Error("failed to close")
	}
}

// CloseAll calls Close on all io.Closer implementations passed as arguments.
func CloseAll(cs ...io.Closer) {
	for _, c := range cs {
		Close(c)
	}
}

// Close does nothing.
func (nc *noopCloser) Close() error {
	return nil
}

// CloseFunc wraps te provider function as an io.Closer implementation. This can be used when you
// have a type that has some form of close method that does not directly implement io.Closer.
func CloseFunc(fn func() error) io.Closer {
	return &functionCloser{fn: fn}
}

// Close calls the desired function.
func (c *functionCloser) Close() error {
	return c.fn()
}

// Add a closer to the MultiCloser. It will be closed when calling MultiCloser.Close.
func (c *MultiCloser) Add(cl io.Closer) {
	*c = append(*c, cl)
}

// Close all io.Closer implementations in the MultiCloser. Returns any non-nil errors as a
// multi-error.
func (c *MultiCloser) Close() error {
	errs := make([]error, len(*c))
	for i, closer := range *c {
		errs[i] = closer.Close()
	}

	return multierror.Append(nil, errs...)
}
