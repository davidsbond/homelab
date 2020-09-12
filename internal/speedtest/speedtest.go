// Package speedtest contains methods for performing internet speed tests using speedtest.net.
package speedtest

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/distance"
	"pkg.dsb.dev/health"
)

const (
	baseURL      = "https://speedtest.net"
	randomURIFmt = "/speedtest/random%vx%v.jpg"
)

// ErrStatus is the error given when an HTTP call to speedtest.net returns an error
// status code.
var ErrStatus = errors.New("unexpected status code")

type (
	// The Tester type is responsible for performing internet speed tests.
	Tester struct {
		http *http.Client
	}

	// The Results type contains the results of a speed test.
	Results struct {
		Latency  float64
		Download float64
		Upload   float64
	}
)

// New returns a new instance of the Tester type that can perform a speed test.
func New() *Tester {
	ts := &Tester{
		http: &http.Client{
			Timeout: time.Minute * 10,
		},
	}

	health.AddCheck("speedtest", ts.Ping)
	return ts
}

// Test performs a speed test and returns the results. Can be cancelled using the given context.
func (t *Tester) Test(ctx context.Context) (Results, error) {
	client, err := t.client(ctx)
	if err != nil {
		return Results{}, fmt.Errorf("failed to get client: %w", err)
	}

	servers, err := t.servers(ctx)
	if err != nil {
		return Results{}, fmt.Errorf("failed to list servers: %w", err)
	}

	n, err := nearestServer(client, servers)
	if err != nil {
		return Results{}, fmt.Errorf("failed to calculate nearest server: %w", err)
	}

	nearest := servers[n]
	latency, err := t.latency(ctx, nearest)
	if err != nil {
		return Results{}, fmt.Errorf("failed to get latency: %w", err)
	}

	download, err := t.downloadSpeed(ctx, nearest, latency)
	if err != nil {
		return Results{}, fmt.Errorf("failed to get download speed: %w", err)
	}

	upload, err := t.uploadSpeed(ctx, nearest, latency)
	if err != nil {
		return Results{}, fmt.Errorf("failed to get upload speed: %w", err)
	}

	return Results{
		Latency:  float64(latency),
		Download: download,
		Upload:   upload,
	}, nil
}

// Ping returns a non-nil error if speedtest.net appears to be
// down.
func (t *Tester) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := t.client(ctx)
	return err
}

func (t *Tester) client(ctx context.Context) (Client, error) {
	const uri = "/speedtest-config.php"

	var resp GetClientsResponse
	if err := t.get(ctx, uri, &resp); err != nil {
		return Client{}, err
	}

	return resp.Clients[0], nil
}

func (t *Tester) servers(ctx context.Context) ([]Server, error) {
	const uri = "/speedtest-servers-static.php"

	var resp GetServersResponse
	if err := t.get(ctx, uri, &resp); err != nil {
		return nil, err
	}

	return resp.Servers, nil
}

func (t *Tester) latency(ctx context.Context, server Server) (time.Duration, error) {
	const uri = "/speedtest/latency.txt"
	const attempts = 3

	u, err := url.Parse(server.URL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse server URL: %w", err)
	}
	u.Path = uri

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	latency := time.Duration(math.MaxInt64)
	for i := 0; i < attempts; i++ {
		start := time.Now()
		if err := t.do(req, nil); err != nil {
			return 0, fmt.Errorf("failed to perform latency test: %w", err)
		}
		end := time.Now()

		if end.Sub(start) < latency {
			latency = end.Sub(start)
		}
	}

	return latency, nil
}

func (t *Tester) downloadSpeed(ctx context.Context, server Server, latency time.Duration) (float64, error) {
	u, err := url.Parse(server.URL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse server URL: %w", err)
	}

	warmUpSpeed, err := t.downloadWarmup(ctx, latency, u)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate download warmup speed: %w", err)
	}

	workload, weight, skip := downloadWorkload(warmUpSpeed)
	if skip {
		return warmUpSpeed, nil
	}

	sizes := [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	size := sizes[weight]

	u.Path = fmt.Sprintf(randomURIFmt, size, size)

	grp, _ := errgroup.WithContext(ctx)
	start := time.Now()
	for i := 0; i < workload; i++ {
		grp.Go(func() error {
			resp, err := t.http.Get(u.String())
			if err != nil {
				return err
			}
			defer closers.Close(resp.Body)

			_, err = ioutil.ReadAll(resp.Body)
			return err
		})
	}

	if err = grp.Wait(); err != nil {
		return 0, err
	}
	end := time.Now()

	reqMB := sizes[weight] * sizes[weight] * 2 / 1000 / 1000
	dlSpeed := float64(reqMB) * 8 * float64(workload) / end.Sub(start).Seconds()
	return dlSpeed, nil
}

func (t *Tester) uploadSpeed(ctx context.Context, server Server, latency time.Duration) (float64, error) {
	u, err := url.Parse(server.URL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse server URL: %w", err)
	}

	warmUpSpeed, err := t.uploadWarmup(ctx, latency, u)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate upload warmup speed: %w", err)
	}

	workload, weight, skip := uploadWorkload(warmUpSpeed)
	if skip {
		return warmUpSpeed, nil
	}

	sizes := [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000}
	size := sizes[weight]
	values := url.Values{}
	values.Add("content", strings.Repeat("0123456789", size*100-51))

	grp, _ := errgroup.WithContext(ctx)
	start := time.Now()
	for i := 0; i < workload; i++ {
		grp.Go(func() error {
			resp, err := t.http.PostForm(u.String(), values)
			if err != nil {
				return err
			}
			defer closers.Close(resp.Body)

			_, err = ioutil.ReadAll(resp.Body)
			return err
		})
	}

	if err = grp.Wait(); err != nil {
		return 0, err
	}
	end := time.Now()

	reqMB := float64(sizes[weight]) / 1000
	ulSpeed := reqMB * 8 * float64(workload) / end.Sub(start).Seconds()
	return ulSpeed, nil
}

func (t *Tester) downloadWarmup(ctx context.Context, latency time.Duration, u *url.URL) (float64, error) {
	const attempts = 2
	const requestSize = 1.125 * 8 * 2 // 1.125MB for each request (750 * 750 * 2)
	const downloadSize = 500

	grp, _ := errgroup.WithContext(ctx)
	start := time.Now()

	u.Path = fmt.Sprintf(randomURIFmt, downloadSize, downloadSize)
	for i := 0; i < attempts; i++ {
		grp.Go(func() error {
			resp, err := t.http.Get(u.String())
			if err != nil {
				return err
			}
			defer closers.Close(resp.Body)

			_, err = ioutil.ReadAll(resp.Body)
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		return 0, err
	}

	end := time.Now()

	return requestSize / end.Sub(start.Add(latency)).Seconds(), nil
}

func (t *Tester) uploadWarmup(ctx context.Context, latency time.Duration, u *url.URL) (float64, error) {
	const attempts = 2
	const uploadSize = 800*100 - 51
	const uri = "/speedtest/upload"
	const requestSize = 1.0 * 8 * 2 // 1.0 MB for each request

	grp, _ := errgroup.WithContext(ctx)
	start := time.Now()

	u.Path = uri
	values := url.Values{}
	values.Add("content", strings.Repeat("0123456789", uploadSize))

	for i := 0; i < attempts; i++ {
		grp.Go(func() error {
			resp, err := t.http.PostForm(u.String(), values)
			if err != nil {
				return err
			}
			defer closers.Close(resp.Body)

			_, err = ioutil.ReadAll(resp.Body)
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		return 0, err
	}

	end := time.Now()

	return requestSize / end.Sub(start.Add(latency)).Seconds(), nil
}

func uploadWorkload(speed float64) (workload, weight int, skip bool) {
	switch {
	case speed > 10:
		workload = 16
		weight = 9
	case speed > 4:
		workload = 8
		weight = 9
	case speed > 2.5:
		workload = 4
		weight = 5
	default:
		skip = true
	}

	return
}

func downloadWorkload(speed float64) (workload, weight int, skip bool) {
	switch {
	case speed > 10:
		workload = 16
		weight = 4
	case speed > 4:
		workload = 8
		weight = 4
	case speed > 2.5:
		workload = 4
		weight = 4
	default:
		skip = true
	}

	return
}

func (t *Tester) get(ctx context.Context, uri string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+uri, nil)
	if err != nil {
		return fmt.Errorf("failed to build HTTP request: %w", err)
	}

	return t.do(req, out)
}

func (t *Tester) do(req *http.Request, out interface{}) error {
	res, err := t.http.Do(req)
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

func nearestServer(client Client, servers []Server) (int, error) {
	clientLatitude, err := strconv.ParseFloat(client.Latitude, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid client latitude %s: %w", client.Latitude, err)
	}

	clientLongitude, err := strconv.ParseFloat(client.Longitude, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid client longitude %s: %w", client.Longitude, err)
	}

	var nearest int
	var dist float64
	for i, server := range servers {
		serverLatitude, err := strconv.ParseFloat(server.Latitude, 64)
		if err != nil {
			return -1, fmt.Errorf("invalid server latitude %s: %w", server.Latitude, err)
		}

		serverLongitude, err := strconv.ParseFloat(server.Longitude, 64)
		if err != nil {
			return -1, fmt.Errorf("invalid server longitude %s: %w", server.Longitude, err)
		}

		d := distance.Between(clientLatitude, clientLongitude, serverLatitude, serverLongitude)
		if d < dist {
			nearest = i
			dist = d
		}
	}

	return nearest, nil
}
