package constant

import "time"

const (
	UserCacheTTL       = 5 * time.Minute
	IntrospectCacheTTL = 10 * time.Second
	RevocationTTL      = 10 * time.Minute
	RevocationSkew     = 30 * time.Second
	RiskDefaultTTL     = 24 * time.Hour
	RiskMaxTTL         = 30 * 24 * time.Hour
)
