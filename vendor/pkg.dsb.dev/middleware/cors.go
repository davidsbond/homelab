package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

var (
	allowedHeaders = []string{
		"Authorization",
		"X-Request-ID",
		"Origin",
		"Content-Type",
		"Accept",
		"Accept-Language",
		"Content-Language",
		"Range",
	}

	allowedMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}

	allowedOrigins = []string{
		"http://localhost:9080",
		"https://localhost:9080",
	}
)

// CORS is a middleware function that manages cross-origin requests.
func CORS(h http.Handler) http.Handler {
	mw := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders(allowedHeaders),
		handlers.AllowedMethods(allowedMethods),
		handlers.AllowedOrigins(allowedOrigins),
	)

	return mw(h)
}
