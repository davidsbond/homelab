// Package environment contains utilities for interacting with the application's environment.
package environment

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"pkg.dsb.dev/logging"

	"go.uber.org/automaxprocs/maxprocs"
)

var (
	// ApplicationName is the name of the application, injected at build time.
	ApplicationName string

	// ApplicationDescription is the description of the application, injected at build time.
	ApplicationDescription string

	// Version is the version of the application, injected at build time.
	Version string

	// The timestamp the application was compiled, injected at build time.
	compiled string
)

// Compiled returns the time.Time representation of the application's
// build timestamp.
func Compiled() time.Time {
	unix, err := strconv.ParseInt(compiled, 10, 64)
	if err != nil {
		logging.WithError(err).Warn("invalid compile time, using now")
		return time.Now()
	}

	return time.Unix(unix, 0)
}

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

// SetMaxProcsToCPUQuota sets GOMAXPROCS to match the Linux container CPU quota (if any).
func SetMaxProcsToCPUQuota() {
	if autoMaxProcsDisabled {
		return
	}

	if _, err := maxprocs.Set(); err != nil {
		logging.WithError(err).Warn("failed to set GOMAXPROCS")
	}
}
