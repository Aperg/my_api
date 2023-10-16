package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"cmd/main.go/internal/pkg/grps_logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"cmd/main.go/internal/api"
	"cmd/main.go/internal/logger"
	"cmd/main.go/internal/service/user_request"
	desc "cmd/main.go/pkg/my-api"

	"cmd/main.go/internal/config"
)

const grpcServerStartLogTag = "GrpcServer.Start()"

// GrpcServer is gRPC server
type GrpcServer struct {
	userRequestService user_request.ServiceInterface
}

// NewGrpcServer returns gRPC server with supporting of batch listing
func NewGrpcServer(userRequestService user_request.ServiceInterface) *GrpcServer {
	return &GrpcServer{userRequestService: userRequestService}
}

// Start method runs server
func (s *GrpcServer) Start(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)

	gatewayServer := createGatewayServer(ctx, grpcAddr, gatewayAddr)

	go func() {

		logger.InfoKV(ctx, fmt.Sprintf("%s: gateway server is running on", grpcServerStartLogTag),
			"address", grpcAddr)
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed running gateway server")
			cancel()

		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	//nolint
	defer l.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			grpc_zap.PayloadUnaryServerInterceptor(logger.Clone(ctx), grps_logger.ServerPayloadLoggingDecider()),
			grps_logger.UnaryServerInterceptor(),
		)),
	)

	desc.RegisterApiServiceServer(grpcServer, api.NewApiService(s.userRequestService))

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("%s: GRPC server is listening on", grpcServerStartLogTag),
			"address", grpcAddr,
		)
		if err := grpcServer.Serve(l); err != nil {
			log.Fatal().Err(err).Msg("Failed running gRPC server")
		}

	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		log.Info().Msg("The service is ready to accept requests")
	}()

	if cfg.Project.Debug {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		logger.InfoKV(ctx, fmt.Sprintf("%s: signal.Notify", grpcServerStartLogTag), "quit", v)
	case done := <-ctx.Done():
		logger.InfoKV(ctx, fmt.Sprintf("%s: ctx.Done", grpcServerStartLogTag), "done", done)
	}

	isReady.Store(false)

	if err := gatewayServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, fmt.Sprintf("%s: gatewayServer.Shutdown failed", grpcServerStartLogTag), "err", err)
	} else {
		logger.Info(ctx, fmt.Sprintf("%s: gatewayServer shut down correctly", grpcServerStartLogTag))
	}

	grpcServer.GracefulStop()
	logger.Info(ctx, fmt.Sprintf("%s: grpcServer shut down correctly", grpcServerStartLogTag))

	return nil
}
