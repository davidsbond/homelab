// Package health contains helpers for constructing health checks.
package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/environment"
)

type (
	health struct {
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

// Serve adds the health check endpoint to the given router.
func Serve(r *mux.Router) {
	if disabled {
		return
	}

	r.HandleFunc("/__/health", func(w http.ResponseWriter, r *http.Request) {
		h := &health{
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

// AddCheck adds a health check to the health status. If the provided function returns
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
