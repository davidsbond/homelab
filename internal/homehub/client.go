// Package homehub contains utilities for scraping metrics from a BT home hub.
package homehub

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/health"
)

type (
	// The Client type is used to interact with the BT home hub via
	// HTTP.
	Client struct {
		url  *url.URL
		http *http.Client
	}

	// The Summary type contains metrics scraped from the BT home hub.
	Summary struct {
		Uptime    float64
		BytesUp   float64
		BytesDown float64
	}
)

var numberRegexp = regexp.MustCompile(`\d+`)

// ErrStatus is the error given when an HTTP call to the homehub returns an error
// status code.
var ErrStatus = errors.New("unexpected status code")

// ErrInvalidUpDownBytes is the error given when the up/down bytes in the homehub's
// XML response cannot be parsed.
var ErrInvalidUpDownBytes = errors.New("invalid value for up/down bytes")

// New creates a new instance of the Client type that will scrape a BT home hub at the
// provided url.
func New(urlStr string) (*Client, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		url: u,
		http: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	health.AddCheck(urlStr, c.Ping)
	return c, c.Ping()
}

// Ping returns a non-nil error if the root path of the configured URL returns a non
// success status code when an HTTP GET call is made to it.
func (cl *Client) Ping() error {
	const uri = "/"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return cl.get(ctx, uri, nil)
}

// Summary returns a summary of scraped metrics from the wan_conn.xml endpoint on the
// home hub.
func (cl *Client) Summary(ctx context.Context) (*Summary, error) {
	var resp GetStatusResponse
	const uri = "/nonAuth/wan_conn.xml"

	if err := cl.get(ctx, uri, &resp); err != nil {
		return nil, err
	}

	return parseSummary(resp)
}

func parseSummary(resp GetStatusResponse) (*Summary, error) {
	var s Summary
	var err error

	s.Uptime, err = strconv.ParseFloat(resp.Sysuptime.Value, 64)
	if err != nil {
		return nil, err
	}

	// The XML contains URL encoded values.
	rawWanConnVolume, err := url.QueryUnescape(resp.WanConnVolumeList.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wan conn volume: %w", err)
	}

	// For this particular value, it's a URL-encoded JSON array of numbers, as strings,
	// separated by semi-colons. It's easier to just regex out the numbers.
	matches := numberRegexp.FindAllString(rawWanConnVolume, -1)
	if len(matches) < 6 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUpDownBytes, rawWanConnVolume)
	}

	rawDownBytes := matches[4]
	rawUpBytes := matches[5]

	s.BytesUp, err = strconv.ParseFloat(rawUpBytes, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse up bytes (%s): %w", rawUpBytes, err)
	}

	s.BytesDown, err = strconv.ParseFloat(rawDownBytes, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse down bytes (%s): %w", rawDownBytes, err)
	}

	return &s, nil
}

func (cl *Client) get(ctx context.Context, uri string, out interface{}) error {
	u := *cl.url
	u.Path = uri

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
		if err := xml.NewDecoder(res.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}
