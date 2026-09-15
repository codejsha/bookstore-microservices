package protostub

import (
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/orderpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

type OrderGrpcClient struct {
	grpcCfg *config.GrpcConfig
	conn    *grpc.ClientConn
	Client  orderpb.OrderServiceClient
}

func NewOrderGrpcClient(
	grpcCfg *config.GrpcConfig,
	telemetryManager *support.TelemetryManager,
) (*OrderGrpcClient, error) {
	target := net.JoinHostPort(grpcCfg.OrderServer.Host, grpcCfg.OrderServer.Port)
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	option := grpc.WithStatsHandler(otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(telemetryManager.TraceProvider),
		otelgrpc.WithMeterProvider(telemetryManager.MeterProvider),
	))

	conn, err := grpc.NewClient(target, credentials, option)
	if err != nil {
		return nil, fmt.Errorf("create order grpc client: %w", err)
	}

	return &OrderGrpcClient{
		grpcCfg: grpcCfg,
		conn:    conn,
		Client:  orderpb.NewOrderServiceClient(conn),
	}, nil
}

func (c *OrderGrpcClient) Close() error {
	return c.conn.Close()
}
