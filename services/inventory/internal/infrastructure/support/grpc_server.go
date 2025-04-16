package support

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/pb/stockpb"
	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/protosvc"
)

func NewGrpcServer(
	grpcCfg *config.GrpcConfig,
	stockGrpcServer *protosvc.StockGrpcServer,
) *GrpcServer {
	return &GrpcServer{
		grpcCfg:         grpcCfg,
		stockGrpcServer: stockGrpcServer,
	}
}

type GrpcServer struct {
	server          *grpc.Server
	grpcCfg         *config.GrpcConfig
	stockGrpcServer *protosvc.StockGrpcServer
}

func (s *GrpcServer) Run() {
	serverCfg := s.grpcCfg.Server
	addr := net.JoinHostPort(serverCfg.Host, serverCfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	s.server = grpc.NewServer()
	s.registerGrpcServers()

	reflection.Register(s.server)
	err = s.server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (s *GrpcServer) registerGrpcServers() {
	stockpb.RegisterStockServiceServer(s.server, s.stockGrpcServer)
}
