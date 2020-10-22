// Package closers contains utilities for dealing with io.Closer implementations
package closers

import (
	"io"
	"log"
)

type (
	noopCloser struct{}

	functionCloser struct {
		fn func() error
	}
)

// Noop is an io.Closer implementation that does nothing.
var Noop io.Closer = &noopCloser{}

// Close the given closer. If a non-nil error is returned, it is logged.
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err.Error())
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
