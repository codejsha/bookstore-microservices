package protostub

import (
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/pb/paymentpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/config"
)

func NewPaymentGrpcClient(
	grpcCfg *config.GrpcConfig,
) *PaymentGrpcClient {
	target := net.JoinHostPort("", grpcCfg.PaymentServer.Port)
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

	client := &PaymentGrpcClient{
		grpcCfg: grpcCfg,
		Client:  paymentpb.NewPaymentServiceClient(conn),
	}
	return client
}

type PaymentGrpcClient struct {
	grpcCfg *config.GrpcConfig
	Client  paymentpb.PaymentServiceClient
}
