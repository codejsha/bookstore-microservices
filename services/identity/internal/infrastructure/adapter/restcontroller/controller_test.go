package restcontroller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

type stubUseCase struct {
	registerUser    func(context.Context, command.UserRegisterCommand) (*aggregate.UserAggregate, error)
	findAllUsers    func(context.Context, option.UserQueryOption) (int64, []*aggregate.UserAggregate, error)
	findUser        func(context.Context, string) (*aggregate.UserAggregate, error)
	findUserByEmail func(context.Context, string) (*aggregate.UserAggregate, error)
	updateUser      func(context.Context, string, command.UserUpdateCommand) (*aggregate.UserAggregate, error)
	updateUserRoles func(context.Context, string, command.UserRolesCommand) (*aggregate.UserAggregate, error)
	suspendUser     func(context.Context, string) (*aggregate.UserAggregate, error)
	reactivateUser  func(context.Context, string) (*aggregate.UserAggregate, error)
	deactivateUser  func(context.Context, string) (*aggregate.UserAggregate, error)
	syncUserFromIdp func(context.Context, string) (*aggregate.UserAggregate, error)
}

var _ usecase.IdentityUseCase = (*stubUseCase)(nil)

func (s *stubUseCase) RegisterUser(ctx context.Context, cmd command.UserRegisterCommand) (*aggregate.UserAggregate, error) {
	return s.registerUser(ctx, cmd)
}
func (s *stubUseCase) FindAllUsers(ctx context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
	return s.findAllUsers(ctx, opt)
}
func (s *stubUseCase) FindUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	return s.findUser(ctx, uid)
}
func (s *stubUseCase) FindUserByEmail(ctx context.Context, email string) (*aggregate.UserAggregate, error) {
	return s.findUserByEmail(ctx, email)
}
func (s *stubUseCase) UpdateUser(ctx context.Context, uid string, cmd command.UserUpdateCommand) (*aggregate.UserAggregate, error) {
	return s.updateUser(ctx, uid, cmd)
}
func (s *stubUseCase) UpdateUserRoles(ctx context.Context, uid string, cmd command.UserRolesCommand) (*aggregate.UserAggregate, error) {
	return s.updateUserRoles(ctx, uid, cmd)
}
func (s *stubUseCase) SuspendUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	return s.suspendUser(ctx, uid)
}
func (s *stubUseCase) ReactivateUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	return s.reactivateUser(ctx, uid)
}
func (s *stubUseCase) DeactivateUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	return s.deactivateUser(ctx, uid)
}
func (s *stubUseCase) SyncUserFromIdp(ctx context.Context, idpId string) (*aggregate.UserAggregate, error) {
	return s.syncUserFromIdp(ctx, idpId)
}

func newAggregate(idpId, email, first, last, phone, status string) *aggregate.UserAggregate {
	idp := idpId
	var phonePtr *string
	if phone != "" {
		phonePtr = &phone
	}
	res := &repo.UserResult{
		Id: 1, IdpId: &idp, Email: email, FirstName: first, LastName: last,
		Phone: phonePtr, Roles: []string{"PROFILE", "VIEW"}, Status: status,
		CreatedAt: time.Now(),
	}
	return aggregate.NewUserAggregate(res)
}

// ─── tests ──────────────────────────────────────────────────────────────────

func TestUsersRegister_WhenRequestValid_ReturnsRegisteredUser(t *testing.T) {
	captured := command.UserRegisterCommand{}
	use := &stubUseCase{
		registerUser: func(_ context.Context, cmd command.UserRegisterCommand) (*aggregate.UserAggregate, error) {
			captured = cmd
			return newAggregate("idp-1", cmd.Email, cmd.FirstName, cmd.LastName, "", "ACTIVE"), nil
		},
	}
	ctrl := NewUserController(use)
	roles := []openapi.AuthRole{openapi.AUTHROLE_PROFILE, openapi.AUTHROLE_ORDER}
	req := openapi.UserRegisterRequest{
		Email: "u@x.com", Password: "secret", FirstName: "F", LastName: "L", Roles: &roles,
	}
	resp, err := ctrl.UsersRegister(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Email != "u@x.com" || resp.FirstName != "F" {
		t.Errorf("resp = %+v", resp)
	}
	if captured.Email != "u@x.com" || len(captured.Roles) != 2 {
		t.Errorf("captured = %+v", captured)
	}
	if captured.Roles[0] != "PROFILE" || captured.Roles[1] != "ORDER" {
		t.Errorf("Roles = %+v", captured.Roles)
	}
}

func TestUsersGetAll_WhenFiltersGiven_ReturnsPagedUsers(t *testing.T) {
	use := &stubUseCase{
		findAllUsers: func(_ context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
			if opt.Email() == nil || *opt.Email() != "x@x.com" {
				t.Errorf("Email = %v", opt.Email())
			}
			return 1, []*aggregate.UserAggregate{newAggregate("idp-1", "x@x.com", "F", "L", "", "ACTIVE")}, nil
		},
	}
	ctrl := NewUserController(use)
	email := "x@x.com"
	resp, err := ctrl.UsersGetAll(context.Background(), &email, nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("resp = %+v", resp)
	}
	if resp.Items[0].Email != "x@x.com" || resp.Items[0].Status != openapi.USERSTATUS_ACTIVE {
		t.Errorf("Items[0] = %+v", resp.Items[0])
	}
}

func TestUsersGetAll_WhenUsecaseFails_PropagatesError(t *testing.T) {
	wantErr := errors.New("repo down")
	use := &stubUseCase{
		findAllUsers: func(context.Context, option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
			return 0, nil, wantErr
		},
	}
	ctrl := NewUserController(use)
	_, err := ctrl.UsersGetAll(context.Background(), nil, nil, nil, nil, nil, nil, nil)
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestUsersUpdate_WhenRequestValid_ReturnsUpdatedUser(t *testing.T) {
	use := &stubUseCase{
		updateUser: func(_ context.Context, uid string, cmd command.UserUpdateCommand) (*aggregate.UserAggregate, error) {
			if uid != "uid-1" {
				t.Errorf("uid = %q", uid)
			}
			if cmd.FirstName == nil || *cmd.FirstName != "New" {
				t.Errorf("FirstName = %v", cmd.FirstName)
			}
			return newAggregate("idp-1", "u@x.com", *cmd.FirstName, "L", "", "ACTIVE"), nil
		},
	}
	ctrl := NewUserController(use)
	first := "New"
	resp, err := ctrl.UsersUpdate(context.Background(), "uid-1", openapi.UserUpdateRequest{FirstName: &first})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.FirstName != "New" {
		t.Errorf("FirstName = %q", resp.FirstName)
	}
}

func TestUsersDeactivate_WhenUsecaseSucceeds_ReturnsDeactivatedStatus(t *testing.T) {
	use := &stubUseCase{
		deactivateUser: func(context.Context, string) (*aggregate.UserAggregate, error) {
			return newAggregate("idp-1", "u@x.com", "F", "L", "", "ACTIVE"), nil
		},
	}
	ctrl := NewUserController(use)
	resp, err := ctrl.UsersDeactivate(context.Background(), "uid-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Status != openapi.USERSTATUS_DEACTIVATED {
		t.Errorf("Status = %v, want DEACTIVATED (controller must override)", resp.Status)
	}
}

func TestUsersSuspend_WhenUsecaseSucceeds_ReturnsSuspendedStatus(t *testing.T) {
	use := &stubUseCase{
		suspendUser: func(context.Context, string) (*aggregate.UserAggregate, error) {
			return newAggregate("idp-1", "u@x.com", "F", "L", "", "ACTIVE"), nil
		},
	}
	ctrl := NewUserController(use)
	resp, err := ctrl.UsersSuspend(context.Background(), "uid-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Status != openapi.USERSTATUS_SUSPENDED {
		t.Errorf("Status = %v, want SUSPENDED", resp.Status)
	}
}

func TestUsersReactivate_WhenUsecaseSucceeds_ReturnsActiveStatus(t *testing.T) {
	use := &stubUseCase{
		reactivateUser: func(context.Context, string) (*aggregate.UserAggregate, error) {
			return newAggregate("idp-1", "u@x.com", "F", "L", "", "SUSPENDED"), nil
		},
	}
	ctrl := NewUserController(use)
	resp, err := ctrl.UsersReactivate(context.Background(), "uid-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Status != openapi.USERSTATUS_ACTIVE {
		t.Errorf("Status = %v, want ACTIVE", resp.Status)
	}
}

func TestUsersUpdateRoles_WhenRequestValid_PassesRolesToUsecase(t *testing.T) {
	captured := command.UserRolesCommand{}
	use := &stubUseCase{
		updateUserRoles: func(_ context.Context, _ string, cmd command.UserRolesCommand) (*aggregate.UserAggregate, error) {
			captured = cmd
			return newAggregate("idp-1", "u@x.com", "F", "L", "", "ACTIVE"), nil
		},
	}
	ctrl := NewUserController(use)
	_, err := ctrl.UsersUpdateRoles(context.Background(), "uid-1", openapi.UserRolesRequest{
		Roles: []openapi.AuthRole{openapi.AUTHROLE_MANAGE, openapi.AUTHROLE_VIEW},
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(captured.Roles) != 2 || captured.Roles[0] != "MANAGE" || captured.Roles[1] != "VIEW" {
		t.Errorf("Roles = %+v", captured.Roles)
	}
}

func TestUsersSyncFromIdp_WhenStatusMissing_ReturnsActiveStatus(t *testing.T) {
	use := &stubUseCase{
		syncUserFromIdp: func(_ context.Context, idpId string) (*aggregate.UserAggregate, error) {
			if idpId != "idp-99" {
				t.Errorf("idpId = %q", idpId)
			}
			return newAggregate(idpId, "u@x.com", "F", "L", "", ""), nil
		},
	}
	ctrl := NewUserController(use)
	resp, err := ctrl.UsersSyncFromIdp(context.Background(), "idp-99")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Status != openapi.USERSTATUS_ACTIVE {
		t.Errorf("default status = %v, want ACTIVE (parseStatusOrDefault)", resp.Status)
	}
}

// ─── helpers ───────────────────────────────────────────────────────────────

// The table also covers the known statuses, each of which must map to its own value.
func TestParseStatusOrDefault_WhenStatusBlankOrUnknown_ReturnsActive(t *testing.T) {
	cases := []struct {
		in   string
		want openapi.UserStatus
	}{
		{"", openapi.USERSTATUS_ACTIVE},
		{"BOGUS", openapi.USERSTATUS_ACTIVE},
		{"ACTIVE", openapi.USERSTATUS_ACTIVE},
		{"SUSPENDED", openapi.USERSTATUS_SUSPENDED},
		{"DEACTIVATED", openapi.USERSTATUS_DEACTIVATED},
	}
	for _, c := range cases {
		if got := parseStatusOrDefault(c.in); got != c.want {
			t.Errorf("parseStatusOrDefault(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

// The populated case is asserted too: it must come back as the role names in order.
func TestAuthRolesToStrings_WhenRolesNil_ReturnsNil(t *testing.T) {
	if got := authRolesToStrings(nil); got != nil {
		t.Errorf("nil -> %v, want nil", got)
	}
	in := []openapi.AuthRole{openapi.AUTHROLE_ORDER, openapi.AUTHROLE_VIEW}
	if got := authRolesToStrings(&in); len(got) != 2 || got[0] != "ORDER" || got[1] != "VIEW" {
		t.Errorf("got = %v", got)
	}
}

var _ = pagination.NewPageOption
