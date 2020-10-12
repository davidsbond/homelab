package server

import (
	"context"
	"errors"

	"github.com/gorilla/mux"
)

type (
	// The HTTPOption type represents a server option that can be used to modify
	// the HTTP server configuration.
	HTTPOption func(c *httpConfig)
)

// WithHandler sets the server's handler to the one provided.
func WithHandler(r *mux.Router) HTTPOption {
	return func(c *httpConfig) {
		c.handler = r
	}
}

// WithSecurity enables/disables the server's SSL features.
func WithSecurity(flag bool) HTTPOption {
	return func(c *httpConfig) {
		c.securityEnabled = flag
	}
}

// WithTLSCache determines the location on the filesystem that TLS certificates
// will be cached.
func WithTLSCache(dir string) HTTPOption {
	return func(c *httpConfig) {
		c.tlsCache = dir
	}
}

// ErrInvalidTLSHost is the error given when an inbound TLS request does not match
// the server's host policy.
var ErrInvalidTLSHost = errors.New("invalid tls host")

// WithAllowedHosts sets hosts that are allowed to make TLS requests.
func WithAllowedHosts(hosts ...string) HTTPOption {
	return func(c *httpConfig) {
		c.tlsHostPolicy = func(ctx context.Context, host string) error {
			for _, h := range hosts {
				if h == host {
					return nil
				}
			}

			return ErrInvalidTLSHost
		}
	}
}
