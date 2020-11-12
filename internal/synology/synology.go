// Package synology contains an HTTP client implementation that can scrape metrics from a synology
// diskstation NAS.
package synology

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/health"
)

const (
	systemInfoURI = "/webman/modules/SystemInfoApp/SystemInfo.cgi"
	loginURI      = "/webapi/auth.cgi"
)

type (
	// The Client type is responsible for communicating with the synology NAS.
	Client struct {
		http     *http.Client
		url      *url.URL
		token    string
		user     string
		password string
		cookies  []*http.Cookie
	}
)

// New creates a new client for the synology NAS. Returns an error if authentication
// fails.
func New(urlStr, user, password string) (*Client, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		http:     &http.Client{Timeout: time.Minute},
		url:      u,
		user:     user,
		password: password,
	}

	health.AddCheck(u.Host, c.Ping)
	return c, c.Ping()
}

// Ping determines if the connection to the synology is healthy by authenticating.
func (cl *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return cl.authenticate(ctx)
}

type (
	// SystemInfo contains an overview of information on the synology
	// NAS.
	SystemInfo struct {
		Disks   []Disk
		Volumes []Volume
		Uptime  float64
	}

	// The Disk type contains information on an individual disk in the NAS.
	Disk struct {
		Name        string
		Temperature float64
		Size        float64
	}

	// The Volume type contains information on an individual volume on a disk
	// in the NAS.
	Volume struct {
		Size float64
		Used float64
		Name string
	}
)

// SystemInfo returns the current system information from the NAS.
func (cl *Client) SystemInfo(ctx context.Context) (*SystemInfo, error) {
	var resp GetSystemInfoResponse
	err := cl.get(ctx, systemInfoURI, url.Values{
		"query": []string{"systemHealth"},
	}, &resp)
	if err != nil {
		return nil, err
	}

	info := &SystemInfo{}
	for _, dsk := range resp.Disks {
		disk := Disk{
			Name:        dsk.Name,
			Temperature: dsk.Temp,
		}

		disk.Size, err = strconv.ParseFloat(dsk.SizeTotal, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s size: %w", dsk.Name, err)
		}

		info.Disks = append(info.Disks, disk)
	}

	for _, vol := range resp.VolInfo {
		volume := Volume{
			Name: vol.Name,
		}

		volume.Size, err = strconv.ParseFloat(vol.TotalSize, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s size: %w", vol.Name, err)
		}

		volume.Used, err = strconv.ParseFloat(vol.UsedSize, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s used size: %w", vol.Name, err)
		}

		info.Volumes = append(info.Volumes, volume)
	}

	uptime, err := parseUptime(resp.Optime)
	if err != nil {
		return nil, err
	}

	info.Uptime = uptime.Seconds()
	return info, nil
}

func (cl *Client) authenticate(ctx context.Context) error {
	var resp LoginResponse
	err := cl.get(ctx, loginURI, url.Values{
		"api":     []string{"SYNO.API.Auth"},
		"version": []string{"3"},
		"method":  []string{"login"},
		"account": []string{cl.user},
		"passwd":  []string{cl.password},
	}, &resp)
	if err != nil {
		return err
	}

	cl.token = resp.Data.Sid
	return nil
}

func (cl *Client) get(ctx context.Context, uri string, values url.Values, out interface{}) error {
	u := *cl.url
	u.Path = uri

	if cl.token != "" {
		values.Set("SynoToken", cl.token)
	}

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
	for _, cookie := range cl.cookies {
		req.AddCookie(cookie)
	}

	res, err := cl.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer closers.Close(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var apiError APIError
	if err := json.Unmarshal(body, &apiError); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	if apiError.Err.Code > 0 || apiError.Errno.Key != "" {
		return apiError
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	if len(res.Cookies()) > 0 {
		cl.cookies = res.Cookies()
	}

	return nil
}

// ErrInvalidUptime is the error given when the uptime value cannot be parsed.
var ErrInvalidUptime = errors.New("invalid uptime")

func parseUptime(str string) (time.Duration, error) {
	parts := strings.Split(str, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("%w: %s", ErrInvalidUptime, str)
	}

	hours, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidUptime, str)
	}

	mins, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidUptime, str)
	}

	secs, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidUptime, str)
	}

	result := (time.Hour * time.Duration(hours)) + (time.Minute * time.Duration(mins)) + (time.Second * time.Duration(secs))
	return result, nil
}
