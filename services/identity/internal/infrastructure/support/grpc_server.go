package support

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/protosvc"
)

func NewGrpcServer(
	grpcCfg *config.GrpcConfig,
	userGrpcServer *protosvc.UserGrpcServer,
) *GrpcServer {
	return &GrpcServer{
		grpcCfg:        grpcCfg,
		userGrpcServer: userGrpcServer,
	}
}

type GrpcServer struct {
	server         *grpc.Server
	grpcCfg        *config.GrpcConfig
	userGrpcServer *protosvc.UserGrpcServer
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
	userpb.RegisterUserServiceServer(s.server, s.userGrpcServer)
}
