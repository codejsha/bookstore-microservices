package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"

	idp "github.com/codejsha/bookstore-microservices/identity/generated/application/port/idpapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

// ─── stubs ──────────────────────────────────────────────────────────────────

type stubUserRepo struct {
	findAllFn       func(ctx context.Context, opt option.UserQueryOption) (int64, []*repo.UserResult, error)
	findByUidFn     func(ctx context.Context, uid string) (*repo.UserResult, error)
	fetchByEmailFn  func(ctx context.Context, email string) ([]*repo.UserResult, error)
	upsertFn        func(ctx context.Context, profile repo.UserUpsert) (*repo.UserResult, error)
	updateProfileFn func(ctx context.Context, idpUid string, patch repo.UserProfileUpdate) error
	updateStatusFn  func(ctx context.Context, idpUid, status string) error
	updateRolesFn   func(ctx context.Context, idpUid string, roles []string) error
	softDeleteFn    func(ctx context.Context, idpUid string) error
}

func (s *stubUserRepo) FindAll(ctx context.Context, opt option.UserQueryOption) (int64, []*repo.UserResult, error) {
	return s.findAllFn(ctx, opt)
}
func (s *stubUserRepo) FindByUid(ctx context.Context, uid string) (*repo.UserResult, error) {
	return s.findByUidFn(ctx, uid)
}
func (s *stubUserRepo) FetchByEmail(ctx context.Context, email string) ([]*repo.UserResult, error) {
	return s.fetchByEmailFn(ctx, email)
}
func (s *stubUserRepo) Upsert(ctx context.Context, profile repo.UserUpsert) (*repo.UserResult, error) {
	if s.upsertFn != nil {
		return s.upsertFn(ctx, profile)
	}
	return nil, nil
}
func (s *stubUserRepo) UpdateProfile(ctx context.Context, idpUid string, patch repo.UserProfileUpdate) error {
	if s.updateProfileFn != nil {
		return s.updateProfileFn(ctx, idpUid, patch)
	}
	return nil
}
func (s *stubUserRepo) UpdateStatus(ctx context.Context, idpUid, status string) error {
	if s.updateStatusFn != nil {
		return s.updateStatusFn(ctx, idpUid, status)
	}
	return nil
}
func (s *stubUserRepo) UpdateRoles(ctx context.Context, idpUid string, roles []string) error {
	if s.updateRolesFn != nil {
		return s.updateRolesFn(ctx, idpUid, roles)
	}
	return nil
}
func (s *stubUserRepo) SoftDelete(ctx context.Context, idpUid string) error {
	if s.softDeleteFn != nil {
		return s.softDeleteFn(ctx, idpUid)
	}
	return nil
}

type stubUsersClient struct {
	listUsersFn         func(ctx context.Context, realm, email string) ([]idp.UserRepresentation, error)
	createUserFn        func(ctx context.Context, realm string, req idp.UserRepresentation) error
	getUserFn           func(ctx context.Context, realm, userId string) (idp.UserRepresentation, error)
	updateUserFn        func(ctx context.Context, realm, userId string, req idp.UserRepresentation) error
	deleteUserFn        func(ctx context.Context, realm, userId string) error
	logoutUserFn        func(ctx context.Context, realm, userId string) error
	setUserRealmRolesFn func(ctx context.Context, realm, userId string, roles []string) error
	logouts             []string
	realmRoles          map[string][]string
}

func (s *stubUsersClient) ListUsers(ctx context.Context, realm, email string) ([]idp.UserRepresentation, error) {
	return s.listUsersFn(ctx, realm, email)
}
func (s *stubUsersClient) CreateUser(ctx context.Context, realm string, req idp.UserRepresentation) error {
	return s.createUserFn(ctx, realm, req)
}
func (s *stubUsersClient) GetUser(ctx context.Context, realm, userId string) (idp.UserRepresentation, error) {
	return s.getUserFn(ctx, realm, userId)
}
func (s *stubUsersClient) UpdateUser(ctx context.Context, realm, userId string, req idp.UserRepresentation) error {
	return s.updateUserFn(ctx, realm, userId, req)
}
func (s *stubUsersClient) DeleteUser(ctx context.Context, realm, userId string) error {
	return s.deleteUserFn(ctx, realm, userId)
}
func (s *stubUsersClient) LogoutUser(ctx context.Context, realm, userId string) error {
	s.logouts = append(s.logouts, userId)
	if s.logoutUserFn != nil {
		return s.logoutUserFn(ctx, realm, userId)
	}
	return nil
}
func (s *stubUsersClient) SetUserRealmRoles(ctx context.Context, realm, userId string, roles []string) error {
	if s.setUserRealmRolesFn != nil {
		if err := s.setUserRealmRolesFn(ctx, realm, userId, roles); err != nil {
			return err
		}
	}
	if s.realmRoles == nil {
		s.realmRoles = map[string][]string{}
	}
	s.realmRoles[userId] = roles
	return nil
}

type stubRevoker struct {
	revokeFn func(ctx context.Context, sub string) error
	revoked  []string
}

func (s *stubRevoker) Revoke(ctx context.Context, sub string) error {
	s.revoked = append(s.revoked, sub)
	if s.revokeFn != nil {
		return s.revokeFn(ctx, sub)
	}
	return nil
}

func newSvc(userRepo *stubUserRepo, usersClient *stubUsersClient) *identityService {
	return newSvcWithRevoker(userRepo, usersClient, &stubRevoker{})
}

func newSvcWithRevoker(
	userRepo *stubUserRepo,
	usersClient *stubUsersClient,
	revoker *stubRevoker,
) *identityService {
	if userRepo == nil {
		userRepo = &stubUserRepo{}
	}
	if usersClient == nil {
		usersClient = &stubUsersClient{}
	}
	return &identityService{
		cfg:         &config.Config{Keycloak: &sharedconfig.KeycloakConfig{Realm: "test-realm"}},
		userRepo:    userRepo,
		usersClient: usersClient,
		revoker:     revoker,
	}
}

func ptrStr(s string) *string { return &s }

// ─── RegisterUser ───────────────────────────────────────────────────────────

func TestRegisterUser_Success(t *testing.T) {
	createCalled := false
	listCount := 0
	users := &stubUsersClient{
		listUsersFn: func(_ context.Context, realm, email string) ([]idp.UserRepresentation, error) {
			listCount++
			if realm != "test-realm" || email != "u@example.com" {
				t.Errorf("ListUsers args = %q/%q", realm, email)
			}
			if listCount == 1 {
				return nil, nil
			}
			id := "idp-id-1"
			fn := "Alice"
			ln := "Smith"
			roles := []string{"PROFILE", "ORDER", "VIEW"}
			return []idp.UserRepresentation{{Id: &id, Email: &email, FirstName: &fn, LastName: &ln, RealmRoles: &roles}}, nil
		},
		createUserFn: func(_ context.Context, realm string, req idp.UserRepresentation) error {
			createCalled = true
			if req.Username == nil || *req.Username != "u@example.com" {
				t.Errorf("CreateUser username = %v", req.Username)
			}
			if req.RealmRoles != nil {
				t.Errorf("CreateUser RealmRoles = %v, want nil", req.RealmRoles)
			}
			if req.Credentials == nil || len(*req.Credentials) != 1 || *(*req.Credentials)[0].Value != "secret" {
				t.Errorf("CreateUser credentials malformed: %+v", req.Credentials)
			}
			return nil
		},
	}
	svc := newSvc(nil, users)

	agg, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{
		Email: "u@example.com", Password: "secret", FirstName: "Alice", LastName: "Smith",
	})
	if err != nil {
		t.Fatalf("RegisterUser err: %v", err)
	}
	if !createCalled {
		t.Error("CreateUser was not invoked")
	}
	if listCount != 2 {
		t.Errorf("ListUsers calls = %d, want 2 (pre + post)", listCount)
	}
	if got := users.realmRoles["idp-id-1"]; len(got) != 3 {
		t.Errorf("bound realm roles = %v, want the 3 defaults", got)
	}
	if agg.Email() != "u@example.com" || agg.FirstName() != "Alice" || agg.LastName() != "Smith" {
		t.Errorf("agg fields = email=%q first=%q last=%q", agg.Email(), agg.FirstName(), agg.LastName())
	}
}

func TestRegisterUser_AlreadyExists(t *testing.T) {
	users := &stubUsersClient{
		listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) {
			return []idp.UserRepresentation{{Id: ptrStr("existing")}}, nil
		},
		createUserFn: func(context.Context, string, idp.UserRepresentation) error {
			t.Error("CreateUser should not be called when user already exists")
			return nil
		},
	}
	svc := newSvc(nil, users)
	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{Email: "dup@x.com"})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("err = %v, want already-exists message", err)
	}
}

func TestRegisterUser_CreateFails(t *testing.T) {
	users := &stubUsersClient{
		listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) { return nil, nil },
		createUserFn: func(context.Context, string, idp.UserRepresentation) error {
			return errors.New("kc down")
		},
	}
	svc := newSvc(nil, users)
	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{Email: "x@x.com"})
	if err == nil || !strings.Contains(err.Error(), "failed to create user in IdP") {
		t.Errorf("err = %v, want create wrap", err)
	}
}

func TestRegisterUser_NotFoundAfterCreate(t *testing.T) {
	calls := 0
	users := &stubUsersClient{
		listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) {
			calls++
			return nil, nil
		},
		createUserFn: func(context.Context, string, idp.UserRepresentation) error { return nil },
	}
	svc := newSvc(nil, users)
	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{Email: "ghost@x.com"})
	if err == nil || !strings.Contains(err.Error(), "not found after creation") {
		t.Errorf("err = %v, want post-create not-found", err)
	}
	if calls != 3 {
		t.Errorf("ListUsers calls = %d, want 3 (pre + post + compensation re-list)", calls)
	}
}

func TestRegisterUser_NotFoundAfterCreateCompensatesLocatedOrphan(t *testing.T) {
	id := "orphan-id"
	em := "ghost@x.com"
	calls := 0
	var deleted string
	users := &stubUsersClient{
		listUsersFn: func(_ context.Context, _ string, _ string) ([]idp.UserRepresentation, error) {
			calls++
			switch calls {
			case 1, 2:
				return nil, nil
			default:
				return []idp.UserRepresentation{{Id: &id, Email: &em}}, nil
			}
		},
		createUserFn: func(context.Context, string, idp.UserRepresentation) error { return nil },
		deleteUserFn: func(_ context.Context, _ string, userId string) error {
			deleted = userId
			return nil
		},
	}
	svc := newSvc(nil, users)
	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{Email: em})
	if err == nil || !strings.Contains(err.Error(), "not found after creation") {
		t.Errorf("err = %v, want post-create not-found", err)
	}
	if deleted != id {
		t.Errorf("compensating DeleteUser userId = %q, want %q", deleted, id)
	}
}

func TestRegisterUser_RespectsExplicitRoles(t *testing.T) {
	var captured idp.UserRepresentation
	id := "idp-1"
	fn, ln := "F", "L"
	em := "u@x.com"
	users := &stubUsersClient{
		listUsersFn: func(_ context.Context, _ string, email string) ([]idp.UserRepresentation, error) {
			if email == em {
				return nil, nil
			}
			return []idp.UserRepresentation{{Id: &id, FirstName: &fn, LastName: &ln, Email: &em}}, nil
		},
		createUserFn: func(_ context.Context, _ string, req idp.UserRepresentation) error {
			captured = req
			return nil
		},
	}
	calls := 0
	users.listUsersFn = func(context.Context, string, string) ([]idp.UserRepresentation, error) {
		calls++
		if calls == 1 {
			return nil, nil
		}
		return []idp.UserRepresentation{{Id: &id, FirstName: &fn, LastName: &ln, Email: &em}}, nil
	}
	svc := newSvc(nil, users)
	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{
		Email: em, Password: "p", FirstName: fn, LastName: ln, Roles: []string{"VIEW"},
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.RealmRoles != nil {
		t.Errorf("CreateUser RealmRoles = %v, want nil (Keycloak ignores it)", captured.RealmRoles)
	}
	got := users.realmRoles[id]
	if len(got) != 1 || got[0] != "VIEW" {
		t.Errorf("bound realm roles = %v, want [VIEW]", got)
	}
}

func TestRegisterUser_RejectsElevatedRoles(t *testing.T) {
	for _, role := range []string{"MANAGE", "SYSTEM", "ADMIN", "manage"} {
		t.Run(role, func(t *testing.T) {
			created := false
			users := &stubUsersClient{
				listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) {
					return nil, nil
				},
				createUserFn: func(context.Context, string, idp.UserRepresentation) error {
					created = true
					return nil
				},
			}
			svc := newSvc(nil, users)

			_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{
				Email: "u@x.com", Password: "p", FirstName: "F", LastName: "L",
				Roles: []string{role},
			})

			if !errors.Is(err, ErrElevatedRoleOnRegister) {
				t.Fatalf("RegisterUser(roles=[%s]) err = %v, want ErrElevatedRoleOnRegister", role, err)
			}
			if created {
				t.Error("Keycloak user was created even though the role was refused")
			}
		})
	}
}

func TestRegisterUser_RejectsElevatedRoleMixedWithBaseline(t *testing.T) {
	users := &stubUsersClient{
		listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) {
			return nil, nil
		},
	}
	svc := newSvc(nil, users)

	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{
		Email: "u@x.com", Password: "p", FirstName: "F", LastName: "L",
		Roles: []string{"VIEW", "PROFILE", "SYSTEM"},
	})

	if !errors.Is(err, ErrElevatedRoleOnRegister) {
		t.Fatalf("RegisterUser err = %v, want ErrElevatedRoleOnRegister", err)
	}
}

func TestRegisterUser_RoleBindingFailureDeletesIdpUser(t *testing.T) {
	id := "idp-1"
	em := "u@x.com"
	calls := 0
	deleted := ""
	users := &stubUsersClient{
		createUserFn: func(context.Context, string, idp.UserRepresentation) error { return nil },
		setUserRealmRolesFn: func(context.Context, string, string, []string) error {
			return errors.New("role ADMIN does not exist in realm test-realm")
		},
		deleteUserFn: func(_ context.Context, _ string, userId string) error {
			deleted = userId
			return nil
		},
	}
	users.listUsersFn = func(context.Context, string, string) ([]idp.UserRepresentation, error) {
		calls++
		if calls == 1 {
			return nil, nil
		}
		return []idp.UserRepresentation{{Id: &id, Email: &em}}, nil
	}

	svc := newSvc(nil, users)
	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{
		Email: em, Password: "p", FirstName: "F", LastName: "L", Roles: []string{"VIEW"},
	})
	if err == nil {
		t.Fatal("RegisterUser succeeded, want failure when realm roles cannot be bound")
	}
	if deleted != id {
		t.Errorf("compensating DeleteUser userId = %q, want %q", deleted, id)
	}
}

func TestRegisterUser_ResponseCarriesRolesAndTimestamps(t *testing.T) {
	id := "idp-id-1"
	em := "u@example.com"
	listCount := 0
	users := &stubUsersClient{
		listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) {
			listCount++
			if listCount == 1 {
				return nil, nil
			}
			return []idp.UserRepresentation{{Id: &id, Email: &em}}, nil
		},
		createUserFn: func(context.Context, string, idp.UserRepresentation) error { return nil },
	}
	svc := newSvc(nil, users)

	agg, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{
		Email: em, Password: "secret", FirstName: "Alice", LastName: "Smith",
		Roles: []string{"VIEW"},
	})
	if err != nil {
		t.Fatalf("RegisterUser err: %v", err)
	}
	if got := agg.Roles(); len(got) != 1 {
		t.Errorf("response roles = %v, want the single granted VIEW role", got)
	}
	if agg.CreatedAt().IsZero() {
		t.Error("response createdAt is zero; want the synced/local timestamp")
	}
}

// ─── Query ──────────────────────────────────────────────────────────────────

func TestFindAllUsers(t *testing.T) {
	results := []*repo.UserResult{
		{Id: 1, Email: "a@x.com", FirstName: "A", LastName: "B", Roles: []string{"VIEW"}},
		{Id: 2, Email: "b@x.com", FirstName: "C", LastName: "D", Roles: []string{}},
	}
	userRepo := &stubUserRepo{
		findAllFn: func(context.Context, option.UserQueryOption) (int64, []*repo.UserResult, error) {
			return 2, results, nil
		},
	}
	svc := newSvc(userRepo, nil)
	total, aggs, err := svc.FindAllUsers(context.Background(), option.NewUserQueryOption())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if total != 2 || len(aggs) != 2 {
		t.Fatalf("got total=%d len=%d", total, len(aggs))
	}
	if aggs[0].Email() != "a@x.com" {
		t.Errorf("aggs[0].Email = %q", aggs[0].Email())
	}
}

func TestFindAllUsers_RepoError(t *testing.T) {
	want := errors.New("db error")
	svc := newSvc(&stubUserRepo{
		findAllFn: func(context.Context, option.UserQueryOption) (int64, []*repo.UserResult, error) {
			return 0, nil, want
		},
	}, nil)
	_, _, err := svc.FindAllUsers(context.Background(), option.NewUserQueryOption())
	if !errors.Is(err, want) {
		t.Errorf("err = %v, want %v", err, want)
	}
}

func TestFindUserByEmail(t *testing.T) {
	userRepo := &stubUserRepo{
		fetchByEmailFn: func(_ context.Context, email string) ([]*repo.UserResult, error) {
			if email != "x@x.com" {
				t.Errorf("email = %q", email)
			}
			return []*repo.UserResult{{Id: 1, Email: email}}, nil
		},
	}
	svc := newSvc(userRepo, nil)
	got, err := svc.FindUserByEmail(context.Background(), "x@x.com")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got.Email() != "x@x.com" {
		t.Errorf("email = %q", got.Email())
	}
}

func TestFindUserByEmail_NotFound(t *testing.T) {
	userRepo := &stubUserRepo{
		fetchByEmailFn: func(context.Context, string) ([]*repo.UserResult, error) { return nil, nil },
	}
	svc := newSvc(userRepo, nil)
	_, err := svc.FindUserByEmail(context.Background(), "nope@x.com")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("err = %v, want not-found", err)
	}
}

// ─── User management ────────────────────────────────────────────────────────

func TestUpdateUser_AppliesNameChanges(t *testing.T) {
	idpId := "idp-1"
	captured := ""
	user := &repo.UserResult{Id: 1, IdpId: &idpId, Email: "u@x.com", FirstName: "Old", LastName: "Name"}
	usersClient := &stubUsersClient{
		getUserFn: func(_ context.Context, realm, userId string) (idp.UserRepresentation, error) {
			if userId != idpId {
				t.Errorf("getUser id = %q, want %q", userId, idpId)
			}
			fn := "Old"
			ln := "Name"
			return idp.UserRepresentation{Id: &userId, FirstName: &fn, LastName: &ln}, nil
		},
		updateUserFn: func(_ context.Context, realm, userId string, req idp.UserRepresentation) error {
			if req.FirstName == nil || req.LastName == nil {
				t.Errorf("update missing names: %+v", req)
				return nil
			}
			captured = *req.FirstName + "/" + *req.LastName
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
	}
	svc := newSvc(userRepo, usersClient)
	first := "New"
	last := "Person"
	_, err := svc.UpdateUser(context.Background(), "uid", command.UserUpdateCommand{FirstName: &first, LastName: &last})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured != "New/Person" {
		t.Errorf("captured = %q, want New/Person", captured)
	}
}

func TestUpdateUser_NoIdpIdSkipsKeycloak(t *testing.T) {
	user := &repo.UserResult{Id: 1, IdpId: nil, Email: "x@x.com"}
	usersClient := &stubUsersClient{
		getUserFn: func(context.Context, string, string) (idp.UserRepresentation, error) {
			t.Error("GetUser should not be called when IdpId is nil")
			return idp.UserRepresentation{}, nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
	}
	svc := newSvc(userRepo, usersClient)
	first := "Z"
	_, err := svc.UpdateUser(context.Background(), "uid", command.UserUpdateCommand{FirstName: &first})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestSuspendUser(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId}
	called := false
	usersClient := &stubUsersClient{
		updateUserFn: func(_ context.Context, _ string, _ string, req idp.UserRepresentation) error {
			called = true
			if req.Enabled == nil || *req.Enabled {
				t.Errorf("Enabled = %v, want false", req.Enabled)
			}
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
	}
	svc := newSvc(userRepo, usersClient)
	_, err := svc.SuspendUser(context.Background(), "uid")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !called {
		t.Error("UpdateUser was not called")
	}
}

func TestReactivateUser(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId}
	var enabled *bool
	usersClient := &stubUsersClient{
		updateUserFn: func(_ context.Context, _ string, _ string, req idp.UserRepresentation) error {
			enabled = req.Enabled
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
	}
	svc := newSvc(userRepo, usersClient)
	_, err := svc.ReactivateUser(context.Background(), "uid")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if enabled == nil || !*enabled {
		t.Errorf("Enabled = %v, want true", enabled)
	}
}

func TestDeactivateUser(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId}
	deletedId := ""
	usersClient := &stubUsersClient{
		deleteUserFn: func(_ context.Context, _ string, userId string) error {
			deletedId = userId
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
	}
	svc := newSvc(userRepo, usersClient)
	_, err := svc.DeactivateUser(context.Background(), "uid")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if deletedId != idpId {
		t.Errorf("DeleteUser id = %q, want %q", deletedId, idpId)
	}
}

func TestDeactivateUser_KeycloakError(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId}
	usersClient := &stubUsersClient{
		deleteUserFn: func(context.Context, string, string) error { return errors.New("kc fail") },
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
	}
	svc := newSvc(userRepo, usersClient)
	_, err := svc.DeactivateUser(context.Background(), "uid")
	if err == nil || !strings.Contains(err.Error(), "failed to deactivate user in IdP") {
		t.Errorf("err = %v, want deactivate wrap", err)
	}
}

func TestSyncUserFromIdp(t *testing.T) {
	id := "idp-1"
	em := "u@x.com"
	fn, ln := "Alice", "Lee"
	usersClient := &stubUsersClient{
		getUserFn: func(_ context.Context, _ string, userId string) (idp.UserRepresentation, error) {
			if userId != id {
				t.Errorf("userId = %q, want %q", userId, id)
			}
			roles := []string{"PROFILE"}
			return idp.UserRepresentation{Id: &userId, Email: &em, FirstName: &fn, LastName: &ln, RealmRoles: &roles}, nil
		},
	}
	svc := newSvc(nil, usersClient)
	got, err := svc.SyncUserFromIdp(context.Background(), id)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got.Email() != em || got.FirstName() != fn || got.LastName() != ln {
		t.Errorf("agg = email=%q first=%q last=%q", got.Email(), got.FirstName(), got.LastName())
	}
}

func TestSyncUserFromIdp_KeycloakError(t *testing.T) {
	usersClient := &stubUsersClient{
		getUserFn: func(context.Context, string, string) (idp.UserRepresentation, error) {
			return idp.UserRepresentation{}, errors.New("kc dead")
		},
	}
	svc := newSvc(nil, usersClient)
	_, err := svc.SyncUserFromIdp(context.Background(), "idp-x")
	if err == nil || !strings.Contains(err.Error(), "failed to get user from IdP") {
		t.Errorf("err = %v, want get-user wrap", err)
	}
}

// ─── Local DB sync ─────────────────────────────────────────────────────────

func TestRegisterUser_SyncsToLocalDB(t *testing.T) {
	id := "idp-1"
	em := "u@x.com"
	fn, ln := "Alice", "Smith"
	roles := []string{"VIEW"}
	calls := 0
	users := &stubUsersClient{
		listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) {
			calls++
			if calls == 1 {
				return nil, nil
			}
			return []idp.UserRepresentation{{Id: &id, Email: &em, FirstName: &fn, LastName: &ln, RealmRoles: &roles}}, nil
		},
		createUserFn: func(context.Context, string, idp.UserRepresentation) error { return nil },
	}
	var captured repo.UserUpsert
	userRepo := &stubUserRepo{
		upsertFn: func(_ context.Context, p repo.UserUpsert) (*repo.UserResult, error) {
			captured = p
			return &repo.UserResult{Id: 99, IdpId: &p.IdpUid, Email: p.Email}, nil
		},
	}
	svc := newSvc(userRepo, users)
	if _, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{
		Email: em, Password: "p", FirstName: fn, LastName: ln,
	}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.IdpUid != id || captured.Email != em || captured.FirstName != fn || captured.Status != "ACTIVE" {
		t.Errorf("captured = %+v", captured)
	}
}

func TestRegisterUser_LocalSyncFailsBubblesUp(t *testing.T) {
	id := "idp-1"
	em := "u@x.com"
	calls := 0
	compensatedId := ""
	users := &stubUsersClient{
		listUsersFn: func(context.Context, string, string) ([]idp.UserRepresentation, error) {
			calls++
			if calls == 1 {
				return nil, nil
			}
			return []idp.UserRepresentation{{Id: &id, Email: &em}}, nil
		},
		createUserFn: func(context.Context, string, idp.UserRepresentation) error { return nil },
		deleteUserFn: func(_ context.Context, _ string, userId string) error {
			compensatedId = userId
			return nil
		},
	}
	userRepo := &stubUserRepo{
		upsertFn: func(context.Context, repo.UserUpsert) (*repo.UserResult, error) {
			return nil, errors.New("db down")
		},
	}
	svc := newSvc(userRepo, users)
	_, err := svc.RegisterUser(context.Background(), command.UserRegisterCommand{Email: em})
	if err == nil || !strings.Contains(err.Error(), "failed to sync user to local DB") {
		t.Errorf("err = %v, want local-sync wrap", err)
	}
	if compensatedId != id {
		t.Errorf("compensation DeleteUser id = %q, want %q", compensatedId, id)
	}
}

func TestUpdateUser_PropagatesToLocalDB(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId, Email: "u@x.com"}
	usersClient := &stubUsersClient{
		getUserFn: func(_ context.Context, _ string, userId string) (idp.UserRepresentation, error) {
			return idp.UserRepresentation{Id: &userId}, nil
		},
		updateUserFn: func(context.Context, string, string, idp.UserRepresentation) error { return nil },
	}
	var capturedIdp string
	var capturedPatch repo.UserProfileUpdate
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
		updateProfileFn: func(_ context.Context, idp string, p repo.UserProfileUpdate) error {
			capturedIdp = idp
			capturedPatch = p
			return nil
		},
	}
	svc := newSvc(userRepo, usersClient)
	first := "New"
	if _, err := svc.UpdateUser(context.Background(), "uid", command.UserUpdateCommand{FirstName: &first}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if capturedIdp != idpId {
		t.Errorf("UpdateProfile idpUid = %q, want %q", capturedIdp, idpId)
	}
	if capturedPatch.FirstName == nil || *capturedPatch.FirstName != "New" {
		t.Errorf("UpdateProfile FirstName = %v", capturedPatch.FirstName)
	}
}

func TestUpdateUserRoles_UpdatesKeycloakAndLocalDB(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId, Email: "u@x.com"}
	var capturedDbRoles []string
	usersClient := &stubUsersClient{
		updateUserFn: func(_ context.Context, _ string, _ string, _ idp.UserRepresentation) error {
			t.Error("UpdateUser must not be used to bind realm roles: Keycloak ignores realmRoles there")
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
		updateRolesFn: func(_ context.Context, idp string, roles []string) error {
			if idp != idpId {
				t.Errorf("UpdateRoles idpUid = %q", idp)
			}
			capturedDbRoles = roles
			return nil
		},
	}
	svc := newSvc(userRepo, usersClient)
	roles := []string{"MANAGE", "ORDER"}
	if _, err := svc.UpdateUserRoles(context.Background(), "uid", command.UserRolesCommand{Roles: roles}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := usersClient.realmRoles[idpId]; len(got) != 2 || got[0] != "MANAGE" {
		t.Errorf("Keycloak roles = %v", got)
	}
	if len(capturedDbRoles) != 2 || capturedDbRoles[1] != "ORDER" {
		t.Errorf("DB roles = %v", capturedDbRoles)
	}
}

func TestSuspendUser_UpdatesStatusInLocalDB(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId}
	var capturedStatus string
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
		updateStatusFn: func(_ context.Context, _ string, status string) error {
			capturedStatus = status
			return nil
		},
	}
	usersClient := &stubUsersClient{
		updateUserFn: func(context.Context, string, string, idp.UserRepresentation) error { return nil },
	}
	svc := newSvc(userRepo, usersClient)
	if _, err := svc.SuspendUser(context.Background(), "uid"); err != nil {
		t.Fatalf("err: %v", err)
	}
	if capturedStatus != "SUSPENDED" {
		t.Errorf("status = %q, want SUSPENDED", capturedStatus)
	}
}

func TestReactivateUser_UpdatesStatusInLocalDB(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId}
	var capturedStatus string
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) { return user, nil },
		updateStatusFn: func(_ context.Context, _ string, status string) error {
			capturedStatus = status
			return nil
		},
	}
	usersClient := &stubUsersClient{
		updateUserFn: func(context.Context, string, string, idp.UserRepresentation) error { return nil },
	}
	svc := newSvc(userRepo, usersClient)
	if _, err := svc.ReactivateUser(context.Background(), "uid"); err != nil {
		t.Fatalf("err: %v", err)
	}
	if capturedStatus != "ACTIVE" {
		t.Errorf("status = %q, want ACTIVE", capturedStatus)
	}
}

func TestDeactivateUser_SoftDeletesLocalRow(t *testing.T) {
	idpId := "idp-1"
	user := &repo.UserResult{Id: 1, IdpId: &idpId, Status: "ACTIVE"}
	deleted := ""
	userRepo := &stubUserRepo{
		findByUidFn:  func(context.Context, string) (*repo.UserResult, error) { return user, nil },
		softDeleteFn: func(_ context.Context, idp string) error { deleted = idp; return nil },
	}
	usersClient := &stubUsersClient{
		deleteUserFn: func(context.Context, string, string) error { return nil },
	}
	svc := newSvc(userRepo, usersClient)
	got, err := svc.DeactivateUser(context.Background(), "uid")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if deleted != idpId {
		t.Errorf("SoftDelete idpUid = %q", deleted)
	}
	if got.Status() != "DEACTIVATED" {
		t.Errorf("returned aggregate Status = %q, want DEACTIVATED", got.Status())
	}
}

func TestSyncUserFromIdp_UpsertsLocalDB(t *testing.T) {
	id := "idp-1"
	em := "u@x.com"
	enabled := false
	fn, ln := "A", "B"
	roles := []string{"VIEW"}
	usersClient := &stubUsersClient{
		getUserFn: func(_ context.Context, _ string, userId string) (idp.UserRepresentation, error) {
			return idp.UserRepresentation{Id: &userId, Email: &em, FirstName: &fn, LastName: &ln, RealmRoles: &roles, Enabled: &enabled}, nil
		},
	}
	var captured repo.UserUpsert
	userRepo := &stubUserRepo{
		upsertFn: func(_ context.Context, p repo.UserUpsert) (*repo.UserResult, error) {
			captured = p
			return nil, nil
		},
	}
	svc := newSvc(userRepo, usersClient)
	if _, err := svc.SyncUserFromIdp(context.Background(), id); err != nil {
		t.Fatalf("err: %v", err)
	}
	if captured.IdpUid != id || captured.Email != em || captured.Status != "SUSPENDED" {
		t.Errorf("captured = %+v (Enabled=false should map to SUSPENDED)", captured)
	}
}

// ─── session revocation ─────────────────────────────────────────────────────

func TestSuspendUser_RevokesSessions(t *testing.T) {
	idpId := "idp-1"
	usersClient := &stubUsersClient{
		updateUserFn: func(context.Context, string, string, idp.UserRepresentation) error { return nil },
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId}, nil
		},
	}

	revoker := &stubRevoker{}
	if _, err := newSvcWithRevoker(userRepo, usersClient, revoker).SuspendUser(context.Background(), "uid"); err != nil {
		t.Fatalf("SuspendUser: %v", err)
	}
	if len(usersClient.logouts) != 1 || usersClient.logouts[0] != idpId {
		t.Fatalf("logouts = %v, want [%q]", usersClient.logouts, idpId)
	}
	if len(revoker.revoked) != 1 || revoker.revoked[0] != idpId {
		t.Fatalf("revoked = %v, want [%q]", revoker.revoked, idpId)
	}
}

func TestUpdateUserRoles_RevokesSessions(t *testing.T) {
	idpId := "idp-1"
	usersClient := &stubUsersClient{
		getUserFn: func(context.Context, string, string) (idp.UserRepresentation, error) {
			return idp.UserRepresentation{}, nil
		},
		updateUserFn: func(context.Context, string, string, idp.UserRepresentation) error { return nil },
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId}, nil
		},
	}

	revoker := &stubRevoker{}
	_, err := newSvcWithRevoker(userRepo, usersClient, revoker).
		UpdateUserRoles(context.Background(), "uid", command.UserRolesCommand{Roles: []string{"USER"}})
	if err != nil {
		t.Fatalf("UpdateUserRoles: %v", err)
	}
	if len(usersClient.logouts) != 1 || usersClient.logouts[0] != idpId {
		t.Fatalf("logouts = %v, want [%q]", usersClient.logouts, idpId)
	}
	if len(revoker.revoked) != 1 || revoker.revoked[0] != idpId {
		t.Fatalf("revoked = %v, want [%q]", revoker.revoked, idpId)
	}
}

func TestSuspendUser_DenylistFailureAbortsBeforeLocalWrite(t *testing.T) {
	idpId := "idp-1"
	usersClient := &stubUsersClient{
		updateUserFn: func(context.Context, string, string, idp.UserRepresentation) error { return nil },
	}
	revoker := &stubRevoker{
		revokeFn: func(context.Context, string) error { return errors.New("valkey down") },
	}
	statusWritten := false
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId}, nil
		},
		updateStatusFn: func(context.Context, string, string) error {
			statusWritten = true
			return nil
		},
	}

	svc := newSvcWithRevoker(userRepo, usersClient, revoker)
	if _, err := svc.SuspendUser(context.Background(), "uid"); err == nil {
		t.Fatal("expected the denylist failure to surface")
	}
	if statusWritten {
		t.Error("local status was written despite the denylist write failing")
	}
}

func TestSuspendUser_LogoutFailureAbortsBeforeLocalWrite(t *testing.T) {
	idpId := "idp-1"
	usersClient := &stubUsersClient{
		updateUserFn: func(context.Context, string, string, idp.UserRepresentation) error { return nil },
		logoutUserFn: func(context.Context, string, string) error { return errors.New("keycloak down") },
	}
	statusWritten := false
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId}, nil
		},
		updateStatusFn: func(context.Context, string, string) error {
			statusWritten = true
			return nil
		},
	}

	if _, err := newSvc(userRepo, usersClient).SuspendUser(context.Background(), "uid"); err == nil {
		t.Fatal("expected the logout failure to surface")
	}
	if statusWritten {
		t.Error("local status was written despite the sessions surviving")
	}
}

func TestDeactivateUser_DoesNotLogoutSeparately(t *testing.T) {
	idpId := "idp-1"
	usersClient := &stubUsersClient{
		deleteUserFn: func(context.Context, string, string) error { return nil },
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId}, nil
		},
	}

	if _, err := newSvc(userRepo, usersClient).DeactivateUser(context.Background(), "uid"); err != nil {
		t.Fatalf("DeactivateUser: %v", err)
	}
	if len(usersClient.logouts) != 0 {
		t.Fatalf("logouts = %v, want none", usersClient.logouts)
	}
}

// ─── write-order divergence (L2) ─────────────────────────────────────────────

func TestSuspendUser_LocalWriteRetriedOnceThenSucceeds(t *testing.T) {
	idpId := "idp-1"
	disableCalls := 0
	usersClient := &stubUsersClient{
		updateUserFn: func(_ context.Context, _ string, _ string, req idp.UserRepresentation) error {
			disableCalls++
			if req.Enabled == nil || *req.Enabled {
				t.Errorf("Keycloak update Enabled = %v, want false (disable only, never re-enable)", req.Enabled)
			}
			return nil
		},
	}
	statusCalls := 0
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId}, nil
		},
		updateStatusFn: func(context.Context, string, string) error {
			statusCalls++
			if statusCalls == 1 {
				return errors.New("transient db error")
			}
			return nil
		},
	}
	if _, err := newSvc(userRepo, usersClient).SuspendUser(context.Background(), "uid"); err != nil {
		t.Fatalf("SuspendUser: %v", err)
	}
	if statusCalls != 2 {
		t.Errorf("UpdateStatus calls = %d, want 2 (fail then retry-success)", statusCalls)
	}
	if disableCalls != 1 {
		t.Errorf("Keycloak update calls = %d, want 1 (disable only, no re-enable)", disableCalls)
	}
}

func TestSuspendUser_LocalWriteFailsTwiceLeavesKeycloakDisabled(t *testing.T) {
	idpId := "idp-1"
	disableCalls := 0
	reEnabled := false
	usersClient := &stubUsersClient{
		updateUserFn: func(_ context.Context, _ string, _ string, req idp.UserRepresentation) error {
			if req.Enabled != nil && *req.Enabled {
				reEnabled = true
			} else {
				disableCalls++
			}
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId}, nil
		},
		updateStatusFn: func(context.Context, string, string) error {
			return errors.New("db down")
		},
	}
	_, err := newSvc(userRepo, usersClient).SuspendUser(context.Background(), "uid")
	if err == nil || !strings.Contains(err.Error(), "please retry") {
		t.Fatalf("err = %v, want a retry-instructing error", err)
	}
	if reEnabled {
		t.Error("Keycloak was re-enabled after a failed local write; locked-out is the safe direction")
	}
	if disableCalls != 1 {
		t.Errorf("Keycloak disable calls = %d, want 1", disableCalls)
	}
}

func TestUpdateUserRoles_LocalWriteFailureCompensatesKeycloak(t *testing.T) {
	idpId := "idp-1"
	setCalls := 0
	usersClient := &stubUsersClient{
		setUserRealmRolesFn: func(_ context.Context, _ string, _ string, _ []string) error {
			setCalls++
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId, Roles: []string{"VIEW"}}, nil
		},
		updateRolesFn: func(context.Context, string, []string) error {
			return errors.New("db down")
		},
	}
	_, err := newSvc(userRepo, usersClient).
		UpdateUserRoles(context.Background(), "uid", command.UserRolesCommand{Roles: []string{"MANAGE", "ORDER"}})
	if err == nil || !strings.Contains(err.Error(), "failed to update roles in local DB") {
		t.Fatalf("err = %v, want local-DB wrap", err)
	}
	if setCalls != 2 {
		t.Errorf("SetUserRealmRoles calls = %d, want 2 (apply + compensate)", setCalls)
	}
	if got := usersClient.realmRoles[idpId]; len(got) != 1 || got[0] != "VIEW" {
		t.Errorf("post-compensation Keycloak roles = %v, want [VIEW] (restored)", got)
	}
}

func TestUpdateUser_LocalWriteFailureCompensatesKeycloak(t *testing.T) {
	idpId := "idp-1"
	oldFirst, oldLast := "Old", "Name"
	var updates [][2]string
	usersClient := &stubUsersClient{
		getUserFn: func(_ context.Context, _ string, userId string) (idp.UserRepresentation, error) {
			return idp.UserRepresentation{Id: &userId, FirstName: &oldFirst, LastName: &oldLast}, nil
		},
		updateUserFn: func(_ context.Context, _ string, _ string, req idp.UserRepresentation) error {
			f, l := "", ""
			if req.FirstName != nil {
				f = *req.FirstName
			}
			if req.LastName != nil {
				l = *req.LastName
			}
			updates = append(updates, [2]string{f, l})
			return nil
		},
	}
	userRepo := &stubUserRepo{
		findByUidFn: func(context.Context, string) (*repo.UserResult, error) {
			return &repo.UserResult{Id: 1, IdpId: &idpId, Email: "u@x.com"}, nil
		},
		updateProfileFn: func(context.Context, string, repo.UserProfileUpdate) error {
			return errors.New("db down")
		},
	}
	newFirst := "New"
	_, err := newSvc(userRepo, usersClient).
		UpdateUser(context.Background(), "uid", command.UserUpdateCommand{FirstName: &newFirst})
	if err == nil || !strings.Contains(err.Error(), "failed to update user in local DB") {
		t.Fatalf("err = %v, want local-DB wrap", err)
	}
	if len(updates) != 2 {
		t.Fatalf("Keycloak update calls = %v, want 2 (apply + compensate)", updates)
	}
	if updates[0] != [2]string{"New", "Name"} {
		t.Errorf("first update = %v, want [New Name]", updates[0])
	}
	if updates[1] != [2]string{"Old", "Name"} {
		t.Errorf("compensation update = %v, want the previous [Old Name]", updates[1])
	}
}
