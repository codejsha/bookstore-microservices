package protostub

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

const (
	actorUidMetadataKey     = "x-actor-uid"
	actorAdminMetadataKey   = "x-actor-admin"
	actorAdminMetadataValue = "true"
)

func withActor(ctx context.Context) context.Context {
	caller := httpx.CallerFromContext(ctx)
	if caller.UserUid == "" {
		return ctx
	}
	ctx = metadata.AppendToOutgoingContext(ctx, actorUidMetadataKey, caller.UserUid)
	if caller.IsAdmin {
		ctx = metadata.AppendToOutgoingContext(ctx, actorAdminMetadataKey, actorAdminMetadataValue)
	}
	return ctx
}
