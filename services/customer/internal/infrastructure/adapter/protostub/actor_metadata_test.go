package protostub

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

func TestWithActor(t *testing.T) {
	t.Run("attaches x-actor-uid for a non-admin caller and omits x-actor-admin", func(t *testing.T) {
		ctx := httpx.WithCaller(context.Background(), httpx.Caller{UserUid: "u-1", IsAdmin: false})
		md, ok := metadata.FromOutgoingContext(withActor(ctx))
		if !ok {
			t.Fatal("no outgoing metadata set")
		}
		uid := md.Get(actorUidMetadataKey)
		if len(uid) != 1 || uid[0] != "u-1" {
			t.Errorf("%s = %v, want [u-1]", actorUidMetadataKey, uid)
		}
		if admin := md.Get(actorAdminMetadataKey); len(admin) != 0 {
			t.Errorf("%s = %v, want none for a non-admin caller", actorAdminMetadataKey, admin)
		}
	})

	t.Run("attaches x-actor-admin=true for an admin caller", func(t *testing.T) {
		ctx := httpx.WithCaller(context.Background(), httpx.Caller{UserUid: "admin-1", IsAdmin: true})
		md, ok := metadata.FromOutgoingContext(withActor(ctx))
		if !ok {
			t.Fatal("no outgoing metadata set")
		}
		uid := md.Get(actorUidMetadataKey)
		if len(uid) != 1 || uid[0] != "admin-1" {
			t.Errorf("%s = %v, want [admin-1]", actorUidMetadataKey, uid)
		}
		admin := md.Get(actorAdminMetadataKey)
		if len(admin) != 1 || admin[0] != actorAdminMetadataValue {
			t.Errorf("%s = %v, want [%s]", actorAdminMetadataKey, admin, actorAdminMetadataValue)
		}
	})

	t.Run("no metadata when the request carried no principal", func(t *testing.T) {
		if _, ok := metadata.FromOutgoingContext(withActor(context.Background())); ok {
			t.Error("outgoing metadata was set for an unauthenticated context, want none")
		}
	})
}
