package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"pkg.dsb.dev/requestid"
)

// RequestID is a middleware that reuses or creates a request id for each HTTP request. The
// id is also added to the request context using the requestid package.
func RequestID() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := requestid.Extract(r.Context(), r.Header, w.Header())
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
