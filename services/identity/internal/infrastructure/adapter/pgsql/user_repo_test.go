package pgsql

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/identity/generated/infrastructure/port/entity"
)

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
		{"whenRolesPresent_returnsRoles", `{"values":["STAFF","MANAGE"]}`, []string{"STAFF", "MANAGE"}},
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
