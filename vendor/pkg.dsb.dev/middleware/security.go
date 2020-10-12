package middleware

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
)

// Security is an HTTP middleware that applies HTTP security features to inbound
// requests. If 'enabled' is false, SSL features are disabled.
func Security(enabled bool) mux.MiddlewareFunc {
	const stsSeconds = 86400

	mw := secure.New(secure.Options{
		BrowserXssFilter:     true,
		ContentTypeNosniff:   true,
		FrameDeny:            true,
		IsDevelopment:        !enabled,
		SSLRedirect:          true,
		SSLTemporaryRedirect: false,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		HostsProxyHeaders:    []string{"X-Forwarded-Host"},
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           stsSeconds,
	})

	return mw.Handler
}
