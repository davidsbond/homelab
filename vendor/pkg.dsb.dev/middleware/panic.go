package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"pkg.dsb.dev/transport"
)

// Panic is a middleware function that handles runtime panics spawning from HTTP requests. If a panic
// occurs, the goroutine is recovered and an internal server error is written to the client.
func Panic() mux.MiddlewareFunc {
	t := transport.HTTP{}

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				p := recover()
				if p == nil {
					return
				}

				t.ErrorWithStack(r.Context(), w, http.StatusInternalServerError, "%v", p)
			}()

			handler.ServeHTTP(w, r)
		})
	}
}
