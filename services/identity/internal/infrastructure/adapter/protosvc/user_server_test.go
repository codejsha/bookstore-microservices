package protosvc

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

type stubUseCase struct {
	findAllUsers func(context.Context, option.UserQueryOption) (int64, []*aggregate.UserAggregate, error)
	findUser     func(context.Context, string) (*aggregate.UserAggregate, error)
}

var _ usecase.IdentityUseCase = (*stubUseCase)(nil)

func (s *stubUseCase) FindAllUsers(ctx context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
	return s.findAllUsers(ctx, opt)
}
func (s *stubUseCase) FindUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	return s.findUser(ctx, uid)
}

func (s *stubUseCase) RegisterUser(context.Context, command.UserRegisterCommand) (*aggregate.UserAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) FindUserByEmail(context.Context, string) (*aggregate.UserAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) UpdateUser(context.Context, string, command.UserUpdateCommand) (*aggregate.UserAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) UpdateUserRoles(context.Context, string, command.UserRolesCommand) (*aggregate.UserAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) SuspendUser(context.Context, string) (*aggregate.UserAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) ReactivateUser(context.Context, string) (*aggregate.UserAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) DeactivateUser(context.Context, string) (*aggregate.UserAggregate, error) {
	return nil, nil
}
func (s *stubUseCase) SyncUserFromIdp(context.Context, string) (*aggregate.UserAggregate, error) {
	return nil, nil
}

func newAggregate(idpId, email string, roles []string) *aggregate.UserAggregate {
	idp := idpId
	return aggregate.NewUserAggregate(&repo.UserResult{
		Id: 1, IdpId: &idp, Email: email, FirstName: "F", LastName: "L",
		Roles: roles, Status: "ACTIVE", CreatedAt: time.Now(),
	})
}

func TestUserGrpcServer_WhenUserExists_ReturnsUser(t *testing.T) {
	use := &stubUseCase{
		findUser: func(_ context.Context, uid string) (*aggregate.UserAggregate, error) {
			if uid != "u-1" {
				t.Errorf("uid = %q", uid)
			}
			return newAggregate("idp-1", "u@x.com", []string{"PROFILE", "ORDER"}), nil
		},
	}
	srv := NewUserGrpcServer(use)
	resp, err := srv.FindUser(context.Background(), &userpb.FindUserRequest{Uid: "u-1"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	u := resp.GetUser()
	if u.GetUid() != "idp-1" || u.GetEmail() != "u@x.com" {
		t.Errorf("user = %+v", u)
	}
	if len(u.GetRoles()) != 2 || u.GetRoles()[0] != "PROFILE" || u.GetRoles()[1] != "ORDER" {
		t.Errorf("Roles = %v", u.GetRoles())
	}
	if u.GetStatus() != userpb.UserStatus_USER_STATUS_ACTIVE {
		t.Errorf("Status = %v, want ACTIVE", u.GetStatus())
	}
}

func TestUserGrpcServer_WhenUserMissing_ReturnsNotFoundStatus(t *testing.T) {
	use := &stubUseCase{
		findUser: func(context.Context, string) (*aggregate.UserAggregate, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	srv := NewUserGrpcServer(use)
	_, err := srv.FindUser(context.Background(), &userpb.FindUserRequest{Uid: "u-1"})
	if status.Code(err) != codes.NotFound {
		t.Errorf("code = %v, want NotFound", status.Code(err))
	}
}

func TestUserGrpcServer_WhenLookupFailsUnexpectedly_ReturnsInternalStatusWithoutDetail(t *testing.T) {
	use := &stubUseCase{
		findUser: func(context.Context, string) (*aggregate.UserAggregate, error) {
			return nil, errors.New("secret internal db dsn leaked")
		},
	}
	srv := NewUserGrpcServer(use)
	_, err := srv.FindUser(context.Background(), &userpb.FindUserRequest{Uid: "u-1"})
	if status.Code(err) != codes.Internal {
		t.Errorf("code = %v, want Internal", status.Code(err))
	}
	if status.Convert(err).Message() != "internal error" {
		t.Errorf("message = %q, want generic 'internal error' (must not leak raw error)", status.Convert(err).Message())
	}
}

func TestUserGrpcServer_WhenFiltersGiven_ForwardsThemToTheUsecase(t *testing.T) {
	use := &stubUseCase{
		findAllUsers: func(_ context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
			if opt.Email() == nil || *opt.Email() != "x@x.com" {
				t.Errorf("Email = %v", opt.Email())
			}
			if opt.Name() == nil || *opt.Name() != "Alice" {
				t.Errorf("Name = %v", opt.Name())
			}
			if opt.Phone() == nil || *opt.Phone() != "+1" {
				t.Errorf("Phone = %v", opt.Phone())
			}
			return 2, []*aggregate.UserAggregate{
				newAggregate("idp-1", "a@x.com", []string{"VIEW"}),
				newAggregate("idp-2", "b@x.com", []string{"VIEW"}),
			}, nil
		},
	}
	srv := NewUserGrpcServer(use)
	resp, err := srv.ListUsers(context.Background(), &userpb.ListUsersRequest{
		Email:    "x@x.com",
		Name:     "Alice",
		Phone:    "+1",
		PageSize: 50,
		OrderBy:  "email",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.GetTotalSize() != 2 || len(resp.GetUsers()) != 2 {
		t.Errorf("resp = %+v", resp)
	}
	if resp.GetUsers()[0].GetUid() != "idp-1" {
		t.Errorf("user[0].Uid = %q", resp.GetUsers()[0].GetUid())
	}
}

func TestNewPageOption_WhenPagingGiven_ReturnsOption(t *testing.T) {
	p := newPageOption(0, "", "")
	if p.GetSize() != 0 {
		t.Errorf("zero pageSize -> Size %d", p.GetSize())
	}
	p = newPageOption(50, "ignored", "title")
	if p.GetSize() != 50 || p.GetSort() != "title" {
		t.Errorf("got Size=%d Sort=%q, want 50/title", p.GetSize(), p.GetSort())
	}
}
