package protostub

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/support"
)

type UserGrpcClient struct {
	grpcCfg *config.GrpcConfig
	conn    *grpc.ClientConn
	Client  userpb.UserServiceClient
}

func NewUserGrpcClient(
	grpcCfg *config.GrpcConfig,
	telemetryManager *support.TelemetryManager,
) (*UserGrpcClient, error) {
	target := net.JoinHostPort(grpcCfg.UserServer.Host, grpcCfg.UserServer.Port)
	conn, err := grpc.NewClient(target, dialOptions(telemetryManager)...)
	if err != nil {
		return nil, fmt.Errorf("create user grpc client: %w", err)
	}

	return &UserGrpcClient{
		grpcCfg: grpcCfg,
		conn:    conn,
		Client:  userpb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserGrpcClient) Close() error {
	return c.conn.Close()
}
