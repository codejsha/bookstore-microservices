package support

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
)

func newFakeValkeyRiskCache(t *testing.T) (*RiskCache, *fakeValkey) {
	t.Helper()
	c, valkey := newFakeValkeyClient(t)
	return NewRiskCache(c), valkey
}

func riskEntry(sub string, level security.RiskLevel) security.RiskEntry {
	return security.RiskEntry{
		Sub:       sub,
		Level:     level,
		Reason:    "test",
		FlaggedBy: "admin-1",
		FlaggedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		ExpiresAt: time.Date(2026, 10, 1, 0, 15, 0, 0, time.UTC),
	}
}

func entrySubs(entries []security.RiskEntry) []string {
	subs := make([]string, len(entries))
	for i, entry := range entries {
		subs[i] = entry.Sub
	}
	sort.Strings(subs)
	return subs
}

func TestRiskCacheFlag_entryGiven_storesItWithTtlAndIndexesTheSubject(t *testing.T) {
	r, valkey := newFakeValkeyRiskCache(t)
	entry := riskEntry("sub-1", security.RiskLevelBlock)

	if err := r.Flag(context.Background(), entry, 15*time.Minute); err != nil {
		t.Fatalf("flag: %v", err)
	}

	if _, ok := valkey.value("identity:risk:sub-1"); !ok {
		t.Fatal("identity:risk:sub-1 missing, want the entry stored under the manual key")
	}
	if got := valkey.ttl("identity:risk:sub-1"); got != 15*time.Minute {
		t.Errorf("ttl = %v, want 15m", got)
	}
	if want := []string{"sub-1"}; !reflect.DeepEqual(valkey.members(riskIndexKey), want) {
		t.Errorf("index = %v, want %v", valkey.members(riskIndexKey), want)
	}
	got, err := r.Check(context.Background(), "sub-1")
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if got == nil || *got != entry {
		t.Fatalf("checked = %+v, want %+v", got, entry)
	}
}

func TestRiskCacheUnflag_manualAndAutoFlagsPresent_removesBothAndTheIndexEntry(t *testing.T) {
	r, valkey := newFakeValkeyRiskCache(t)
	ctx := context.Background()
	if err := r.Flag(ctx, riskEntry("sub-1", security.RiskLevelBlock), time.Minute); err != nil {
		t.Fatalf("flag sub-1: %v", err)
	}
	if err := r.Flag(ctx, riskEntry("sub-2", security.RiskLevelRestrict), time.Minute); err != nil {
		t.Fatalf("flag sub-2: %v", err)
	}
	valkey.put("identity:risk:auto:sub-1", riskEntryJSON(t, "sub-1", security.RiskLevelRestrict))

	if err := r.Unflag(ctx, "sub-1"); err != nil {
		t.Fatalf("unflag: %v", err)
	}

	for _, key := range []string{"identity:risk:sub-1", "identity:risk:auto:sub-1"} {
		if _, ok := valkey.value(key); ok {
			t.Errorf("%s still stored, want it removed", key)
		}
	}
	if want := []string{"sub-2"}; !reflect.DeepEqual(valkey.members(riskIndexKey), want) {
		t.Errorf("index = %v, want %v", valkey.members(riskIndexKey), want)
	}
	got, err := r.Check(ctx, "sub-1")
	if err != nil || got != nil {
		t.Fatalf("check = (%+v, %v), want the subject cleared", got, err)
	}
	if other, err := r.Check(ctx, "sub-2"); err != nil || other == nil {
		t.Fatalf("check sub-2 = (%+v, %v), want the other subject still flagged", other, err)
	}
}

func TestRiskCacheCheck_onlyAutoFlagPresent_fallsBackToTheAutoEntry(t *testing.T) {
	r, valkey := newFakeValkeyRiskCache(t)
	valkey.put("identity:risk:auto:sub-1", riskEntryJSON(t, "sub-1", security.RiskLevelRestrict))

	got, err := r.Check(context.Background(), "sub-1")
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if got == nil || got.Level != security.RiskLevelRestrict || got.Sub != "sub-1" {
		t.Fatalf("checked = %+v, want the auto restrict entry", got)
	}
}

func TestRiskCacheCheck_manualAndAutoFlagsPresent_manualFlagWins(t *testing.T) {
	r, valkey := newFakeValkeyRiskCache(t)
	valkey.put("identity:risk:sub-1", riskEntryJSON(t, "sub-1", security.RiskLevelBlock))
	valkey.put("identity:risk:auto:sub-1", riskEntryJSON(t, "sub-1", security.RiskLevelRestrict))

	got, err := r.Check(context.Background(), "sub-1")
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if got == nil || got.Level != security.RiskLevelBlock {
		t.Fatalf("checked = %+v, want the manual block entry", got)
	}
	if lookups := valkey.commandsNamed("GET"); len(lookups) != 1 {
		t.Errorf("lookups = %v, want the auto key left unread", lookups)
	}
}

func TestRiskCacheCheck_subjectNotFlagged_nilEntryWithoutError(t *testing.T) {
	r, _ := newFakeValkeyRiskCache(t)

	got, err := r.Check(context.Background(), "sub-1")

	if err != nil || got != nil {
		t.Fatalf("check = (%+v, %v), want (nil, nil)", got, err)
	}
}

func TestRiskCacheCheck_storedEntryCorrupt_reportsErrorInsteadOfClearingTheSubject(t *testing.T) {
	r, valkey := newFakeValkeyRiskCache(t)
	valkey.put("identity:risk:sub-1", "not-json")

	got, err := r.Check(context.Background(), "sub-1")

	if err == nil || errors.Is(err, redis.Nil) {
		t.Fatalf("err = %v, want a decode error", err)
	}
	if got != nil {
		t.Errorf("checked = %+v, want nil", got)
	}
}

func TestRiskCache_storeUnavailable_reportsError(t *testing.T) {
	ctx := context.Background()
	calls := map[string]func(r *RiskCache) error{
		"flag":   func(r *RiskCache) error { return r.Flag(ctx, riskEntry("sub-1", security.RiskLevelBlock), time.Minute) },
		"unflag": func(r *RiskCache) error { return r.Unflag(ctx, "sub-1") },
		"check": func(r *RiskCache) error {
			_, err := r.Check(ctx, "sub-1")
			return err
		},
		"list": func(r *RiskCache) error {
			_, err := r.List(ctx)
			return err
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			r, _ := newFakeValkeyRiskCache(t)
			if err := r.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}

			if err := call(r); !errors.Is(err, redis.ErrClosed) {
				t.Fatalf("err = %v, want redis.ErrClosed", err)
			}
		})
	}
}

func TestRiskCacheList_manualAutoAndExpiredSubjects_listsLiveFlagsAndPrunesTheExpired(t *testing.T) {
	r, valkey := newFakeValkeyRiskCache(t)
	ctx := context.Background()
	if err := r.Flag(ctx, riskEntry("manual", security.RiskLevelBlock), time.Minute); err != nil {
		t.Fatalf("flag: %v", err)
	}
	valkey.put("identity:risk:auto:auto", riskEntryJSON(t, "auto", security.RiskLevelRestrict))
	if err := r.SAdd(ctx, riskIndexKey, "auto", "expired"); err != nil {
		t.Fatalf("index: %v", err)
	}

	entries, err := r.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if want := []string{"auto", "manual"}; !reflect.DeepEqual(entrySubs(entries), want) {
		t.Fatalf("listed = %v, want %v", entrySubs(entries), want)
	}
	if want := []string{"auto", "manual"}; !reflect.DeepEqual(valkey.members(riskIndexKey), want) {
		t.Errorf("index = %v, want the expired subject pruned leaving %v", valkey.members(riskIndexKey), want)
	}
}

type scriptedRiskBatchStore struct {
	subs       []string
	membersErr error
	replies    [][]any
	mgetErrs   []error
	calls      int
	removed    []string
}

func (s *scriptedRiskBatchStore) SMembers(context.Context, string) ([]string, error) {
	return s.subs, s.membersErr
}

func (s *scriptedRiskBatchStore) MGet(context.Context, ...string) ([]any, error) {
	i := s.calls
	s.calls++
	var err error
	if i < len(s.mgetErrs) {
		err = s.mgetErrs[i]
	}
	if err != nil {
		return nil, err
	}
	return s.replies[i], nil
}

func (s *scriptedRiskBatchStore) SRem(_ context.Context, _ string, members ...string) error {
	s.removed = append(s.removed, members...)
	return nil
}

func TestRiskCacheList_indexLookupFails_reportsError(t *testing.T) {
	failure := errors.New("valkey down")
	r := &RiskCache{batch: &scriptedRiskBatchStore{membersErr: failure}}

	entries, err := r.List(context.Background())

	if !errors.Is(err, failure) || entries != nil {
		t.Fatalf("list = (%+v, %v), want the index failure", entries, err)
	}
}

func TestRiskCacheList_manualBatchFails_reportsError(t *testing.T) {
	failure := errors.New("valkey down")
	r := &RiskCache{batch: &scriptedRiskBatchStore{subs: []string{"a"}, mgetErrs: []error{failure}}}

	entries, err := r.List(context.Background())

	if !errors.Is(err, failure) || entries != nil {
		t.Fatalf("list = (%+v, %v), want the batch failure", entries, err)
	}
}

func TestRiskCacheList_autoBatchFails_reportsErrorInsteadOfAPartialList(t *testing.T) {
	failure := errors.New("valkey down")
	store := &scriptedRiskBatchStore{
		subs:     []string{"a", "b"},
		replies:  [][]any{{riskEntryJSON(t, "a", security.RiskLevelBlock), nil}},
		mgetErrs: []error{nil, failure},
	}
	r := &RiskCache{batch: store}

	entries, err := r.List(context.Background())

	if !errors.Is(err, failure) || entries != nil {
		t.Fatalf("list = (%+v, %v), want the batch failure", entries, err)
	}
	if len(store.removed) != 0 {
		t.Errorf("removed = %v, want nothing pruned on a failed lookup", store.removed)
	}
}

func TestRiskCacheList_autoEntryCorrupt_reportsErrorWithoutPruning(t *testing.T) {
	store := &scriptedRiskBatchStore{
		subs:    []string{"a", "b"},
		replies: [][]any{{nil, nil}, {"not-json", nil}},
	}
	r := &RiskCache{batch: store}

	_, err := r.List(context.Background())

	if err == nil || !strings.Contains(err.Error(), "subject a") {
		t.Fatalf("err = %v, want the corrupt subject named", err)
	}
	if len(store.removed) != 0 {
		t.Errorf("removed = %v, want nothing pruned", store.removed)
	}
}

func TestRiskCacheList_manualEntryCorrupt_reportsError(t *testing.T) {
	store := &scriptedRiskBatchStore{subs: []string{"a"}, replies: [][]any{{"not-json"}}}
	r := &RiskCache{batch: store}

	_, err := r.List(context.Background())

	if err == nil {
		t.Fatal("err = nil, want a decode error")
	}
	if store.calls != 1 {
		t.Errorf("batches = %d, want the auto lookup skipped", store.calls)
	}
}

func TestDecodeRiskValues_mixedReplyTypes_decodesStringsAndBytesAndCollectsMissing(t *testing.T) {
	values := []any{
		riskEntryJSON(t, "a", security.RiskLevelBlock),
		nil,
		[]byte(riskEntryJSON(t, "c", security.RiskLevelRestrict)),
	}

	entries, missing, err := decodeRiskValues([]string{"a", "b", "c"}, values)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if want := []string{"a", "c"}; !reflect.DeepEqual(entrySubs(entries), want) {
		t.Errorf("entries = %v, want %v", entrySubs(entries), want)
	}
	if want := []string{"b"}; !reflect.DeepEqual(missing, want) {
		t.Errorf("missing = %v, want %v", missing, want)
	}
}

func TestDecodeRiskValues_unexpectedReplyType_reportsErrorNamingTheSubject(t *testing.T) {
	_, _, err := decodeRiskValues([]string{"a"}, []any{int64(42)})

	if err == nil || !strings.Contains(err.Error(), "subject a") {
		t.Fatalf("err = %v, want the subject named", err)
	}
}

func TestProvideRiskPorts_sharedRiskCache_checkerSeesWhatTheStoreFlagged(t *testing.T) {
	r, _ := newFakeValkeyRiskCache(t)
	store, checker := ProvideRiskStore(r), ProvideRiskChecker(r)
	entry := riskEntry("sub-1", security.RiskLevelBlock)

	if err := store.Flag(context.Background(), entry, time.Minute); err != nil {
		t.Fatalf("flag: %v", err)
	}
	got, err := checker.Check(context.Background(), "sub-1")
	if err != nil {
		t.Fatalf("check: %v", err)
	}

	if got == nil || *got != entry {
		t.Fatalf("checked = %+v, want %+v", got, entry)
	}
}
