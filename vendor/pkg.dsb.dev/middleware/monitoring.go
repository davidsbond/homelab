package middleware

import (
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
)

// Monitoring is an HTTP middleware that will write panics to the application
// monitoring provider.
func Monitoring() mux.MiddlewareFunc {
	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle
}
