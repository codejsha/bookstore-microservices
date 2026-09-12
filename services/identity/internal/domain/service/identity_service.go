package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/ptr"

	idp "github.com/codejsha/bookstore-microservices/identity/generated/application/port/idpapi"
	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/usecase"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/aggregate"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

const (
	statusActive      = string(constant.USERSTATUS_ACTIVE_VALUE)
	statusSuspended   = string(constant.USERSTATUS_SUSPENDED_VALUE)
	statusDeactivated = string(constant.USERSTATUS_DEACTIVATED_VALUE)
)

var _ usecase.IdentityUseCase = (*identityService)(nil)

type identityService struct {
	cfg         *config.Config
	userRepo    repo.UserRepo
	usersClient security.UsersClient
	revoker     security.SessionRevoker
}

func NewIdentityService(
	cfg *config.Config,
	userRepo repo.UserRepo,
	usersClient security.UsersClient,
	revoker security.SessionRevoker,
) usecase.IdentityUseCase {
	return &identityService{
		cfg:         cfg,
		userRepo:    userRepo,
		usersClient: usersClient,
		revoker:     revoker,
	}
}

func (s identityService) revokeSessions(ctx context.Context, idpUid string) error {
	if err := s.usersClient.LogoutUser(ctx, s.cfg.Keycloak.Realm, idpUid); err != nil {
		return fmt.Errorf("failed to terminate IdP sessions: %w", err)
	}
	if err := s.revoker.Revoke(ctx, idpUid); err != nil {
		return fmt.Errorf("failed to denylist revoked subject: %w", err)
	}
	return nil
}

func (s identityService) deleteOrphanedIdpUser(ctx context.Context, idpUid string) {
	if err := s.usersClient.DeleteUser(ctx, s.cfg.Keycloak.Realm, idpUid); err != nil {
		log.Printf("identity: failed to compensate orphaned Keycloak user %s: %v", idpUid, err)
	}
}

func (s identityService) deleteOrphanedIdpUserByEmail(ctx context.Context, email string) {
	users, err := s.usersClient.ListUsers(ctx, s.cfg.Keycloak.Realm, email)
	if err != nil {
		log.Printf("identity: failed to re-list orphaned Keycloak user for %s during compensation: %v", email, err)
		return
	}
	if len(users) == 0 {
		log.Printf("identity: no orphaned Keycloak user found for %s during compensation", email)
		return
	}
	for _, u := range users {
		if u.Id == nil {
			continue
		}
		s.deleteOrphanedIdpUser(ctx, *u.Id)
	}
}

// ─── User registration ─────────────────────────────────────────────────────

var ErrElevatedRoleOnRegister = errors.New("identity: registration may not grant elevated roles")
var ErrUserAlreadyExists = errors.New("identity: user already exists")
var selfServiceRolesAllowed = map[string]bool{
	string(openapi.AUTHROLE_USER): true,
}

func selfServiceRoles(requested []string) ([]string, error) {
	if len(requested) == 0 {
		return []string{
			string(openapi.AUTHROLE_USER),
		}, nil
	}
	for _, role := range requested {
		if !selfServiceRolesAllowed[strings.ToUpper(role)] {
			return nil, fmt.Errorf("%w: %q", ErrElevatedRoleOnRegister, role)
		}
	}
	return requested, nil
}

func (s identityService) RegisterUser(ctx context.Context, cmd command.UserRegisterCommand) (*aggregate.UserAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	users, err := s.usersClient.ListUsers(ctx, s.cfg.Keycloak.Realm, cmd.Email)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, fmt.Errorf("%w with email %s", ErrUserAlreadyExists, cmd.Email)
	}

	roles, err := selfServiceRoles(cmd.Roles)
	if err != nil {
		return nil, err
	}

	credentials := []idp.CredentialRepresentation{
		{
			Type:      ptr.String("password"),
			Value:     &cmd.Password,
			Temporary: ptr.Bool(false),
		},
	}
	request := idp.UserRepresentation{
		Username:    &cmd.Email,
		FirstName:   &cmd.FirstName,
		LastName:    &cmd.LastName,
		Email:       &cmd.Email,
		Enabled:     ptr.Bool(true),
		Credentials: &credentials,
	}
	err = s.usersClient.CreateUser(ctx, s.cfg.Keycloak.Realm, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in IdP: %w", err)
	}

	users, err = s.usersClient.ListUsers(ctx, s.cfg.Keycloak.Realm, cmd.Email)
	if err != nil {
		s.deleteOrphanedIdpUserByEmail(ctx, cmd.Email)
		return nil, fmt.Errorf("failed to look up user after creation in IdP: %w", err)
	}
	if len(users) == 0 {
		s.deleteOrphanedIdpUserByEmail(ctx, cmd.Email)
		return nil, fmt.Errorf("user not found after creation")
	}

	if users[0].Id == nil {
		s.deleteOrphanedIdpUserByEmail(ctx, cmd.Email)
		return nil, fmt.Errorf("created user has no IdP id")
	}
	if err := s.usersClient.SetUserRealmRoles(ctx, s.cfg.Keycloak.Realm, *users[0].Id, roles); err != nil {
		s.deleteOrphanedIdpUser(ctx, *users[0].Id)
		return nil, fmt.Errorf("failed to assign realm roles in IdP: %w", err)
	}

	localRow, err := s.upsertLocal(ctx, users[0], statusActive, roles)
	if err != nil {
		s.deleteOrphanedIdpUser(ctx, *users[0].Id)
		return nil, fmt.Errorf("failed to sync user to local DB: %w", err)
	}

	userAgg := &aggregate.UserAggregate{}
	userAgg.FromIdp(0, users[0])
	created := time.Now()
	updated := created
	if localRow != nil {
		if !localRow.CreatedAt.IsZero() {
			created = localRow.CreatedAt
		}
		if localRow.UpdatedAt != nil {
			updated = *localRow.UpdatedAt
		}
	}
	userAgg.ApplyRegistration(roles, created, updated)
	return userAgg, nil
}

// ─── User query ─────────────────────────────────────────────────────────────

func (s identityService) FindAllUsers(ctx context.Context, opt option.UserQueryOption) (int64, []*aggregate.UserAggregate, error) {
	total, users, err := s.userRepo.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}

	userAggs := make([]*aggregate.UserAggregate, len(users))
	for i, user := range users {
		userAggs[i] = aggregate.NewUserAggregate(user)
	}

	return total, userAggs, nil
}

func (s identityService) FindUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return aggregate.NewUserAggregate(user), nil
}

func (s identityService) FindUserByEmail(ctx context.Context, email string) (*aggregate.UserAggregate, error) {
	users, err := s.userRepo.FetchByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found with email %s", email)
	}
	return aggregate.NewUserAggregate(users[0]), nil
}

// ─── User management ────────────────────────────────────────────────────────

func (s identityService) UpdateUser(ctx context.Context, uid string, cmd command.UserUpdateCommand) (*aggregate.UserAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user.IdpId != nil {
		idpUser, err := s.usersClient.GetUser(ctx, s.cfg.Keycloak.Realm, *user.IdpId)
		if err != nil {
			return nil, fmt.Errorf("failed to get user from IdP: %w", err)
		}
		previous := idpUser
		if cmd.FirstName != nil {
			idpUser.FirstName = cmd.FirstName
		}
		if cmd.LastName != nil {
			idpUser.LastName = cmd.LastName
		}
		if err := s.usersClient.UpdateUser(ctx, s.cfg.Keycloak.Realm, *user.IdpId, idpUser); err != nil {
			return nil, fmt.Errorf("failed to update user in IdP: %w", err)
		}

		if err := s.userRepo.UpdateProfile(ctx, *user.IdpId, repo.UserProfileUpdate{
			FirstName: cmd.FirstName,
			LastName:  cmd.LastName,
			Phone:     cmd.Phone,
		}); err != nil {
			if compErr := s.usersClient.UpdateUser(ctx, s.cfg.Keycloak.Realm, *user.IdpId, previous); compErr != nil {
				logrus.WithContext(ctx).WithError(compErr).Errorf(
					"identity: DIVERGED profile for uid %s: local write failed (%v) and Keycloak compensation failed", uid, err)
			}
			return nil, fmt.Errorf("failed to update user in local DB: %w", err)
		}
	}

	updatedUser, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return aggregate.NewUserAggregate(updatedUser), nil
}

func (s identityService) UpdateUserRoles(ctx context.Context, uid string, cmd command.UserRolesCommand) (*aggregate.UserAggregate, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}

	if user.IdpId != nil {
		previousRoles := append([]string(nil), user.Roles...)
		if err := s.usersClient.SetUserRealmRoles(ctx, s.cfg.Keycloak.Realm, *user.IdpId, cmd.Roles); err != nil {
			return nil, fmt.Errorf("failed to update roles in IdP: %w", err)
		}
		if err := s.revokeSessions(ctx, *user.IdpId); err != nil {
			return nil, fmt.Errorf("failed to revoke sessions after role change: %w", err)
		}
		if err := s.userRepo.UpdateRoles(ctx, *user.IdpId, cmd.Roles); err != nil {
			if compErr := s.usersClient.SetUserRealmRoles(ctx, s.cfg.Keycloak.Realm, *user.IdpId, previousRoles); compErr != nil {
				logrus.WithContext(ctx).WithError(compErr).Errorf(
					"identity: DIVERGED roles for uid %s: local write failed (%v) and Keycloak compensation failed", uid, err)
			}
			return nil, fmt.Errorf("failed to update roles in local DB: %w", err)
		}
	}

	updatedUser, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return aggregate.NewUserAggregate(updatedUser), nil
}

func (s identityService) SuspendUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}

	if user.IdpId != nil {
		err = s.usersClient.UpdateUser(ctx, s.cfg.Keycloak.Realm, *user.IdpId, idp.UserRepresentation{Enabled: ptr.Bool(false)})
		if err != nil {
			return nil, fmt.Errorf("failed to suspend user in IdP: %w", err)
		}
		if err := s.revokeSessions(ctx, *user.IdpId); err != nil {
			return nil, fmt.Errorf("failed to revoke sessions after suspend: %w", err)
		}
		if err := s.userRepo.UpdateStatus(ctx, *user.IdpId, statusSuspended); err != nil {
			if retryErr := s.userRepo.UpdateStatus(ctx, *user.IdpId, statusSuspended); retryErr != nil {
				logrus.WithContext(ctx).WithError(retryErr).Errorf(
					"identity: DIVERGED suspend for uid %s: Keycloak disabled but local status write failed after retry; leaving Keycloak disabled", uid)
				return nil, fmt.Errorf("user %s suspended in IdP but local status update failed, please retry: %w", uid, retryErr)
			}
		}
	}

	updatedUser, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return aggregate.NewUserAggregate(updatedUser), nil
}

func (s identityService) ReactivateUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}

	if user.IdpId != nil {
		err = s.usersClient.UpdateUser(ctx, s.cfg.Keycloak.Realm, *user.IdpId, idp.UserRepresentation{Enabled: ptr.Bool(true)})
		if err != nil {
			return nil, fmt.Errorf("failed to reactivate user in IdP: %w", err)
		}
		if err := s.userRepo.UpdateStatus(ctx, *user.IdpId, statusActive); err != nil {
			return nil, fmt.Errorf("failed to update user status in local DB: %w", err)
		}
	}

	updatedUser, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return aggregate.NewUserAggregate(updatedUser), nil
}

func (s identityService) DeactivateUser(ctx context.Context, uid string) (*aggregate.UserAggregate, error) {
	user, err := s.userRepo.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user.IdpId != nil {
		err = s.usersClient.DeleteUser(ctx, s.cfg.Keycloak.Realm, *user.IdpId)
		if err != nil {
			return nil, fmt.Errorf("failed to deactivate user in IdP: %w", err)
		}
		if err := s.userRepo.SoftDelete(ctx, *user.IdpId); err != nil {
			return nil, fmt.Errorf("failed to soft-delete user in local DB: %w", err)
		}
	}
	user.Status = statusDeactivated
	return aggregate.NewUserAggregate(user), nil
}

// ─── Keycloak sync ──────────────────────────────────────────────────────────

func (s identityService) SyncUserFromIdp(ctx context.Context, idpId string) (*aggregate.UserAggregate, error) {
	idpUser, err := s.usersClient.GetUser(ctx, s.cfg.Keycloak.Realm, idpId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from IdP: %w", err)
	}

	roles, err := s.usersClient.GetUserRealmRoles(ctx, s.cfg.Keycloak.Realm, idpId)
	if err != nil {
		return nil, fmt.Errorf("failed to get realm roles from IdP: %w", err)
	}
	idpUser.RealmRoles = &roles

	status := statusActive
	if idpUser.Enabled != nil && !*idpUser.Enabled {
		status = statusSuspended
	}
	if _, err := s.upsertLocal(ctx, idpUser, status, roles); err != nil {
		return nil, fmt.Errorf("failed to upsert user to local DB: %w", err)
	}

	userAgg := &aggregate.UserAggregate{}
	userAgg.FromIdp(0, idpUser)
	return userAgg, nil
}

func (s identityService) upsertLocal(
	ctx context.Context,
	idpUser idp.UserRepresentation,
	status string,
	roles []string,
) (*repo.UserResult, error) {
	if idpUser.Id == nil {
		return nil, fmt.Errorf("IdP user is missing id")
	}
	profile := repo.UserUpsert{
		IdpUid: *idpUser.Id,
		Status: status,
	}
	if idpUser.Email != nil {
		profile.Email = *idpUser.Email
	}
	if idpUser.FirstName != nil {
		profile.FirstName = *idpUser.FirstName
	}
	if idpUser.LastName != nil {
		profile.LastName = *idpUser.LastName
	}
	switch {
	case roles != nil:
		profile.Roles = roles
	case idpUser.RealmRoles != nil:
		profile.Roles = *idpUser.RealmRoles
	}
	return s.userRepo.Upsert(ctx, profile)
}
