package protostub

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

func TestWithUserAuthorization_WhenBearerTokenPresent_AttachesAuthorizationMetadata(t *testing.T) {
	t.Run("whenBearerTokenPresent_attachesAuthorizationMetadata", func(t *testing.T) {
		ctx := httpx.WithBearerToken(context.Background(), "Bearer abc.def.ghi")
		md, ok := metadata.FromOutgoingContext(withUserAuthorization(ctx))
		if !ok {
			t.Fatal("no outgoing metadata set")
		}
		got := md.Get(userAuthorizationMetadataKey)
		if len(got) != 1 || got[0] != "Bearer abc.def.ghi" {
			t.Errorf("authorization metadata = %v, want [Bearer abc.def.ghi]", got)
		}
	})

	t.Run("whenNoToken_attachesNoMetadata", func(t *testing.T) {
		if _, ok := metadata.FromOutgoingContext(withUserAuthorization(context.Background())); ok {
			t.Error("outgoing metadata was set for an unauthenticated context, want none")
		}
	})
}

func TestToCustomerAggregate_WhenProtoHasPhoneAndRoles_MapsEveryField(t *testing.T) {
	t.Run("whenPhoneAndRolesPresent_mapsEveryField", func(t *testing.T) {
		got := toCustomerAggregate(&userpb.User{
			Uid:       "u-1",
			Email:     "a@b.com",
			FirstName: "A",
			LastName:  "B",
			Phone:     "+1-555",
			Roles:     []string{"ORDER", "VIEW", "BOGUS"},
		})
		if got.Uid != "u-1" || got.Email != "a@b.com" {
			t.Errorf("identity fields = %+v", got)
		}
		if got.Phone == nil || *got.Phone != "+1-555" {
			t.Errorf("phone = %v, want +1-555", got.Phone)
		}
		want := []constant.AuthRole{constant.AUTHROLE_ORDER, constant.AUTHROLE_VIEW, constant.AUTHROLE_UNKNOWN}
		if !reflect.DeepEqual(got.Roles, want) {
			t.Errorf("roles = %+v, want %+v", got.Roles, want)
		}
	})

	t.Run("whenPhoneMissing_returnsNilPhone", func(t *testing.T) {
		got := toCustomerAggregate(&userpb.User{Uid: "u-2", Email: "x@y.com"})
		if got.Phone != nil {
			t.Errorf("phone = %v, want nil for empty string", got.Phone)
		}
		if len(got.Roles) != 0 {
			t.Errorf("roles = %v, want empty", got.Roles)
		}
	})
}
