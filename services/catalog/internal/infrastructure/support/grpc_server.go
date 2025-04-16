package support

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/authorpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/bookpb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/pb/categorypb"
	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/protosvc"
)

func NewGrpcServer(
	grpcCfg *config.GrpcConfig,
	authorGrpcServer *protosvc.AuthorGrpcServer,
	bookGrpcServer *protosvc.BookGrpcServer,
	categoryGrpcServer *protosvc.CategoryGrpcServer,
) *GrpcServer {
	return &GrpcServer{
		grpcCfg:            grpcCfg,
		authorGrpcServer:   authorGrpcServer,
		bookGrpcServer:     bookGrpcServer,
		categoryGrpcServer: categoryGrpcServer,
	}
}

type GrpcServer struct {
	server             *grpc.Server
	grpcCfg            *config.GrpcConfig
	authorGrpcServer   *protosvc.AuthorGrpcServer
	bookGrpcServer     *protosvc.BookGrpcServer
	categoryGrpcServer *protosvc.CategoryGrpcServer
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
	authorpb.RegisterAuthorServiceServer(s.server, s.authorGrpcServer)
	bookpb.RegisterBookServiceServer(s.server, s.bookGrpcServer)
	categorypb.RegisterCategoryServiceServer(s.server, s.categoryGrpcServer)
}
