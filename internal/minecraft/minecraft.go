// Package minecraft contains utilities for obtaining statistics on minecraft servers.
package minecraft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"pkg.dsb.dev/closers"
)

type (
	// The Client type is used to obtain statistical information on Minecraft servers. It makes HTTP requests to
	// https://ac.mcsrvstat.us, using their API to get server info.
	Client struct {
		http *http.Client
		addr string
	}
)

const baseURL = "https://api.mcsrvstat.us/2/"

// NewClient returns a new instance of the Client type that will check the status of the provided minecraft
// server address.
func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
		http: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetStatus returns a ServerStatus instance containing statistical information about the desired Minecraft
// server.
func (c *Client) GetStatus(ctx context.Context) (*ServerStatus, error) {
	urlStr := baseURL + c.addr
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	var resp ServerStatus
	if err = c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ErrStatus is the error given when an HTTP call to the minecraft status endpoint returns an error
// status code.
var ErrStatus = errors.New("unexpected status code")

func (c *Client) do(req *http.Request, out interface{}) error {
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer closers.Close(res.Body)

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		return fmt.Errorf("%w: %v", ErrStatus, res.StatusCode)
	}

	if out != nil {
		if err = json.NewDecoder(res.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}
