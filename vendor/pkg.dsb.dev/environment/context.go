// Package environment contains utilities for interacting with the application's environment.
package environment

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// NewContext returns a new context.Context implementation that will cancel when the
// current process receives an exit signal. Should be used for graceful shutdowns.
func NewContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sCh := make(chan os.Signal, 1)
		signal.Notify(sCh, os.Interrupt, syscall.SIGTERM)
		<-sCh
		cancel()
	}()

	return ctx
}
