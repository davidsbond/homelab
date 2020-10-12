package middleware

import (
	"pkg.dsb.dev/tracing"

	"github.com/gorilla/mux"
)

// Tracing returns a middleware function that adds opentracing to inbound HTTP requests. It will also
// add request ids and user identifiers to the span tags.
func Tracing() mux.MiddlewareFunc {
	return tracing.WrapHTTPHandler
}
