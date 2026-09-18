package protostub

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

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
	conn, err := grpc.NewClient(target, dialOptions(telemetryManager)...)
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
