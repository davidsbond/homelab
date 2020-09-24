// Package weather contains utilities for obtaining weather statistics from weatherapi.com.
package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"pkg.dsb.dev/health"

	"pkg.dsb.dev/closers"
)

const baseURL = "https://api.weatherapi.com"

type (
	// The Client type is used to interact with the weather API via HTTP.
	Client struct {
		http   *http.Client
		url    *url.URL
		apiKey string
	}
)

// ErrStatus is the error given when an HTTP call to weatherapi.com returns an error
// status code.
var ErrStatus = errors.New("unexpected status code")

// New creates a new Client instance that will authenticate using the given api key.
func New(apiKey string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	cl := &Client{
		url: u,
		http: &http.Client{
			Timeout: time.Second * 10,
		},
		apiKey: apiKey,
	}

	health.AddCheck(baseURL, cl.Ping)
	return cl, cl.Ping()
}

// GetWeather returns the current weather for a given location. Postcodes & city names are valid
// locations.
func (cl *Client) GetWeather(ctx context.Context, location string) (*Weather, error) {
	const uri = "/v1/current.json"

	var resp GetCurrentWeatherResponse
	if err := cl.get(ctx, uri, url.Values{"q": []string{location}}, &resp); err != nil {
		return nil, err
	}

	return &resp.Current, nil
}

// Ping determines if the connection to the API is healthy. If not, a non-nil error is returned.
func (cl *Client) Ping() error {
	const uri = "/v1/search.json"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return cl.get(ctx, uri, url.Values{
		"q": []string{"London"},
	}, nil)
}

func (cl *Client) get(ctx context.Context, uri string, values url.Values, out interface{}) error {
	u := *cl.url
	u.Path = uri

	values.Set("key", cl.apiKey)
	u.RawQuery = values.Encode()

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
