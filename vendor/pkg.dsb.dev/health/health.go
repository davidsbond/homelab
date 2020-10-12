// Package health contains helpers for constructing health checks.
package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/environment"
)

type (
	status struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Version     string    `json:"version"`
		Compiled    time.Time `json:"compiled"`
		Healthy     bool      `json:"healthy"`
		Checks      []*check  `json:"checks"`
	}

	check struct {
		Name    string       `json:"name"`
		Healthy bool         `json:"healthy"`
		Message string       `json:"message,omitempty"`
		Func    func() error `json:"-"`
	}
)

var (
	checks   []*check
	disabled bool
)

// RegisterGRPCHealthServer registers a gRPC health server implementation onto the
// provided gRPC server.
func RegisterGRPCHealthServer(svr *grpc.Server) *health.Server {
	hsvr := health.NewServer()
	grpc_health.RegisterHealthServer(svr, hsvr)
	return hsvr
}

// ServeGRPC periodically updates the status of the health server based on the result of the current
// health. If any components are marked as down, the status is changed to NOT_SERVING. This function
// blocks until the provided context is cancelled.
func ServeGRPC(ctx context.Context, svr *health.Server) error {
	if disabled {
		return nil
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			status := grpc_health.HealthCheckResponse_SERVING
			if err := PerformCheck(ctx); err != nil {
				status = grpc_health.HealthCheckResponse_NOT_SERVING
			}

			svr.SetServingStatus("health", status)
		}
	}
}

// ServeHTTP adds the status check endpoint to the given router.
func ServeHTTP(r *mux.Router) {
	if disabled {
		return
	}

	r.HandleFunc("/__/health", func(w http.ResponseWriter, r *http.Request) {
		h := &status{
			Name:        environment.ApplicationName,
			Description: environment.ApplicationDescription,
			Version:     environment.Version,
			Compiled:    environment.Compiled(),
			Healthy:     true,
			Checks:      checks,
		}

		for _, check := range checks {
			check.Healthy = true
			if err := check.Func(); err != nil {
				check.Healthy = false
				h.Healthy = false
				check.Message = err.Error()
			}
		}

		if !h.Healthy {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application-json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(h); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// AddCheck adds a status check to the status status. If the provided function returns
// an error the check is deemed unhealthy.
func AddCheck(name string, fn func() error) {
	checks = append(checks, &check{
		Name: name,
		Func: fn,
	})
}

// ErrDown is the error returned when an HTTP health check returns a non-200
// status.
var ErrDown = errors.New("application down")

// PerformCheck performs an HTTP health check on localhost:8081/__/health. It returns an error if the health
// check endpoint returns a non-200 status code.
func PerformCheck(ctx context.Context) error {
	cl := &http.Client{
		Timeout: time.Minute,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8081/__/health", nil)
	if err != nil {
		return err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer closers.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ErrDown
	}

	return nil
}
