// Package server contains utilities for managing different kinds of servers. Currently supports HTTP & gRPC.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"pkg.dsb.dev/logging"
	"pkg.dsb.dev/middleware"
)

type (
	httpConfig struct {
		handler         *mux.Router
		securityEnabled bool
		tlsCache        string
		tlsHostPolicy   autocert.HostPolicy
	}
)

var defaultHTTPConfig = httpConfig{
	handler:  mux.NewRouter(),
	tlsCache: "/cache",
	tlsHostPolicy: func(_ context.Context, _ string) error {
		return nil
	},
}

// ServeHTTP starts the HTTP server. When the provided context is closed, the server is shut down.
func ServeHTTP(ctx context.Context, opts ...HTTPOption) error {
	c := defaultHTTPConfig
	for _, opt := range opts {
		opt(&c)
	}

	c.handler.Use(
		middleware.Security(c.securityEnabled),
		middleware.Metrics(),
		middleware.Tracing(),
		middleware.RequestID(),
		middleware.Panic(),
		middleware.Monitoring(),
	)

	svr := &http.Server{
		Addr:     ":80",
		Handler:  middleware.CORS(c.handler),
		ErrorLog: log.New(ioutil.Discard, "", log.Flags()),
	}

	grp, ctx := errgroup.WithContext(ctx)
	if c.securityEnabled {
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(c.tlsCache),
			HostPolicy: c.tlsHostPolicy,
		}

		svr.Addr = ":443"
		svr.TLSConfig = createTLSConfig(m)
		svr.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))

		csvr := &http.Server{
			Addr:     ":80",
			Handler:  m.HTTPHandler(http.NewServeMux()),
			ErrorLog: log.New(ioutil.Discard, "", log.Flags()),
		}
		grp.Go(csvr.ListenAndServe)
		grp.Go(func() error {
			return waitForHTTPShutdown(ctx, csvr)
		})
	}

	logging.WithField("port", svr.Addr).Info("serving http")
	grp.Go(func() error {
		if c.securityEnabled {
			return svr.ListenAndServeTLS("", "")
		}

		return svr.ListenAndServe()
	})
	grp.Go(func() error {
		return waitForHTTPShutdown(ctx, svr)
	})

	err := grp.Wait()
	switch {
	case errors.Is(err, http.ErrServerClosed):
		logging.Info("server shut down")
		return nil
	case err != nil:
		return err
	default:
		return nil
	}
}

func waitForHTTPShutdown(ctx context.Context, svr *http.Server) error {
	<-ctx.Done()
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return svr.Shutdown(ctx)
}

func createTLSConfig(m *autocert.Manager) *tls.Config {
	return &tls.Config{
		GetCertificate:           m.GetCertificate,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
