package protostub

import (
	"net"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/stockpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
)

func NewStockGrpcClient(
	grpcCfg *config.GrpcConfig,
) *StockGrpcClient {
	target := net.JoinHostPort("", grpcCfg.StockServer.Port)
	credentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	option := grpc.WithStatsHandler(otelgrpc.NewClientHandler())

	conn, err := grpc.NewClient(target, credentials, option)
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

	client := &StockGrpcClient{
		grpcCfg: grpcCfg,
		Client:  stockpb.NewStockServiceClient(conn),
	}
	return client
}

type StockGrpcClient struct {
	grpcCfg *config.GrpcConfig
	Client  stockpb.StockServiceClient
}
