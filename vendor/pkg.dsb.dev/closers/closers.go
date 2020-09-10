// Package closers contains utilities for dealing with io.Closer implementations
package closers

import (
	"io"
	"log"
)

type (
	noopCloser struct{}
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
