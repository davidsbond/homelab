// Package requestid is used to add/extract a request identifier to/from a
// context.Context.
package requestid

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
)

type (
	ctxKey struct{}
)

// Extract the request id from the X-Request-ID header from 'in'. If it does not exist, it is
// generated as a V4 UUID. The X-Request-ID is then set on 'out'. The returned context contains
// the request id and if tracing is enabled the request id is added to the current span.
func Extract(ctx context.Context, in, out http.Header) context.Context {
	const key = "X-Request-ID"

	id := in.Get(key)
	if id == "" {
		id = uuid.Must(uuid.NewV4()).String()
		in.Set(key, id)
	}

	out.Set(key, id)
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("http.request_id", id)
	}

	return ToContext(ctx, id)
}

// ToContext adds a request id to a context.
func ToContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// FromContext obtains a request id from a context.
func FromContext(ctx context.Context) string {
	id, ok := ctx.Value(ctxKey{}).(string)
	if !ok {
		return ""
	}
	return id
}
