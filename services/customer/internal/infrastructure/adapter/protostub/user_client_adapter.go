package protostub

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/userpb"
	port "github.com/codejsha/bookstore-microservices/customer/internal/application/port/protostub"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

var _ port.UserClient = (*userClient)(nil)

const userAuthorizationMetadataKey = "authorization"

func withUserAuthorization(ctx context.Context) context.Context {
	if token := httpx.BearerTokenFromContext(ctx); token != "" {
		return metadata.AppendToOutgoingContext(ctx, userAuthorizationMetadataKey, token)
	}
	return ctx
}

type userClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(c *UserGrpcClient) port.UserClient {
	return &userClient{client: c.Client}
}

func (g *userClient) ListUsers(ctx context.Context, email, name, phone string, pageSize int32) (int64, []*aggregate.CustomerAggregate, error) {
	ctx = withUserAuthorization(ctx)
	resp, err := g.client.ListUsers(ctx, &userpb.ListUsersRequest{
		Email:    email,
		Name:     name,
		Phone:    phone,
		PageSize: pageSize,
	})
	if err != nil {
		return 0, nil, err
	}
	out := make([]*aggregate.CustomerAggregate, len(resp.GetUsers()))
	for i, u := range resp.GetUsers() {
		out[i] = toCustomerAggregate(u)
	}
	return int64(resp.GetTotalSize()), out, nil
}

func (g *userClient) FindUser(ctx context.Context, uid string) (*aggregate.CustomerAggregate, error) {
	ctx = withUserAuthorization(ctx)
	resp, err := g.client.FindUser(ctx, &userpb.FindUserRequest{Uid: uid})
	if err != nil {
		return nil, err
	}
	if resp.GetUser() == nil {
		return nil, nil
	}
	return toCustomerAggregate(resp.GetUser()), nil
}

func toCustomerAggregate(u *userpb.User) *aggregate.CustomerAggregate {
	roles := make([]constant.AuthRole, len(u.GetRoles()))
	for i, r := range u.GetRoles() {
		roles[i] = constant.AuthRoleFromString(r)
	}
	var phone *string
	if p := u.GetPhone(); p != "" {
		phone = &p
	}
	return &aggregate.CustomerAggregate{
		Uid:       u.GetUid(),
		Email:     u.GetEmail(),
		FirstName: u.GetFirstName(),
		LastName:  u.GetLastName(),
		Phone:     phone,
		Roles:     roles,
	}
}
