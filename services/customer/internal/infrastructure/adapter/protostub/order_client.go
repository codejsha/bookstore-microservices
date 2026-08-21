package protostub

import (
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/orderpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
)

func NewOrderGrpcClient(
	grpcCfg *config.GrpcConfig,
) *OrderGrpcClient {
	target := net.JoinHostPort("", grpcCfg.OrderServer.Port)
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

	client := &OrderGrpcClient{
		grpcCfg: grpcCfg,
		Client:  orderpb.NewOrderServiceClient(conn),
	}
	return client
}

type OrderGrpcClient struct {
	grpcCfg *config.GrpcConfig
	Client  orderpb.OrderServiceClient
}
