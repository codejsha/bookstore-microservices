package protosvc

import (
	"context"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
)

var _ userpb.UserServiceServer = (*UserGrpcServer)(nil)

func NewUserGrpcServer(
	identityUseCase usecase.IdentityUseCase,
) *UserGrpcServer {
	return &UserGrpcServer{
		identityUseCase: identityUseCase,
	}
}

type UserGrpcServer struct {
	identityUseCase usecase.IdentityUseCase
	userpb.UnimplementedUserServiceServer
}

func (s UserGrpcServer) FindAllUsers(ctx context.Context, req *userpb.UserFindAllProtoReq) (*userpb.UserFindAllProtoResp, error) {
	total, users, err := s.identityUseCase.FindAllUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	resps := make([]*userpb.UserFindProtoResp, len(users))
	for i, user := range users {
		resps[i] = user.ToUserFindProtoResp()
	}
	response := &userpb.UserFindAllProtoResp{Total: total, Items: resps}
	return response, nil
}

func (s UserGrpcServer) FindUser(ctx context.Context, req *userpb.UserFindProtoReq) (*userpb.UserFindProtoResp, error) {
	// TODO implement me
	panic("implement me")
}

func (s UserGrpcServer) UpdateUser(ctx context.Context, req *userpb.UserUpdateProtoReq) (*userpb.UserUpdateProtoResp, error) {
	// TODO implement me
	panic("implement me")
}

func (s UserGrpcServer) DeleteUser(ctx context.Context, req *userpb.UserDeleteProtoReq) (*userpb.UserDeleteProtoResp, error) {
	// TODO implement me
	panic("implement me")
}
