package protostub

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

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
	conn, err := grpc.NewClient(target, dialOptions(telemetryManager)...)
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
