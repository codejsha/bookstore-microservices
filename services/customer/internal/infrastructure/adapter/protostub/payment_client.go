package protostub

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

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
	conn, err := grpc.NewClient(target, dialOptions(telemetryManager)...)
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
