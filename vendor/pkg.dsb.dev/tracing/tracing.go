// Package tracing contains orchestration code for opentracing. Specifically, it
// is used to build an integration with Jaeger.
package tracing

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/environment"
)

type (
	options struct {
		disabled            bool
		host                string
		sampleRate          float64
		bufferFlushInterval time.Duration
	}
)

var config = &options{
	host:                "jaeger:6831",
	sampleRate:          1,
	bufferFlushInterval: time.Second,
}

// New sets up the global tracer. The io.Closer instance returned is
// used to stop tracing.
func New() (io.Closer, error) {
	if config.disabled {
		return closers.Noop, nil
	}

	sender, err := jaeger.NewUDPTransport(config.host, 0)
	if err != nil {
		return nil, err
	}

	tracer, closer := jaeger.NewTracer(
		environment.ApplicationName,
		jaeger.NewRateLimitingSampler(config.sampleRate),
		jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(config.bufferFlushInterval),
		),
	)

	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

// WithError returns err after adding its details to the span.
func WithError(span opentracing.Span, err error) error {
	if err == nil {
		return nil
	}

	ext.Error.Set(span, true)
	span.LogFields(otlog.Error(err))
	return err
}

// AddError adds the provided error to the span.
func AddError(span opentracing.Span, err error) {
	if err == nil {
		return
	}

	ext.Error.Set(span, true)
	span.LogFields(otlog.Error(err))
}

// WrapHTTPHandler wraps an http.Handler so that all inbound requests create a
// trace. Ignores operational endpoints.
func WrapHTTPHandler(h http.Handler) http.Handler {
	if config.disabled {
		return h
	}

	return nethttp.Middleware(opentracing.GlobalTracer(), h,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			// Use request URI as span name
			return r.RequestURI
		}), nethttp.MWSpanFilter(func(r *http.Request) bool {
			// Ignore operational endpoints.
			return !strings.HasPrefix(r.RequestURI, "/__/")
		}))
}
