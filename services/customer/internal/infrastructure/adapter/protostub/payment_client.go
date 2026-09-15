package protostub

import (
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/paymentpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

type PaymentGrpcClient struct {
	grpcCfg *config.GrpcConfig
	conn    *grpc.ClientConn
	Client  paymentpb.PaymentServiceClient
}

func NewPaymentGrpcClient(
	grpcCfg *config.GrpcConfig,
	telemetryManager *support.TelemetryManager,
) (*PaymentGrpcClient, error) {
	target := net.JoinHostPort(grpcCfg.PaymentServer.Host, grpcCfg.PaymentServer.Port)
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	option := grpc.WithStatsHandler(otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(telemetryManager.TraceProvider),
		otelgrpc.WithMeterProvider(telemetryManager.MeterProvider),
	))

	conn, err := grpc.NewClient(target, credentials, option)
	if err != nil {
		return nil, fmt.Errorf("create payment grpc client: %w", err)
	}

	return &PaymentGrpcClient{
		grpcCfg: grpcCfg,
		conn:    conn,
		Client:  paymentpb.NewPaymentServiceClient(conn),
	}, nil
}

func (c *PaymentGrpcClient) Close() error {
	return c.conn.Close()
}
