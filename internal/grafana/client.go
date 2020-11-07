// Package grafana contains a client implementation for obtaining data sources and dashboards
// from a grafana instance.
package grafana

import (
	"context"
	"time"

	"github.com/grafana-tools/sdk"
	"pkg.dsb.dev/health"
)

type (
	// The Client type contains methods for interacting with the grafana instance.
	Client struct {
		grafana *sdk.Client
	}
)

// NewClient creates a new grafana client for the given url that authenticates using the given API key, returns an
// error if the health check fails.
func NewClient(url, apiKey string) (*Client, error) {
	grafana := sdk.NewClient(url, apiKey, sdk.DefaultHTTPClient)

	cl := &Client{grafana: grafana}
	health.AddCheck(url, cl.Ping)
	return cl, cl.Ping()
}

// Ping returns a non-nil error if the grafana instance is deemed unhealthy.
func (cl *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if _, err := cl.grafana.GetHealth(ctx); err != nil {
		return err
	}

	return nil
}
