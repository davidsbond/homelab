// Package transport contains utility methods for API interfaces that speak to clients
package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"

	"pkg.dsb.dev/logging"
	"pkg.dsb.dev/tracing"
)

type (
	// The HTTP type contains utilities for interacting with HTTP requests.
	HTTP struct{}

	// The Error type represents an error returned by the API.
	Error struct {
		Message string   `json:"message"`
		Stack   []string `json:"stack,omitempty"`
	}
)

func (e Error) Error() string {
	return e.Message
}

// Decode the contents of the request body into the given interface and perform
// validation on the struct.
func (t HTTP) Decode(r *http.Request, out interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return err
	}

	return validator.New().Struct(out)
}

// Respond to an HTTP request with the provided code and JSON-encoded
// response body.
func (t HTTP) Respond(ctx context.Context, w http.ResponseWriter, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Error(ctx, w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// Error returns a JSON-encoded error to the HTTP client.
func (t HTTP) Error(ctx context.Context, w http.ResponseWriter, code int, msg string, args ...interface{}) {
	e := Error{Message: fmt.Sprintf(msg, args...)}
	if span := opentracing.SpanFromContext(ctx); span != nil {
		tracing.AddError(span, e)
	}

	logging.WithError(e).WithFields(map[string]interface{}{
		"status": code,
	}).Error("http request failed")

	t.Respond(ctx, w, e, code)
}

// ErrorWithStack returns a JSON-encoded error to the HTTP client that contains a stack trace. Should be used
// for debugging and handling panics.
func (t HTTP) ErrorWithStack(ctx context.Context, w http.ResponseWriter, code int, msg string, args ...interface{}) {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	str := string(buf[:n])

	e := Error{
		Message: fmt.Sprintf(msg, args...),
		Stack:   strings.Split(str, "\n"),
	}

	logging.WithError(e).WithFields(map[string]interface{}{
		"status": code,
		"stack":  e.Stack,
	}).Error("http request failed")

	t.Respond(ctx, w, e, code)
}

// NotFound returns an http.Handler implementation to be used when a requested
// HTTP resource does not exist.
func NotFound() http.Handler {
	t := HTTP{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(r.Context(), w, http.StatusNotFound, "the requested resource does not exist")
	})
}

// MethodNotAllowed returns an http.Handler implementation to be used when an HTTP call
// for an existing resource is using an invalid method.
func MethodNotAllowed() http.Handler {
	t := HTTP{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error(r.Context(), w, http.StatusMethodNotAllowed, "the provided method is not allowed")
	})
}

// RangeFromValues generates a Content-Range header using the provided
// values.
func RangeFromValues(start, end, size int64) string {
	var builder strings.Builder
	builder.WriteString("bytes ")
	builder.WriteString(strconv.FormatInt(start, 10))
	builder.WriteRune('-')
	builder.WriteString(strconv.FormatInt(end, 10))
	builder.WriteRune('/')
	builder.WriteString(strconv.FormatInt(size, 10))

	return builder.String()
}

var (
	// ErrNoRangeHeader is the error used when attempting to parse a non-existent
	// range header.
	ErrNoRangeHeader = errors.New("no range header provided")

	// ErrInvalidRange is the error used when a range value cannot be parsed.
	ErrInvalidRange = errors.New("invalid range")
)

// RangeFromHeader determines the range of bytes to stream from the Range
// header.
func RangeFromHeader(r *http.Request) (int64, int64, error) {
	const expectedParts = 2

	raw := strings.ReplaceAll(r.Header.Get("Range"), "bytes=", "")
	if raw == "" {
		return 0, 0, ErrNoRangeHeader
	}

	parts := strings.Split(raw, "-")
	if len(parts) != expectedParts {
		return 0, 0, fmt.Errorf("%w: %s", ErrInvalidRange, raw)
	}

	rawStart := parts[0]
	rawEnd := parts[1]

	start, err := strconv.ParseInt(rawStart, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start range %s: %w", rawStart, err)
	}

	if rawEnd == "" {
		return start, -1, nil
	}

	end, err := strconv.ParseInt(rawEnd, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end range %s: %w", rawEnd, err)
	}

	return start, end, nil
}
