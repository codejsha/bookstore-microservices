package pgsql

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/identity/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"
)

const usersSchemaDDL = `CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid TEXT, idp_uid TEXT, email TEXT, first_name TEXT, last_name TEXT, phone TEXT,
	roles TEXT, status TEXT, last_login_at DATETIME,
	created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
	actor INTEGER, version INTEGER
)`

func newUserTestRepo(t *testing.T) (userRepository, *gorm.DB) {
	t.Helper()
	dsn := "file:" + uuid.NewString() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.Exec(usersSchemaDDL).Error; err != nil {
		t.Fatalf("create users schema: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	return userRepository{db: db}, db
}

func seedUser(t *testing.T, db *gorm.DB, email, status string) {
	t.Helper()
	u := &entity.UsersEntity{
		Uid:       uuid.NewString(),
		IdpUid:    uuid.NewString(),
		Email:     email,
		FirstName: "F",
		LastName:  "L",
		Role:      emptyRolesJSON,
		Status:    status,
		CreatedAt: time.Now(),
		Actor:     1,
		Version:   1,
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
}

func userEmails(users []*repo.UserResult) []string {
	out := make([]string, len(users))
	for i, u := range users {
		out[i] = u.Email
	}
	return out
}

func TestUserRepositoryFindAll_secondPage_skipsOneFullPage(t *testing.T) {
	r, db := newUserTestRepo(t)
	for i := 0; i < 7; i++ {
		seedUser(t, db, fmt.Sprintf("u%d@x.com", i), "ACTIVE")
	}
	size, page, sort := int32(3), int32(1), "email,asc"

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, &sort)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 7 {
		t.Errorf("total = %d, want 7", total)
	}
	if want := []string{"u3@x.com", "u4@x.com", "u5@x.com"}; !reflect.DeepEqual(userEmails(got), want) {
		t.Fatalf("emails = %v, want %v", userEmails(got), want)
	}
}

func TestUserRepositoryFindAll_sortOmitted_newestFirst(t *testing.T) {
	r, db := newUserTestRepo(t)
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 3; i++ {
		u := &entity.UsersEntity{
			Uid:       uuid.NewString(),
			IdpUid:    uuid.NewString(),
			Email:     fmt.Sprintf("u%d@x.com", i),
			FirstName: "F",
			LastName:  "L",
			Role:      emptyRolesJSON,
			Status:    "ACTIVE",
			CreatedAt: base.Add(time.Duration(i) * time.Hour),
			Actor:     1,
			Version:   1,
		}
		if err := db.Create(u).Error; err != nil {
			t.Fatalf("seed user: %v", err)
		}
	}
	size, page := int32(3), int32(0)

	_, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithPage(pagination.NewPageOption(&size, &page, nil)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if want := []string{"u3@x.com", "u2@x.com", "u1@x.com"}; !reflect.DeepEqual(userEmails(got), want) {
		t.Fatalf("emails = %v, want %v", userEmails(got), want)
	}
}

func TestUserRepositoryFindAll_statusGiven_keepsOnlyMatchingUsers(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUser(t, db, "a@x.com", "ACTIVE")
	seedUser(t, db, "b@x.com", "SUSPENDED")
	seedUser(t, db, "c@x.com", "ACTIVE")
	seedUser(t, db, "d@x.com", "SUSPENDED")
	status, sort := "SUSPENDED", "email,asc"

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithStatus(&status),
		option.UserQueryOption{}.WithPage(pagination.NewPageOption(nil, nil, &sort)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if want := []string{"b@x.com", "d@x.com"}; !reflect.DeepEqual(userEmails(got), want) {
		t.Fatalf("emails = %v, want %v", userEmails(got), want)
	}
}

func TestUserRepositoryFindAll_statusEmpty_appliesNoStatusFilter(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUser(t, db, "a@x.com", "ACTIVE")
	seedUser(t, db, "b@x.com", "SUSPENDED")
	status := ""

	total, _, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithStatus(&status),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
}

func seedNamedUser(t *testing.T, db *gorm.DB, email, firstName string, createdAt time.Time) {
	t.Helper()
	u := &entity.UsersEntity{
		Uid:       uuid.NewString(),
		IdpUid:    uuid.NewString(),
		Email:     email,
		FirstName: firstName,
		LastName:  "L",
		Role:      emptyRolesJSON,
		Status:    "ACTIVE",
		CreatedAt: createdAt,
		Actor:     1,
		Version:   1,
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
}

func seedSortableUsers(t *testing.T, db *gorm.DB) {
	t.Helper()
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	seedNamedUser(t, db, "a@x.com", "Carol", base.Add(1*time.Hour))
	seedNamedUser(t, db, "b@x.com", "Alice", base.Add(3*time.Hour))
	seedNamedUser(t, db, "c@x.com", "Bob", base.Add(2*time.Hour))
}

func findSorted(t *testing.T, r userRepository, sort string) []string {
	t.Helper()
	_, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
		option.UserQueryOption{}.WithPage(pagination.NewPageOption(nil, nil, &sort)),
	))
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	return userEmails(got)
}

func TestUserRepositoryFindAll_snakeCaseSort_ordersPage(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedSortableUsers(t, db)

	got := findSorted(t, r, "first_name,asc")

	if want := []string{"b@x.com", "c@x.com", "a@x.com"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("emails = %v, want %v", got, want)
	}
}

func TestUserRepositoryFindAll_camelCaseSort_ordersPage(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedSortableUsers(t, db)

	got := findSorted(t, r, "firstName,desc")

	if want := []string{"a@x.com", "c@x.com", "b@x.com"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("emails = %v, want %v", got, want)
	}
}

func TestUserRepositoryFindAll_unknownSortOnly_fallsBackToNewestFirst(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedSortableUsers(t, db)

	got := findSorted(t, r, "bogus,asc;nope")

	if want := []string{"b@x.com", "c@x.com", "a@x.com"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("emails = %v, want %v", got, want)
	}
}

func TestUserRepositoryFindAll_mixedSort_keepsOnlyKnownKeys(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedSortableUsers(t, db)

	got := findSorted(t, r, "bogus,desc;created_at,asc")

	if want := []string{"a@x.com", "c@x.com", "b@x.com"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("emails = %v, want %v", got, want)
	}
}

func TestUserRepositoryFindAll_deactivatedStatus_includesSoftDeletedUsers(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUser(t, db, "a@x.com", "ACTIVE")
	seedUser(t, db, "b@x.com", "ACTIVE")
	var target entity.UsersEntity
	if err := db.Where("email = ?", "b@x.com").First(&target).Error; err != nil {
		t.Fatalf("load seeded user: %v", err)
	}
	if err := r.SoftDelete(context.Background(), target.IdpUid); err != nil {
		t.Fatalf("soft delete: %v", err)
	}
	status := "DEACTIVATED"

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption(
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

func TestUserRepositoryFindAll_noStatus_excludesSoftDeletedUsers(t *testing.T) {
	r, db := newUserTestRepo(t)
	seedUser(t, db, "a@x.com", "ACTIVE")
	seedUser(t, db, "b@x.com", "ACTIVE")
	var target entity.UsersEntity
	if err := db.Where("email = ?", "b@x.com").First(&target).Error; err != nil {
		t.Fatalf("load seeded user: %v", err)
	}
	if err := r.SoftDelete(context.Background(), target.IdpUid); err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	total, got, err := r.FindAll(context.Background(), option.NewUserQueryOption())
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

func TestUserRepositoryUpsert_phoneOmitted_keepsStoredPhone(t *testing.T) {
	r, _ := newUserTestRepo(t)
	phone := "+1-555-0100"
	if _, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: "idp-1", Email: "u@x.com", FirstName: "F", LastName: "L", Phone: &phone, Status: "ACTIVE",
	}); err != nil {
		t.Fatalf("first upsert: %v", err)
	}

	got, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: "idp-1", Email: "u@x.com", FirstName: "F2", LastName: "L", Status: "ACTIVE",
	})
	if err != nil {
		t.Fatalf("second upsert: %v", err)
	}

	if got.Phone == nil || *got.Phone != phone {
		t.Errorf("phone = %v, want stored %q", got.Phone, phone)
	}
	if got.FirstName != "F2" {
		t.Errorf("first name = %q, want F2", got.FirstName)
	}
}

func TestUserRepositoryUpsert_phoneGiven_overwritesStoredPhone(t *testing.T) {
	r, _ := newUserTestRepo(t)
	oldPhone, newPhone := "+1-555-0100", "+1-555-0199"
	if _, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: "idp-1", Email: "u@x.com", FirstName: "F", LastName: "L", Phone: &oldPhone, Status: "ACTIVE",
	}); err != nil {
		t.Fatalf("first upsert: %v", err)
	}

	got, err := r.Upsert(context.Background(), repo.UserUpsert{
		IdpUid: "idp-1", Email: "u@x.com", FirstName: "F", LastName: "L", Phone: &newPhone, Status: "ACTIVE",
	})
	if err != nil {
		t.Fatalf("second upsert: %v", err)
	}

	if got.Phone == nil || *got.Phone != newPhone {
		t.Errorf("phone = %v, want %q", got.Phone, newPhone)
	}
}

func TestResolveSort_snakeAndCamelKeys_keepsBoth(t *testing.T) {
	sortable := sortableWithAliases(map[string]field.OrderExpr{
		"created_at": field.NewTime("users", "created_at"),
		"uid":        field.NewString("users", "idp_uid"),
	})
	sort := "created_at,asc;createdAt,desc;bogus,asc;uid"

	got := resolveSort(pagination.NewPageOption(nil, nil, &sort), sortable).GetSort()

	if want := "created_at,asc;createdAt,desc;uid"; got != want {
		t.Errorf("sort = %q, want %q", got, want)
	}
}

func TestResolveSort_noResolvableKey_usesDefaultSort(t *testing.T) {
	sortable := sortableWithAliases(map[string]field.OrderExpr{
		"created_at": field.NewTime("users", "created_at"),
	})
	sort := "bogus,asc"

	got := resolveSort(pagination.NewPageOption(nil, nil, &sort), sortable).GetSort()

	if got != defaultSort {
		t.Errorf("sort = %q, want %q", got, defaultSort)
	}
}

func TestToUserResult_WhenRowFullyPopulated_MapsEveryField(t *testing.T) {
	now := time.Now()
	updated := now.Add(time.Hour)
	phone := "+1-555"
	got := toUserResult(&entity.UsersEntity{
		Id: 11, IdpUid: "idp-1", Email: "u@x.com", FirstName: "F", LastName: "L",
		Phone: &phone, Role: `{"values":["USER","STAFF"]}`, Status: "ACTIVE",
		LastLoginAt: &updated, CreatedAt: now, UpdatedAt: &updated,
	})
	if got.Id != 11 || got.IdpId == nil || *got.IdpId != "idp-1" {
		t.Errorf("got = %+v", got)
	}
	if got.Email != "u@x.com" || got.Status != "ACTIVE" {
		t.Errorf("scalar fields = %+v", got)
	}
	if got.Phone == nil || *got.Phone != phone {
		t.Errorf("Phone = %v", got.Phone)
	}
	if !reflect.DeepEqual(got.Roles, []string{"USER", "STAFF"}) {
		t.Errorf("Roles = %v", got.Roles)
	}
}

func TestToUserResult_WhenIdpUidEmpty_ReturnsNilIdpId(t *testing.T) {
	got := toUserResult(&entity.UsersEntity{IdpUid: ""})
	if got.IdpId != nil {
		t.Errorf("IdpId = %v, want nil for empty IdpUid", got.IdpId)
	}
}

func TestParseRoles(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"whenInputEmpty_returnsNil", "", nil},
		{"whenJsonMalformed_returnsNil", "not json", nil},
		{"whenRolesPresent_returnsRoles", `{"values":["STAFF","MANAGER"]}`, []string{"STAFF", "MANAGER"}},
		{"whenValuesEmpty_returnsEmptySlice", `{"values":[]}`, []string{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := parseRoles(c.in)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got = %v, want %v", got, c.want)
			}
		})
	}
}

func TestEncodeRoles(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want string
	}{
		{"whenRolesNil_returnsEmptyJsonArray", nil, `{"values":[]}`},
		{"whenRolesEmpty_returnsEmptyJsonArray", []string{}, `{"values":[]}`},
		{"whenRolesPresent_returnsJsonValues", []string{"USER", "STAFF"}, `{"values":["USER","STAFF"]}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := encodeRoles(c.in)
			if got != c.want {
				t.Errorf("got = %q, want %q", got, c.want)
			}
			if !json.Valid([]byte(got)) {
				t.Errorf("got = %q, want valid JSON for a JSONB NOT NULL column", got)
			}
		})
	}
}
