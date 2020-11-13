// Package tracing contains orchestration code for opentracing. Specifically, it
// is used to build an integration with Jaeger.
package tracing

import (
	"context"
	"errors"
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
	"pkg.dsb.dev/logging"
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
			jaeger.ReporterOptions.Logger(logging.JaegerLogger()),
		),
		jaeger.TracerOptions.Logger(logging.JaegerLogger()),
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

// SetSpanTag adds the given key/value pair to the tags for the current span.
func SetSpanTag(ctx context.Context, key string, val interface{}) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	span.SetTag(key, val)
}

// SpanMetadata returns a map[string]string containing the metadata for the provided span. This should be used
// for passing span information across application boundaries.
func SpanMetadata(span opentracing.Span) (map[string]string, error) {
	if span == nil {
		return nil, nil
	}

	carrier := opentracing.TextMapCarrier{}
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	return carrier, err
}

// SpanFromMetadata starts a new span whose name is set to operationName. The previous span is extracted from the given
// map[string]string. If no span metadata is set, a new span is started using whatever is inside the given context.
func SpanFromMetadata(ctx context.Context, operationName string, md map[string]string) (opentracing.Span, context.Context, error) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	switch {
	case errors.Is(err, opentracing.ErrSpanContextNotFound):
		span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
		return span, ctx, nil
	case err != nil:
		return nil, nil, err
	default:
		span := opentracing.StartSpan(operationName, opentracing.ChildOf(spanCtx))
		ctx := opentracing.ContextWithSpan(ctx, span)
		return span, ctx, nil
	}
}
