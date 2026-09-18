package support

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
)

func TestRiskEntry_WhenPayloadFromFlinkAutoFlag_UnmarshalsEveryField(t *testing.T) {
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

type fakeRiskBatchStore struct {
	subs      []string
	values    map[string]string
	mgetCalls [][]string
	removed   []string
}

func (f *fakeRiskBatchStore) SMembers(context.Context, string) ([]string, error) {
	return f.subs, nil
}

func (f *fakeRiskBatchStore) MGet(_ context.Context, keys ...string) ([]any, error) {
	f.mgetCalls = append(f.mgetCalls, keys)
	values := make([]any, len(keys))
	for i, key := range keys {
		if raw, ok := f.values[key]; ok {
			values[i] = raw
		}
	}
	return values, nil
}

func (f *fakeRiskBatchStore) SRem(_ context.Context, _ string, members ...string) error {
	f.removed = append(f.removed, members...)
	return nil
}

func riskEntryJSON(t *testing.T, sub string, level security.RiskLevel) string {
	t.Helper()
	raw, err := json.Marshal(security.RiskEntry{
		Sub:       sub,
		Level:     level,
		Reason:    "test",
		FlaggedBy: "test",
		FlaggedAt: time.Unix(0, 0).UTC(),
	})
	if err != nil {
		t.Fatalf("marshal entry: %v", err)
	}
	return string(raw)
}

func TestRiskCacheList_manyFlaggedSubjects_fetchedInTwoBatches(t *testing.T) {
	store := &fakeRiskBatchStore{
		subs: []string{"a", "b", "c"},
		values: map[string]string{
			riskKey("a"):     riskEntryJSON(t, "a", security.RiskLevelBlock),
			autoRiskKey("b"): riskEntryJSON(t, "b", security.RiskLevelRestrict),
		},
	}
	r := &RiskCache{batch: store}

	entries, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("entries = %+v, want the manual and the auto flag", entries)
	}
	if len(store.mgetCalls) != 2 {
		t.Fatalf("mget calls = %d, want 2 batches instead of one round trip per subject", len(store.mgetCalls))
	}
	wantFirst := []string{riskKey("a"), riskKey("b"), riskKey("c")}
	if len(store.mgetCalls[0]) != len(wantFirst) {
		t.Fatalf("first batch = %v, want %v", store.mgetCalls[0], wantFirst)
	}
	wantSecond := []string{autoRiskKey("b"), autoRiskKey("c")}
	if len(store.mgetCalls[1]) != len(wantSecond) {
		t.Fatalf("second batch = %v, want only the subjects missing a manual flag %v", store.mgetCalls[1], wantSecond)
	}
	if len(store.removed) != 1 || store.removed[0] != "c" {
		t.Fatalf("removed = %v, want the expired subject c pruned from the index", store.removed)
	}
}

func TestRiskCacheList_allSubjectsFlaggedManually_autoBatchSkipped(t *testing.T) {
	store := &fakeRiskBatchStore{
		subs: []string{"a", "b"},
		values: map[string]string{
			riskKey("a"): riskEntryJSON(t, "a", security.RiskLevelBlock),
			riskKey("b"): riskEntryJSON(t, "b", security.RiskLevelRestrict),
		},
	}
	r := &RiskCache{batch: store}

	entries, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("entries = %+v, want 2", entries)
	}
	if len(store.mgetCalls) != 1 {
		t.Fatalf("mget calls = %d, want a single batch when no subject falls through to the auto key", len(store.mgetCalls))
	}
	if len(store.removed) != 0 {
		t.Fatalf("removed = %v, want nothing pruned", store.removed)
	}
}

func TestRiskCacheList_emptyIndex_noValueLookup(t *testing.T) {
	store := &fakeRiskBatchStore{}
	r := &RiskCache{batch: store}

	entries, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("entries = %+v, want none", entries)
	}
	if len(store.mgetCalls) != 0 {
		t.Fatalf("mget calls = %d, want none for an empty index", len(store.mgetCalls))
	}
}

func TestDecodeRiskValues_valueCountMismatch_reportsError(t *testing.T) {
	if _, _, err := decodeRiskValues([]string{"a", "b"}, []any{nil}); err == nil {
		t.Fatal("a short mget reply must not be silently truncated")
	}
}
