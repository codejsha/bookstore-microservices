package restcontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/codejsha/bookstore-microservices/inventory/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/httpx"
)

func requireUidParam(ctx context.Context, field, value string) error {
	if _, err := uuid.Parse(value); err != nil {
		return httpx.MapBusinessError(ctx, fmt.Errorf("%s must be a valid uuid: %w", field, command.ErrInvalidCommand))
	}
	return nil
}

func optionalUidParam(ctx context.Context, field string, value *string) error {
	if value == nil || *value == "" {
		return nil
	}
	return requireUidParam(ctx, field, *value)
}
