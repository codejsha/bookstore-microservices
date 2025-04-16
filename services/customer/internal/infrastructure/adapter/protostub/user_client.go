package protostub

import (
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
)

func NewUserGrpcClient(
	grpcCfg *config.GrpcConfig,
) *UserGrpcClient {
	target := net.JoinHostPort("", grpcCfg.UserServer.Port)
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(target, credentials)
	if err != nil {
		logrus.Errorf("failed to create gRPC client: %v", err)
		return nil
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logrus.Errorf("failed to close gRPC client: %v", err)
		}
	}(conn)

	client := &UserGrpcClient{
		grpcCfg: grpcCfg,
		Client:  userpb.NewUserServiceClient(conn),
	}
	return client
}

type UserGrpcClient struct {
	grpcCfg *config.GrpcConfig
	Client  userpb.UserServiceClient
}
