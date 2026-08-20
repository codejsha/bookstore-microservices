package restcontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/httpx"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

type RiskController struct {
	store security.RiskStore
}

func NewRiskController(store security.RiskStore) openapi.RiskApi {
	return &RiskController{store: store}
}

func (c *RiskController) RiskGetAll(ctx context.Context) (*openapi.RiskEntryListResponse, error) {
	entries, err := c.store.List(ctx)
	if err != nil {
		return nil, httpx.MapError(ctx, err)
	}
	out := make([]openapi.RiskEntryResponse, 0, len(entries))
	for _, e := range entries {
		out = append(out, toRiskEntryResponse(e))
	}
	return &openapi.RiskEntryListResponse{Entries: out}, nil
}

func (c *RiskController) RiskFlagPrincipal(
	ctx context.Context,
	uid string,
	req openapi.RiskFlagRequest,
) (*openapi.RiskEntryResponse, error) {
	if _, err := uuid.Parse(uid); err != nil {
		return nil, httpx.MapError(ctx, fmt.Errorf("%w: uid must be a UUID", httpx.ErrBadRequest))
	}
	level, ok := security.RiskLevelFromString(string(req.Level))
	if !ok {
		return nil, httpx.MapError(ctx, fmt.Errorf("%w: level must be restrict or block", httpx.ErrBadRequest))
	}
	ttl := constant.RiskDefaultTTL
	if req.TtlSeconds != nil && *req.TtlSeconds > 0 {
		ttl = time.Duration(*req.TtlSeconds) * time.Second
	}
	if ttl > constant.RiskMaxTTL {
		return nil, httpx.MapError(ctx, fmt.Errorf("%w: ttl_seconds exceeds the 30-day maximum", httpx.ErrBadRequest))
	}

	now := time.Now()
	entry := security.RiskEntry{
		Sub:       uid,
		Level:     level,
		Reason:    req.Reason,
		FlaggedBy: httpx.CallerUidFromContext(ctx),
		FlaggedAt: now,
		ExpiresAt: now.Add(ttl),
	}
	if err := c.store.Flag(ctx, entry, ttl); err != nil {
		return nil, httpx.MapError(ctx, err)
	}
	logrus.WithContext(ctx).
		WithField("user_uid", uid).WithField("level", string(level)).
		WithField("flagged_by", entry.FlaggedBy).Warn("principal flagged in risk store")
	resp := toRiskEntryResponse(entry)
	return &resp, nil
}

func (c *RiskController) RiskUnflagPrincipal(ctx context.Context, uid string) error {
	if err := c.store.Unflag(ctx, uid); err != nil {
		return httpx.MapError(ctx, err)
	}
	logrus.WithContext(ctx).WithField("user_uid", uid).Info("principal unflagged in risk store")
	return nil
}

func toRiskEntryResponse(e security.RiskEntry) openapi.RiskEntryResponse {
	var flaggedBy *string
	if e.FlaggedBy != "" {
		flaggedBy = &e.FlaggedBy
	}
	return openapi.RiskEntryResponse{
		UserUid:   e.Sub,
		Level:     openapi.RiskLevel(e.Level),
		Reason:    e.Reason,
		FlaggedBy: flaggedBy,
		FlaggedAt: e.FlaggedAt.UTC(),
		ExpiresAt: e.ExpiresAt.UTC(),
	}
}
