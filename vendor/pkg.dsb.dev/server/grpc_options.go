package server

import (
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type (
	// The GRPCOption type is a function that modifies a field on the gRPC configuration.
	GRPCOption func(c *grpcConfig)
)

const (
	defaultMaxConnectionIdle     = 10 * time.Minute
	defaultMaxConnectionAge      = 5 * time.Minute
	defaultMaxConnectionAgeGrace = 5 * time.Minute
	defaultTime                  = 5 * time.Minute
)

var defaultGRPCConfig = grpcConfig{
	serverOptions: []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     defaultMaxConnectionIdle,
			MaxConnectionAge:      defaultMaxConnectionAge,
			MaxConnectionAgeGrace: defaultMaxConnectionAgeGrace,
			Time:                  defaultTime,
		}),
	},
}

// WithServerOptions configures the options to use when calling grpc.NewServer.
func WithServerOptions(opts ...grpc.ServerOption) GRPCOption {
	return func(c *grpcConfig) {
		c.serverOptions = opts
	}
}

// WithService adds a GRPCService implementation whose Register method will be called
// once the server is started.
func WithService(svc GRPCService) GRPCOption {
	return func(c *grpcConfig) {
		c.services = append(c.services, svc)
	}
}
