package security

import "context"

type SessionRevoker interface {
	Revoke(ctx context.Context, sub string) error
}

type RevocationChecker interface {
	IsRevoked(ctx context.Context, sub string, tokenIssuedAt int64) (bool, error)
}
