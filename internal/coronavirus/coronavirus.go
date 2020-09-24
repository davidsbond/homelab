// Package coronavirus contains utilities for obtaining UK coronavirus statistics.
package coronavirus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/health"
)

const baseURL = "https://api.coronavirus.data.gov.uk"

type (
	// The Client type is used to communicate with the government coronavirus API.
	Client struct {
		http *http.Client
		url  *url.URL
	}

	// The Summary type contains all metrics for coronavirus data.
	Summary struct {
		Cases  float64
		Deaths float64
	}
)

// ErrStatus is the error given when an HTTP call to api.coronavirus.data.gov.uk returns an error
// status code.
var ErrStatus = errors.New("unexpected status code")

// New creates a new instance of the Client type.
func New() (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	cl := &Client{
		http: &http.Client{
			Timeout: time.Second * 10,
		},
		url: u,
	}

	health.AddCheck(baseURL, cl.Ping)
	return cl, cl.Ping()
}

// Ping returns a non-nil error if the API is deemed to be down.
func (cl *Client) Ping() error {
	const uri = "/v1/data"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return cl.get(ctx, uri, url.Values{}, nil)
}

// GetSummary returns a summary of coronavirus statistics for the provided date.
func (cl *Client) GetSummary(ctx context.Context, date time.Time) (*Summary, error) {
	values := url.Values{}
	values.Set("structure", `{"date":"date","newCases":"newCasesByPublishDate","newDeaths":"newDeaths28DaysByDeathDate"}`)
	values.Set("filters", fmt.Sprintf("areaType=nation;date=%s", date.Format("2006-01-02")))

	const uri = "/v1/data"
	var resp GetDataResponse
	if err := cl.get(ctx, uri, values, &resp); err != nil {
		return nil, err
	}

	var results Summary
	for _, data := range resp.Data {
		results.Cases += float64(data.NewCases)
		if data.NewDeaths != nil {
			results.Deaths += float64(*data.NewDeaths)
		}
	}

	return &results, nil
}

func (cl *Client) get(ctx context.Context, uri string, values url.Values, out interface{}) error {
	u := *cl.url
	u.Path = uri

	query, err := url.QueryUnescape(values.Encode())
	if err != nil {
		return err
	}

	u.RawQuery = query
	urlStr := u.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to build HTTP request: %w", err)
	}

	return cl.do(req, out)
}

func (cl *Client) do(req *http.Request, out interface{}) error {
	res, err := cl.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer closers.Close(res.Body)

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		return fmt.Errorf("%w: %v", ErrStatus, res.StatusCode)
	}

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}
