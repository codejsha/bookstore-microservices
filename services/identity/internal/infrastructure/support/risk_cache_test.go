package support

import (
	"encoding/json"
	"testing"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
)

func TestRiskEntry_UnmarshalsFlinkAutoFlagPayload(t *testing.T) {
	payload := `{"sub":"7c9e6679-7425-40de-944b-e07fc1f90ae7",` +
		`"level":"restrict",` +
		`"reason":"auto: score 3421 (2801 req, 120 errors, 70 paths in 5m)",` +
		`"flagged_by":"risk-scorer",` +
		`"flagged_at":"2026-08-20T05:00:00Z",` +
		`"expires_at":"2026-08-20T05:15:00Z"}`

	var entry security.RiskEntry
	if err := json.Unmarshal([]byte(payload), &entry); err != nil {
		t.Fatalf("the Flink risk-scorer payload must stay unmarshalable: %v", err)
	}
	if entry.Level != security.RiskLevelRestrict {
		t.Fatalf("level = %q, want restrict", entry.Level)
	}
	if entry.Sub != "7c9e6679-7425-40de-944b-e07fc1f90ae7" || entry.FlaggedBy != "risk-scorer" {
		t.Fatalf("unexpected entry: %+v", entry)
	}
	if entry.ExpiresAt.Sub(entry.FlaggedAt).Minutes() != 15 {
		t.Fatalf("expires_at - flagged_at = %v, want 15m", entry.ExpiresAt.Sub(entry.FlaggedAt))
	}
}
