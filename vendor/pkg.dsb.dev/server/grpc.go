package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"pkg.dsb.dev/health"
	"pkg.dsb.dev/logging"
)

type (
	// The GRPCService interface describes types that can be registered to a gRPC service. Each service
	// implementation should have a Register method that registers the type to the gRPC server.
	GRPCService interface {
		Register(svr *grpc.Server)
	}

	grpcConfig struct {
		services      []GRPCService
		serverOptions []grpc.ServerOption
	}
)

// ServeGRPC starts a gRPC server listening on port 5000 configured using the provided options. This function
// blocks until the provided context is cancelled. On cancellation, the server is gracefully stopped.
func ServeGRPC(ctx context.Context, opts ...GRPCOption) error {
	c := defaultGRPCConfig
	for _, opt := range opts {
		opt(&c)
	}

	svr := grpc.NewServer(c.serverOptions...)
	for _, svc := range c.services {
		svc.Register(svr)
	}

	hsvr := health.RegisterGRPCHealthServer(svr)
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		logging.WithField("port", ":5000").Info("serving gRPC")
		return svr.Serve(lis)
	})
	grp.Go(func() error {
		return health.ServeGRPC(ctx, hsvr)
	})
	grp.Go(func() error {
		<-ctx.Done()
		svr.GracefulStop()
		return nil
	})

	err = grp.Wait()
	switch {
	case errors.Is(err, grpc.ErrServerStopped):
		logging.Info("server shut down")
		return nil
	case err != nil:
		return err
	default:
		return nil
	}
}
