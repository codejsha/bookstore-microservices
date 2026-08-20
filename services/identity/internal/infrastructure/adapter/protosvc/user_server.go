package protosvc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/pagination"
	"github.com/codejsha/shared-library-go/pkg/ptr"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

var _ userpb.UserServiceServer = (*UserGrpcServer)(nil)

type UserGrpcServer struct {
	identityUseCase usecase.IdentityUseCase
	userpb.UnimplementedUserServiceServer
}

func NewUserGrpcServer(
	identityUseCase usecase.IdentityUseCase,
) *UserGrpcServer {
	return &UserGrpcServer{
		identityUseCase: identityUseCase,
	}
}

func (s *UserGrpcServer) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	opts := []option.UserQueryOptionFunc{}
	if req.GetEmail() != "" {
		email := req.GetEmail()
		opts = append(opts, option.UserQueryOption{}.WithEmail(&email))
	}
	if req.GetName() != "" {
		name := req.GetName()
		opts = append(opts, option.UserQueryOption{}.WithName(&name))
	}
	if req.GetPhone() != "" {
		phone := req.GetPhone()
		opts = append(opts, option.UserQueryOption{}.WithPhone(&phone))
	}
	opts = append(opts, option.UserQueryOption{}.WithPage(newPageOption(req.GetPageSize(), req.GetPageToken(), req.GetOrderBy())))
	opt := option.NewUserQueryOption(opts...)

	total, users, err := s.identityUseCase.FindAllUsers(ctx, opt)
	if err != nil {
		return nil, toGrpcError(err)
	}

	resps := make([]*userpb.User, len(users))
	for i, user := range users {
		resps[i] = toUserProto(user)
	}

	return &userpb.ListUsersResponse{
		Users:     resps,
		TotalSize: int32(total),
	}, nil
}

func (s *UserGrpcServer) FindUser(ctx context.Context, req *userpb.FindUserRequest) (*userpb.FindUserResponse, error) {
	user, err := s.identityUseCase.FindUser(ctx, req.GetUid())
	if err != nil {
		return nil, toGrpcError(err)
	}
	return &userpb.FindUserResponse{User: toUserProto(user)}, nil
}

// ─── mapping helpers ────────────────────────────────────────────────────────

func toGrpcError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return status.Error(codes.NotFound, "user not found")
	}
	return status.Error(codes.Internal, "internal error")
}

func toUserProto(a *aggregate.UserAggregate) *userpb.User {
	roles := a.Roles()
	roleStrs := make([]string, len(roles))
	for i, r := range roles {
		roleStrs[i] = r.ToAuthRoleProto()
	}
	phone := ""
	if p := a.Phone(); p != nil {
		phone = *p
	}
	return &userpb.User{
		Uid:       a.Uid(),
		Email:     a.Email(),
		FirstName: a.FirstName(),
		LastName:  a.LastName(),
		Phone:     phone,
		Roles:     roleStrs,
		Status:    toUserStatusProto(a.Status()),
	}
}

func toUserStatusProto(s string) userpb.UserStatus {
	switch s {
	case "ACTIVE":
		return userpb.UserStatus_USER_STATUS_ACTIVE
	case "SUSPENDED":
		return userpb.UserStatus_USER_STATUS_SUSPENDED
	case "DEACTIVATED":
		return userpb.UserStatus_USER_STATUS_INACTIVE
	default:
		return userpb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

func newPageOption(pageSize int32, _pageToken, orderBy string) pagination.PageOption {
	var size *int32
	if pageSize != 0 {
		size = ptr.Int32(pageSize)
	}
	var sortPtr *string
	if orderBy != "" {
		sortPtr = &orderBy
	}
	return pagination.NewPageOption(size, nil, sortPtr)
}
