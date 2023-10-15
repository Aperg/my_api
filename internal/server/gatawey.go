package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"

	desc "cmd/main.go/pkg/my-api"
)

const createGatewayServerLogTag = "createGatewayServer()"

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservice_requests_total",
		Help: "The total number of incoming HTTP requests",
	})
)

func createGatewayServer(ctx context.Context, grpcAddr, gatewayAddr string) *http.Server {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		ctx,
		grpcAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("%s: grpc.DialContext failed", createGatewayServerLogTag),
			"err", err,
		)
	}

	mux := runtime.NewServeMux()
	if err := desc.RegisterApiServiceHandler(ctx, mux, conn); err != nil {
		log.Fatal(ctx, fmt.Sprintf("%s: pb.RegisterBssEquipmentRequestApiServiceHandler failed", createGatewayServerLogTag),
			"err", err,
		)
	}

	gatewayServer := &http.Server{
		Addr: gatewayAddr,
	}

	return gatewayServer
}
