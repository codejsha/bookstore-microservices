package protostub

import (
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/deliverypb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

type DeliveryGrpcClient struct {
	grpcCfg *config.GrpcConfig
	conn    *grpc.ClientConn
	Client  deliverypb.DeliveryServiceClient
}

func NewDeliveryGrpcClient(
	grpcCfg *config.GrpcConfig,
	telemetryManager *support.TelemetryManager,
) (*DeliveryGrpcClient, error) {
	target := net.JoinHostPort(grpcCfg.DeliveryServer.Host, grpcCfg.DeliveryServer.Port)
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	option := grpc.WithStatsHandler(otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(telemetryManager.TraceProvider),
		otelgrpc.WithMeterProvider(telemetryManager.MeterProvider),
	))

	conn, err := grpc.NewClient(target, credentials, option)
	if err != nil {
		return nil, fmt.Errorf("create delivery grpc client: %w", err)
	}

	return &DeliveryGrpcClient{
		grpcCfg: grpcCfg,
		conn:    conn,
		Client:  deliverypb.NewDeliveryServiceClient(conn),
	}, nil
}

func (c *DeliveryGrpcClient) Close() error {
	return c.conn.Close()
}
