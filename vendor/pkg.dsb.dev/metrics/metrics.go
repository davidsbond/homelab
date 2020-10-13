// Package metrics contains helpers for exposing operational metrics for the
// application.
package metrics

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"

	"pkg.dsb.dev/environment"
	"pkg.dsb.dev/logging"
)

var disabled bool

// Serve prometheus metrics via HTTP on the provided router.
func Serve(r *mux.Router) {
	if disabled {
		return
	}

	r.Handle("/__/metrics", promhttp.Handler())
}

// Register a prometheus collector.
func Register(collectors ...prometheus.Collector) {
	prometheus.DefaultRegisterer.MustRegister(collectors...)
}

// AddSQLStats exports prometheus metrics for the provided SQL database
// connection.
func AddSQLStats(db *sql.DB) {
	if disabled {
		return
	}

	if err := prometheus.Register(newSQLStatsCollector(db)); err != nil {
		logging.WithError(err).Error("failed to register metrics")
	}
}

var (
	pushURL      string
	errNoPushURL = errors.New("push gateway url not configured")
)

// Push all registered collectors to the configured push gateway.
func Push() error {
	if disabled {
		return nil
	}

	if pushURL == "" {
		return errNoPushURL
	}

	pc := push.New(pushURL, environment.ApplicationName).Gatherer(prometheus.DefaultGatherer)

	return pc.Push()
}
