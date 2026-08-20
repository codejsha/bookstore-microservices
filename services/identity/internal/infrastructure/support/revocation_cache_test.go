package support

import (
	"testing"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

func TestTokenRevoked(t *testing.T) {
	const cutoff int64 = 1_000_000
	skew := int64(constant.RevocationSkew.Seconds())

	cases := []struct {
		name          string
		tokenIssuedAt int64
		want          bool
	}{
		{
			name:          "token issued well before the cutoff is revoked",
			tokenIssuedAt: cutoff - 60,
			want:          true,
		},
		{
			name:          "token issued exactly at the cutoff is still revoked within the skew window",
			tokenIssuedAt: cutoff,
			want:          true,
		},
		{
			name:          "token issued inside the skew window is revoked (fail-safe)",
			tokenIssuedAt: cutoff + skew - 1,
			want:          true,
		},
		{
			name:          "token issued past the skew window is honoured — this is the re-login case",
			tokenIssuedAt: cutoff + skew + 1,
			want:          false,
		},
		{
			name:          "unknown issued-at is treated as pre-cutoff and revoked",
			tokenIssuedAt: 0,
			want:          true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tokenRevoked(cutoff, tc.tokenIssuedAt); got != tc.want {
				t.Fatalf("tokenRevoked(%d, %d) = %v, want %v", cutoff, tc.tokenIssuedAt, got, tc.want)
			}
		})
	}
}
