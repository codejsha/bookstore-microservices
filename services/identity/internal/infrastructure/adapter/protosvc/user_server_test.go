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
			return newAggregate("idp-1", "u@x.com", []string{"USER", "STAFF"}), nil
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
	if len(u.GetRoles()) != 2 || u.GetRoles()[0] != "USER" || u.GetRoles()[1] != "STAFF" {
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
				newAggregate("idp-1", "a@x.com", []string{"USER"}),
				newAggregate("idp-2", "b@x.com", []string{"USER"}),
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

func TestNewPageOption_tokenGiven_setsPageIndex(t *testing.T) {
	p, err := newPageOption(0, "", "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p.GetSize() != 0 || p.GetPage() != 0 {
		t.Errorf("empty paging -> Size %d Page %d", p.GetSize(), p.GetPage())
	}
	p, err = newPageOption(50, "3", "title")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p.GetSize() != 50 || p.GetPage() != 3 || p.GetSort() != "title" {
		t.Errorf("got Size=%d Page=%d Sort=%q, want 50/3/title", p.GetSize(), p.GetPage(), p.GetSort())
	}
}

func TestNewPageOption_tokenNotAPageIndex_fails(t *testing.T) {
	for _, token := range []string{"abc", "-1", "1.5", "99999999999"} {
		if _, err := newPageOption(20, token, ""); err == nil {
			t.Errorf("token %q accepted", token)
		}
	}
}

func TestNextPageToken(t *testing.T) {
	cases := []struct {
		name  string
		size  int32
		token string
		total int64
		want  string
	}{
		{"firstOfThreePages_pointsToSecond", 20, "", 45, "1"},
		{"secondOfThreePages_pointsToThird", 20, "1", 45, "2"},
		{"lastPage_empty", 20, "2", 45, ""},
		{"singlePage_empty", 20, "", 20, ""},
		{"sizeUnset_usesRepositoryDefault", 0, "", 21, "1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			page, err := newPageOption(c.size, c.token, "")
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if got := nextPageToken(page, c.total); got != c.want {
				t.Errorf("next token = %q, want %q", got, c.want)
			}
		})
	}
}

func TestUserGrpcServer_ListUsers_statusGiven_forwardsDomainStatus(t *testing.T) {
	cases := []struct {
		name string
		in   userpb.UserStatus
		want *string
	}{
		{"active", userpb.UserStatus_USER_STATUS_ACTIVE, strPtr("ACTIVE")},
		{"suspended", userpb.UserStatus_USER_STATUS_SUSPENDED, strPtr("SUSPENDED")},
		{"inactive_mapsToDeactivated", userpb.UserStatus_USER_STATUS_INACTIVE, strPtr("DEACTIVATED")},
		{"unspecified_noFilter", userpb.UserStatus_USER_STATUS_UNSPECIFIED, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got *string
			use := &stubUseCase{
				findAllUsers: func(_ context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
					got = opt.Status()
					return 0, nil, nil
				},
			}
			if _, err := NewUserGrpcServer(use).ListUsers(context.Background(), &userpb.ListUsersRequest{Status: c.in}); err != nil {
				t.Fatalf("err: %v", err)
			}
			if (got == nil) != (c.want == nil) || (got != nil && *got != *c.want) {
				t.Errorf("Status = %v, want %v", got, c.want)
			}
		})
	}
}

func TestUserGrpcServer_ListUsers_tokenGiven_pagesAndSetsNextToken(t *testing.T) {
	var gotPage int32
	use := &stubUseCase{
		findAllUsers: func(_ context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
			gotPage = opt.Page().GetPage()
			return 45, []*aggregate.UserAggregate{newAggregate("idp-1", "a@x.com", []string{"USER"})}, nil
		},
	}
	resp, err := NewUserGrpcServer(use).ListUsers(context.Background(), &userpb.ListUsersRequest{PageSize: 20, PageToken: "1"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotPage != 1 {
		t.Errorf("page forwarded = %d, want 1", gotPage)
	}
	if resp.GetNextPageToken() != "2" {
		t.Errorf("next token = %q, want 2", resp.GetNextPageToken())
	}
}

func TestUserGrpcServer_ListUsers_tokenInvalid_invalidArgument(t *testing.T) {
	use := &stubUseCase{
		findAllUsers: func(context.Context, option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
			t.Error("usecase must not be called")
			return 0, nil, nil
		},
	}
	_, err := NewUserGrpcServer(use).ListUsers(context.Background(), &userpb.ListUsersRequest{PageToken: "abc"})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("code = %v, want InvalidArgument", status.Code(err))
	}
}

func strPtr(v string) *string { return &v }
