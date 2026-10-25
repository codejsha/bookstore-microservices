package protostub

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/fulfillmentpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

type DeliveryGrpcClient struct {
	grpcCfg *config.GrpcConfig
	conn    *grpc.ClientConn
	Client  fulfillmentpb.FulfillmentServiceClient
}

func NewDeliveryGrpcClient(
	grpcCfg *config.GrpcConfig,
	telemetryManager *support.TelemetryManager,
) (*DeliveryGrpcClient, error) {
	target := net.JoinHostPort(grpcCfg.FulfillmentServer.Host, grpcCfg.FulfillmentServer.Port)
	conn, err := grpc.NewClient(target, dialOptions(telemetryManager)...)
	if err != nil {
		return nil, fmt.Errorf("create fulfillment grpc client: %w", err)
	}

	return &DeliveryGrpcClient{
		grpcCfg: grpcCfg,
		conn:    conn,
		Client:  fulfillmentpb.NewFulfillmentServiceClient(conn),
	}, nil
}

func (c *DeliveryGrpcClient) Close() error {
	return c.conn.Close()
}
