package transport

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pkg.dsb.dev/logging"
	"pkg.dsb.dev/tracing"
	v1 "pkg.dsb.dev/transport/v1"
)

type (
	// The GRPC type contains utilities for interacting with gRPC requests.
	GRPC struct{}
)

// Error returns a proto-encoded error to the gRPC client.
func (t GRPC) Error(ctx context.Context, code codes.Code, msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...) // nolint: goerr113

	if span := opentracing.SpanFromContext(ctx); span != nil {
		tracing.AddError(span, err)
	}

	logging.WithError(err).WithFields(map[string]interface{}{
		"status": code,
	}).Error("gRPC request failed")

	return status.Error(code, err.Error())
}

// ErrorWithStack returns a proto-encoded error to the gRPC client that contains a stack trace. Should be used
// for debugging and handling panics. The stack trace is contained within the details field of the grpc status.
func (t GRPC) ErrorWithStack(ctx context.Context, code codes.Code, msg string, args ...interface{}) error {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	str := string(buf[:n])
	details := &v1.GRPCErrorDetails{
		StackTrace: strings.Split(str, "\n"),
	}

	err := fmt.Errorf(msg, args...) // nolint: goerr113
	if span := opentracing.SpanFromContext(ctx); span != nil {
		tracing.AddError(span, err)
	}

	logging.WithError(err).WithFields(map[string]interface{}{
		"status": code,
		"stack":  details.GetStackTrace(),
	}).Error("gRPC request failed")

	st, err := status.New(code, err.Error()).WithDetails(details)
	if err != nil {
		panic(err)
	}

	return st.Err()
}
