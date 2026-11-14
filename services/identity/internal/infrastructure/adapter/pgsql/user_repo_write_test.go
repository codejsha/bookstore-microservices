package pgsql

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/identity/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

var seededAt = time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)

func seedUserRow(t *testing.T, db *gorm.DB, mutate func(*entity.UsersEntity)) *entity.UsersEntity {
	t.Helper()
	updatedAt := seededAt
	u := &entity.UsersEntity{
		Uid:       uuid.NewString(),
		IdpUid:    uuid.NewString(),
		Email:     uuid.NewString() + "@x.com",
		FirstName: "First",
		LastName:  "Last",
		Role:      `{"values":["USER"]}`,
		Status:    "ACTIVE",
		CreatedAt: seededAt,
		UpdatedAt: &updatedAt,
		Actor:     1,
		Version:   1,
	}
	if mutate != nil {
		mutate(u)
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}

func stampedAfterSeed(row entity.UsersEntity) bool {
	return row.UpdatedAt != nil && row.UpdatedAt.After(seededAt)
}

func loadUserRow(t *testing.T, db *gorm.DB, idpUid string) entity.UsersEntity {
	t.Helper()
	var row entity.UsersEntity
	if err := db.Unscoped().Where("idp_uid = ?", idpUid).First(&row).Error; err != nil {
		t.Fatalf("load user %s: %v", idpUid, err)
	}
	return row
}

func TestUserRepositoryFindByUid_existingIdpUid_mapsStoredRow(t *testing.T) {
	r, db := newUserTestRepo(t)
	phone := "+1-555-0100"
	seedUserRow(t, db, nil)
	target := seedUserRow(t, db, func(u *entity.UsersEntity) {
		u.Email = "target@x.com"
		u.Phone = &phone
		u.Role = `{"values":["STAFF","MANAGER"]}`
		u.Status = "SUSPENDED"
	})

	got, err := r.FindByUid(context.Background(), target.IdpUid)
	if err != nil {
		t.Fatalf("find by uid: %v", err)
	}

	if got.IdpId == nil || *got.IdpId != target.IdpUid {
		t.Errorf("idp id = %v, want %s", got.IdpId, target.IdpUid)
	}
	if got.Email != "target@x.com" || got.Status != "SUSPENDED" {
		t.Errorf("result = %+v, want the target row", got)
	}
	if got.Phone == nil || *got.Phone != phone {
		t.Errorf("phone = %v, want %s", got.Phone, phone)
	}
	if want := []string{"STAFF", "MANAGER"}; !reflect.DeepEqual(got.Roles, want) {
		t.Errorf("roles = %v, want %v", got.Roles, want)
	}
}

func TestUserRepositoryFindByUid_unknownIdpUid_reportsRecordNotFound(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUserRow(t, db, nil)

	got, err := r.FindByUid(context.Background(), uuid.NewString())

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("err = %v, want gorm.ErrRecordNotFound", err)
	}
	if got != nil {
		t.Errorf("result = %+v, want nil", got)
	}
}

func TestUserRepositoryFindByUid_internalUidGiven_reportsRecordNotFound(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)

	_, err := r.FindByUid(context.Background(), target.Uid)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("err = %v, want lookups keyed on the idp uid only", err)
	}
}

func TestUserRepositoryFindByUid_softDeletedUser_reportsRecordNotFound(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)
	if err := r.SoftDelete(context.Background(), target.IdpUid); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	_, err := r.FindByUid(context.Background(), target.IdpUid)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("err = %v, want gorm.ErrRecordNotFound", err)
	}
}

func TestUserRepositoryFetchByEmail_emailShared_keepsOnlyLiveUsersWithThatEmail(t *testing.T) {
	r, db := newUserTestRepo(t)
	first := seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "shared@x.com" })
	second := seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "shared@x.com" })
	deleted := seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "shared@x.com" })
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "other@x.com" })
	if err := r.SoftDelete(context.Background(), deleted.IdpUid); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	got, err := r.FetchByEmail(context.Background(), "shared@x.com")
	if err != nil {
		t.Fatalf("fetch by email: %v", err)
	}

	found := map[string]bool{}
	for _, u := range got {
		if u.IdpId == nil {
			t.Fatalf("result = %+v, want an idp id", u)
		}
		found[*u.IdpId] = true
	}
	if want := map[string]bool{first.IdpUid: true, second.IdpUid: true}; !reflect.DeepEqual(found, want) {
		t.Fatalf("idp ids = %v, want %v", found, want)
	}
}

func TestUserRepositoryFetchByEmail_noMatch_emptyList(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUserRow(t, db, nil)

	got, err := r.FetchByEmail(context.Background(), "nobody@x.com")
	if err != nil {
		t.Fatalf("fetch by email: %v", err)
	}

	if got == nil || len(got) != 0 {
		t.Fatalf("results = %#v, want an empty non-nil list", got)
	}
}

func TestUserRepositoryFindAll_emailGiven_keepsOnlyExactEmail(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUser(t, db, "a@x.com", "ACTIVE")
	seedUser(t, db, "aa@x.com", "ACTIVE")
	seedUser(t, db, "b@x.com", "ACTIVE")
	email := "a@x.com"

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithEmail(&email),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if want := []string{"a@x.com"}; !reflect.DeepEqual(userEmails(got), want) {
		t.Fatalf("emails = %v, want %v", userEmails(got), want)
	}
}

func TestUserRepositoryFindAll_phoneGiven_keepsOnlyExactPhone(t *testing.T) {
	r, db := newUserTestRepo(t)
	phone, otherPhone := "+1-555-0100", "+1-555-0199"
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "a@x.com"; u.Phone = &phone })
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "b@x.com"; u.Phone = &otherPhone })
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "c@x.com" })

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithPhone(&phone),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if want := []string{"a@x.com"}; !reflect.DeepEqual(userEmails(got), want) {
		t.Fatalf("emails = %v, want %v", userEmails(got), want)
	}
}

func TestUserRepositoryFindAll_phoneAndStatusGiven_requiresBoth(t *testing.T) {
	r, db := newUserTestRepo(t)
	phone := "+1-555-0100"
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "a@x.com"; u.Phone = &phone })
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "b@x.com"; u.Phone = &phone; u.Status = "SUSPENDED" })
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "c@x.com"; u.Status = "SUSPENDED" })
	status := "SUSPENDED"

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithPhone(&phone),
		option.UserQueryOption{}.WithStatus(&status),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if want := []string{"b@x.com"}; !reflect.DeepEqual(userEmails(got), want) {
		t.Fatalf("emails = %v, want %v", userEmails(got), want)
	}
}

func TestUserRepositoryFindAll_filterMatchesMoreThanOnePage_totalCountsEveryMatch(t *testing.T) {
	r, db := newUserTestRepo(t)
	for i := 0; i < 5; i++ {
		seedUserRow(t, db, func(u *entity.UsersEntity) { u.Status = "SUSPENDED" })
	}
	seedUserRow(t, db, nil)
	status, size, page := "SUSPENDED", int32(2), int32(2)

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithStatus(&status),
		option.UserQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, nil)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if len(got) != 1 {
		t.Fatalf("page = %v, want the single row left on the last page", userEmails(got))
	}
}

func TestUserRepositoryUpsert_unknownIdpUid_insertsRowWithGeneratedUid(t *testing.T) {
	r, db := newUserTestRepo(t)
	phone := "+1-555-0100"
	before := time.Now().UTC().Add(-time.Second)

	got, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: "idp-1", Email: "new@x.com", FirstName: "Ada", LastName: "Lovelace",
		Phone: &phone, Roles: []string{"USER", "STAFF"}, Status: "ACTIVE",
	})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}

	if got.Id == 0 {
		t.Errorf("id = 0, want the stored row id")
	}
	if got.IdpId == nil || *got.IdpId != "idp-1" {
		t.Errorf("idp id = %v, want idp-1", got.IdpId)
	}
	if want := []string{"USER", "STAFF"}; !reflect.DeepEqual(got.Roles, want) {
		t.Errorf("roles = %v, want %v", got.Roles, want)
	}
	row := loadUserRow(t, db, "idp-1")
	parsed, err := uuid.Parse(row.Uid)
	if err != nil {
		t.Fatalf("stored uid %q: %v", row.Uid, err)
	}
	if parsed.Version() != 7 {
		t.Errorf("uid version = %d, want a time-ordered v7 uuid", parsed.Version())
	}
	if row.Email != "new@x.com" || row.FirstName != "Ada" || row.LastName != "Lovelace" || row.Status != "ACTIVE" {
		t.Errorf("row = %+v, want the upserted profile", row)
	}
	if row.Role != `{"values":["USER","STAFF"]}` {
		t.Errorf("stored roles = %s, want the values envelope", row.Role)
	}
	if row.Version != 1 {
		t.Errorf("version = %d, want 1", row.Version)
	}
	if row.CreatedAt.Before(before) {
		t.Errorf("created_at = %v, want it set at insert time", row.CreatedAt)
	}
}

func TestUserRepositoryUpsert_rolesOmitted_insertsEmptyRoleEnvelope(t *testing.T) {
	r, db := newUserTestRepo(t)

	if _, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: "idp-1", Email: "new@x.com", FirstName: "F", LastName: "L", Status: "ACTIVE",
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	if row := loadUserRow(t, db, "idp-1"); row.Role != emptyRolesJSON {
		t.Errorf("stored roles = %s, want %s", row.Role, emptyRolesJSON)
	}
}

func TestUserRepositoryUpsert_existingUserRolesAndStatusOmitted_keepsStoredRolesAndStatus(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, func(u *entity.UsersEntity) {
		u.Role = `{"values":["MANAGER"]}`
		u.Status = "SUSPENDED"
	})

	got, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: target.IdpUid, Email: "renamed@x.com", FirstName: "New", LastName: "Name",
	})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}

	if want := []string{"MANAGER"}; !reflect.DeepEqual(got.Roles, want) {
		t.Errorf("roles = %v, want stored %v", got.Roles, want)
	}
	if got.Status != "SUSPENDED" {
		t.Errorf("status = %q, want stored SUSPENDED", got.Status)
	}
	row := loadUserRow(t, db, target.IdpUid)
	if row.Email != "renamed@x.com" || row.FirstName != "New" || row.LastName != "Name" {
		t.Errorf("row = %+v, want the profile overwritten", row)
	}
	if row.Uid != target.Uid || row.Id != target.Id {
		t.Errorf("row identity = (%d, %s), want it unchanged (%d, %s)", row.Id, row.Uid, target.Id, target.Uid)
	}
	if !stampedAfterSeed(row) {
		t.Errorf("updated_at = %v, want it stamped on update", row.UpdatedAt)
	}
}

func TestUserRepositoryUpsert_existingUserRolesAndStatusGiven_overwritesBoth(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)

	if _, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: target.IdpUid, Email: target.Email, FirstName: "F", LastName: "L",
		Roles: []string{"STAFF"}, Status: "SUSPENDED",
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	row := loadUserRow(t, db, target.IdpUid)
	if row.Role != `{"values":["STAFF"]}` || row.Status != "SUSPENDED" {
		t.Errorf("row = (%s, %s), want the given roles and status", row.Role, row.Status)
	}
	var count int64
	if err := db.Model(&entity.UsersEntity{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("rows = %d, want the existing row updated in place", count)
	}
}

func TestUserRepositoryUpdateProfile_partialPatch_changesOnlyGivenFieldsOfThatUser(t *testing.T) {
	r, db := newUserTestRepo(t)
	phone := "+1-555-0100"
	target := seedUserRow(t, db, func(u *entity.UsersEntity) { u.Phone = &phone })
	bystander := seedUserRow(t, db, nil)
	firstName := "Grace"

	err := r.UpdateProfile(context.Background(), target.IdpUid, repo.UserProfileUpdate{FirstName: &firstName})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}

	row := loadUserRow(t, db, target.IdpUid)
	if row.FirstName != "Grace" {
		t.Errorf("first name = %q, want Grace", row.FirstName)
	}
	if row.LastName != "Last" || row.Phone == nil || *row.Phone != phone {
		t.Errorf("row = %+v, want last name and phone untouched", row)
	}
	if !stampedAfterSeed(row) {
		t.Errorf("updated_at = %v, want it stamped", row.UpdatedAt)
	}
	other := loadUserRow(t, db, bystander.IdpUid)
	if other.FirstName != "First" || stampedAfterSeed(other) {
		t.Errorf("bystander = %+v, want it untouched", other)
	}
}

func TestUserRepositoryUpdateProfile_fullPatch_overwritesNamesAndPhone(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)
	firstName, lastName, phone := "Grace", "Hopper", "+1-555-0199"

	err := r.UpdateProfile(context.Background(), target.IdpUid, repo.UserProfileUpdate{
		FirstName: &firstName, LastName: &lastName, Phone: &phone,
	})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}

	row := loadUserRow(t, db, target.IdpUid)
	if row.FirstName != "Grace" || row.LastName != "Hopper" || row.Phone == nil || *row.Phone != phone {
		t.Errorf("row = %+v, want the full patch applied", row)
	}
	if row.Email != target.Email || row.Status != target.Status {
		t.Errorf("row = %+v, want email and status untouched", row)
	}
}

func TestUserRepositoryUpdateProfile_emptyPatch_stampsOnlyUpdatedAt(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)

	if err := r.UpdateProfile(context.Background(), target.IdpUid, repo.UserProfileUpdate{}); err != nil {
		t.Fatalf("update profile: %v", err)
	}

	row := loadUserRow(t, db, target.IdpUid)
	if row.FirstName != "First" || row.LastName != "Last" || row.Phone != nil {
		t.Errorf("row = %+v, want the profile untouched", row)
	}
	if !stampedAfterSeed(row) {
		t.Errorf("updated_at = %v, want it stamped", row.UpdatedAt)
	}
}

func TestUserRepositoryUpdateStatus_existingUser_changesOnlyThatUsersStatus(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)
	bystander := seedUserRow(t, db, nil)

	if err := r.UpdateStatus(context.Background(), target.IdpUid, "SUSPENDED"); err != nil {
		t.Fatalf("update status: %v", err)
	}

	row := loadUserRow(t, db, target.IdpUid)
	if row.Status != "SUSPENDED" || !stampedAfterSeed(row) {
		t.Errorf("row = %+v, want SUSPENDED with updated_at stamped", row)
	}
	if other := loadUserRow(t, db, bystander.IdpUid); other.Status != "ACTIVE" {
		t.Errorf("bystander status = %q, want ACTIVE", other.Status)
	}
}

func TestUserRepositoryUpdateStatus_unknownIdpUid_changesNothingWithoutError(t *testing.T) {
	r, db := newUserTestRepo(t)
	bystander := seedUserRow(t, db, nil)

	if err := r.UpdateStatus(context.Background(), uuid.NewString(), "SUSPENDED"); err != nil {
		t.Fatalf("update status: %v", err)
	}

	if other := loadUserRow(t, db, bystander.IdpUid); other.Status != "ACTIVE" || stampedAfterSeed(other) {
		t.Errorf("bystander = %+v, want it untouched", other)
	}
}

func TestUserRepositoryUpdateRoles_rolesGiven_storesRoleEnvelopeForThatUser(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)
	bystander := seedUserRow(t, db, nil)

	if err := r.UpdateRoles(context.Background(), target.IdpUid, []string{"STAFF", "MANAGER"}); err != nil {
		t.Fatalf("update roles: %v", err)
	}

	row := loadUserRow(t, db, target.IdpUid)
	if row.Role != `{"values":["STAFF","MANAGER"]}` || !stampedAfterSeed(row) {
		t.Errorf("row = %+v, want the new roles with updated_at stamped", row)
	}
	if other := loadUserRow(t, db, bystander.IdpUid); other.Role != `{"values":["USER"]}` {
		t.Errorf("bystander roles = %s, want them untouched", other.Role)
	}
	got, err := r.FindByUid(context.Background(), target.IdpUid)
	if err != nil {
		t.Fatalf("find by uid: %v", err)
	}
	if want := []string{"STAFF", "MANAGER"}; !reflect.DeepEqual(got.Roles, want) {
		t.Errorf("roles = %v, want %v", got.Roles, want)
	}
}

func TestUserRepositoryUpdateRoles_noRoles_storesEmptyRoleEnvelope(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)

	if err := r.UpdateRoles(context.Background(), target.IdpUid, nil); err != nil {
		t.Fatalf("update roles: %v", err)
	}

	if row := loadUserRow(t, db, target.IdpUid); row.Role != emptyRolesJSON {
		t.Errorf("stored roles = %s, want %s", row.Role, emptyRolesJSON)
	}
}

func TestUserRepositorySoftDelete_existingUser_marksDeactivatedAndKeepsTheRow(t *testing.T) {
	r, db := newUserTestRepo(t)
	target := seedUserRow(t, db, nil)
	bystander := seedUserRow(t, db, nil)

	if err := r.SoftDelete(context.Background(), target.IdpUid); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	row := loadUserRow(t, db, target.IdpUid)
	if row.Status != "DEACTIVATED" {
		t.Errorf("status = %q, want DEACTIVATED", row.Status)
	}
	if !row.DeletedAt.Valid {
		t.Errorf("deleted_at = %+v, want it stamped", row.DeletedAt)
	}
	if !stampedAfterSeed(row) {
		t.Errorf("updated_at = %v, want it stamped", row.UpdatedAt)
	}
	other := loadUserRow(t, db, bystander.IdpUid)
	if other.Status != "ACTIVE" || other.DeletedAt.Valid {
		t.Errorf("bystander = %+v, want it live and ACTIVE", other)
	}
	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption())
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if total != 1 || len(got) != 1 || got[0].Email != bystander.Email {
		t.Errorf("live users = %d %v, want only the bystander", total, userEmails(got))
	}
}

func TestUserRepositorySoftDelete_unknownIdpUid_changesNothingWithoutError(t *testing.T) {
	r, db := newUserTestRepo(t)
	bystander := seedUserRow(t, db, nil)

	if err := r.SoftDelete(context.Background(), uuid.NewString()); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	if other := loadUserRow(t, db, bystander.IdpUid); other.Status != "ACTIVE" || other.DeletedAt.Valid {
		t.Errorf("bystander = %+v, want it untouched", other)
	}
}

func TestUserRepository_connectionClosed_reportsDatabaseError(t *testing.T) {
	ctx := context.Background()
	calls := map[string]func(r userRepository) error{
		"findAll": func(r userRepository) error {
			_, _, err := r.FindAll(ctx, option.NewUserQueryOption())
			return err
		},
		"findByUid": func(r userRepository) error {
			_, err := r.FindByUid(ctx, "idp-1")
			return err
		},
		"fetchByEmail": func(r userRepository) error {
			_, err := r.FetchByEmail(ctx, "a@x.com")
			return err
		},
		"upsert": func(r userRepository) error {
			_, err := r.Upsert(ctx, repo.UserUpsert{IdpUid: "idp-1", Email: "a@x.com", Status: "ACTIVE"})
			return err
		},
		"updateProfile": func(r userRepository) error {
			return r.UpdateProfile(ctx, "idp-1", repo.UserProfileUpdate{})
		},
		"updateStatus": func(r userRepository) error {
			return r.UpdateStatus(ctx, "idp-1", "SUSPENDED")
		},
		"updateRoles": func(r userRepository) error {
			return r.UpdateRoles(ctx, "idp-1", []string{"STAFF"})
		},
		"softDelete": func(r userRepository) error {
			return r.SoftDelete(ctx, "idp-1")
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			r, db := newUserTestRepo(t)
			sqlDB, err := db.DB()
			if err != nil {
				t.Fatalf("sql db: %v", err)
			}
			if err := sqlDB.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}

			err = call(r)

			if err == nil {
				t.Fatal("err = nil, want the database failure surfaced")
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				t.Fatalf("err = %v, want a failure distinct from a missing row", err)
			}
		})
	}
}

func TestUserRepositoryFindAll_nameFilterGiven_matchesFirstOrLastNameIgnoringCase(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "a@x.com"; u.FirstName = "Bilbo"; u.LastName = "Baggins" })
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "b@x.com"; u.FirstName = "Samwise"; u.LastName = "Gamgee" })
	seedUserRow(t, db, func(u *entity.UsersEntity) { u.Email = "c@x.com"; u.FirstName = "Gimli"; u.LastName = "Gloinsson" })
	name := "GAM"

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithName(&name),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 1 || !reflect.DeepEqual(userEmails(got), []string{"b@x.com"}) {
		t.Errorf("total = %d, emails = %v, want 1 and [b@x.com]", total, userEmails(got))
	}
}
