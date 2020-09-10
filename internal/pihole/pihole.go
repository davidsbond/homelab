// Package pihole contains utilities for interacting with the pihole admin API.
package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"pkg.dsb.dev/closers"
)

type (
	// The PiHole type is the API interface for a running pihole instance.
	PiHole struct {
		url  *url.URL
		http *http.Client
	}

	// The Summary type contains the summary returned by the pihole instance.
	Summary struct {
		AdsBlocked       float64 `json:"ads_blocked_today"`
		AdsPercentage    float64 `json:"ads_percentage_today"`
		ClientsEverSeen  float64 `json:"clients_ever_seen"`
		DNSQueries       float64 `json:"dns_queries_today"`
		DomainsBlocked   float64 `json:"domains_being_blocked"`
		QueriesCached    float64 `json:"queries_cached"`
		QueriesForwarded float64 `json:"queries_forwarded"`
		Status           string  `json:"status"`
		UniqueClients    float64 `json:"unique_clients"`
		UniqueDomains    float64 `json:"unique_domains"`
	}
)

// New creates a new instance of the PiHole type that will make requests against the
// given url.
func New(urlStr string) (*PiHole, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &PiHole{
		url: u,
		http: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

// Summary returns the pihole summary.
func (pi *PiHole) Summary(ctx context.Context) (Summary, error) {
	const uri = "/admin/api.php?summaryRaw"

	var resp Summary
	if err := pi.get(ctx, uri, &resp); err != nil {
		return Summary{}, err
	}

	return resp, nil
}

func (pi *PiHole) get(ctx context.Context, uri string, out interface{}) error {
	u := *pi.url

	parts := strings.Split(uri, "?")
	u.Path = parts[0]
	u.RawQuery = parts[1]

	urlStr := u.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to build HTTP request: %w", err)
	}

	return pi.do(req, out)
}

func (pi *PiHole) do(req *http.Request, out interface{}) error {
	res, err := pi.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer closers.Close(res.Body)

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}
