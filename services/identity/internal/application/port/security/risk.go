package security

import (
	"context"
	"time"
)

type RiskLevel string

const (
	RiskLevelRestrict RiskLevel = "restrict"
	RiskLevelBlock    RiskLevel = "block"
)

func RiskLevelFromString(v string) (RiskLevel, bool) {
	switch RiskLevel(v) {
	case RiskLevelRestrict:
		return RiskLevelRestrict, true
	case RiskLevelBlock:
		return RiskLevelBlock, true
	default:
		return "", false
	}
}

type RiskEntry struct {
	Sub       string    `json:"sub"`
	Level     RiskLevel `json:"level"`
	Reason    string    `json:"reason"`
	FlaggedBy string    `json:"flagged_by"`
	FlaggedAt time.Time `json:"flagged_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RiskChecker interface {
	Check(ctx context.Context, sub string) (*RiskEntry, error)
}

type RiskStore interface {
	RiskChecker
	Flag(ctx context.Context, entry RiskEntry, ttl time.Duration) error
	Unflag(ctx context.Context, sub string) error
	List(ctx context.Context) ([]RiskEntry, error)
}
