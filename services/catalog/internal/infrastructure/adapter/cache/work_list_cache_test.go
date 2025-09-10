package cache

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/codejsha/shared-library-go/pkg/pagination"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

type fakeStore struct {
	data  map[string]string
	epoch int64
	errOn bool
}

func newFakeStore() *fakeStore { return &fakeStore{data: map[string]string{}} }

func (s *fakeStore) Get(_ context.Context, key string, dest interface{}) error {
	if s.errOn {
		return errors.New("boom")
	}
	v, ok := s.data[key]
	if !ok {
		return errors.New("nil")
	}
	return json.Unmarshal([]byte(v), dest)
}

func (s *fakeStore) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	if s.errOn {
		return errors.New("boom")
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.data[key] = string(b)
	return nil
}

func (s *fakeStore) Incr(_ context.Context, key string) (int64, error) {
	if s.errOn {
		return 0, errors.New("boom")
	}
	s.epoch++
	s.data[key] = itoa(s.epoch)
	return s.epoch, nil
}

func itoa(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func workOpt(title string, page int32) option.WorkQueryOption {
	t := title
	return option.NewWorkQueryOption(
		option.WorkQueryOption{}.WithTitle(&t),
		option.WorkQueryOption{}.WithPage(pagination.NewPageOption(nil, &page, nil)),
	)
}

// ─── epoch-key construction ────────────────────────────────────────────────

func TestHashWorkQuery_StableAndSensitive(t *testing.T) {
	a := hashWorkQuery(workOpt("hobbit", 1))
	if a != hashWorkQuery(workOpt("hobbit", 1)) {
		t.Error("hash not stable for identical queries")
	}
	if a == hashWorkQuery(workOpt("hobbit", 2)) {
		t.Error("hash insensitive to page")
	}
	if a == hashWorkQuery(workOpt("dune", 1)) {
		t.Error("hash insensitive to title filter")
	}
}

func TestListKey_EmbedsEpochAndPrefix(t *testing.T) {
	store := newFakeStore()
	c := NewWorkListCache(store)
	opt := workOpt("hobbit", 1)

	k0 := c.listKey(context.Background(), opt)
	if !strings.HasPrefix(k0, listKeyPrefix+"0:") {
		t.Errorf("key = %q, want epoch-0 prefix %q", k0, listKeyPrefix+"0:")
	}

	if err := c.Invalidate(context.Background()); err != nil {
		t.Fatalf("invalidate: %v", err)
	}
	k1 := c.listKey(context.Background(), opt)
	if k1 == k0 {
		t.Error("key unchanged after epoch bump")
	}
	if !strings.HasPrefix(k1, listKeyPrefix+"1:") {
		t.Errorf("key = %q, want epoch-1 prefix", k1)
	}
}

// ─── round-trip + invalidation ─────────────────────────────────────────────

func TestWorkListCache_SaveThenLookup(t *testing.T) {
	c := NewWorkListCache(newFakeStore())
	opt := workOpt("hobbit", 1)
	want := []*repo.WorkResult{{Uid: "w-1", Title: "Hobbit"}}

	c.Save(context.Background(), opt, 1, want)
	total, got, ok := c.Lookup(context.Background(), opt)
	if !ok || total != 1 || len(got) != 1 || got[0].Uid != "w-1" {
		t.Fatalf("lookup = (%d, %+v, %v), want cached hit", total, got, ok)
	}

	if err := c.Invalidate(context.Background()); err != nil {
		t.Fatalf("invalidate: %v", err)
	}
	if _, _, ok := c.Lookup(context.Background(), opt); ok {
		t.Error("lookup hit after invalidation, want miss")
	}
}

// ─── graceful degradation ──────────────────────────────────────────────────

func TestWorkListCache_DisabledIsNoOp(t *testing.T) {
	c := NewWorkListCache(nil)
	opt := workOpt("hobbit", 1)

	if _, _, ok := c.Lookup(context.Background(), opt); ok {
		t.Error("disabled cache reported a hit")
	}
	c.Save(context.Background(), opt, 1, []*repo.WorkResult{{Uid: "w-1"}})
	if err := c.Invalidate(context.Background()); err != nil {
		t.Errorf("disabled Invalidate returned %v, want nil", err)
	}
}

func TestWorkListCache_TransportFailuresDegrade(t *testing.T) {
	c := NewWorkListCache(&fakeStore{data: map[string]string{}, errOn: true})
	opt := workOpt("hobbit", 1)

	if _, _, ok := c.Lookup(context.Background(), opt); ok {
		t.Error("Lookup reported hit despite store error")
	}
	c.Save(context.Background(), opt, 1, []*repo.WorkResult{{Uid: "w-1"}})
	if err := c.Invalidate(context.Background()); err == nil {
		t.Error("Invalidate swallowed store error, want it surfaced for logging")
	}
}

// ─── decorator ─────────────────────────────────────────────────────────────

type stubWorkRepo struct {
	calls int
	total int64
	rows  []*repo.WorkResult
	err   error
}

func (s *stubWorkRepo) FindAll(context.Context, option.WorkQueryOption) (int64, []*repo.WorkResult, error) {
	s.calls++
	return s.total, s.rows, s.err
}
func (s *stubWorkRepo) FindOne(context.Context, int64) (*repo.WorkResult, error) { return nil, nil }
func (s *stubWorkRepo) FindByUid(context.Context, string) (*repo.WorkResult, error) {
	return nil, nil
}
func (s *stubWorkRepo) Create(context.Context, command.WorkCreateCommand) (int64, error) {
	return 0, nil
}
func (s *stubWorkRepo) Update(context.Context, int64, command.WorkUpdateCommand) error { return nil }

func TestNewCachingWorkRepo_DisabledReturnsInner(t *testing.T) {
	inner := &stubWorkRepo{}
	if got := NewCachingWorkRepo(inner, NewWorkListCache(nil)); got != repo.WorkRepo(inner) {
		t.Error("disabled cache should return the inner repo unwrapped")
	}
}

func TestCachingWorkRepo_CachesAndDegrades(t *testing.T) {
	inner := &stubWorkRepo{total: 2, rows: []*repo.WorkResult{{Uid: "w-1"}, {Uid: "w-2"}}}
	dec := NewCachingWorkRepo(inner, NewWorkListCache(newFakeStore()))
	opt := workOpt("hobbit", 1)

	if total, _, err := dec.FindAll(context.Background(), opt); err != nil || total != 2 {
		t.Fatalf("first FindAll = (%d, %v)", total, err)
	}
	if _, _, err := dec.FindAll(context.Background(), opt); err != nil {
		t.Fatalf("second FindAll: %v", err)
	}
	if inner.calls != 1 {
		t.Errorf("inner FindAll calls = %d, want 1 (second served from cache)", inner.calls)
	}
}

func TestCachingWorkRepo_StoreErrorFallsThroughToDB(t *testing.T) {
	inner := &stubWorkRepo{total: 1, rows: []*repo.WorkResult{{Uid: "w-1"}}}
	dec := NewCachingWorkRepo(inner, NewWorkListCache(&fakeStore{data: map[string]string{}, errOn: true}))

	total, rows, err := dec.FindAll(context.Background(), workOpt("hobbit", 1))
	if err != nil || total != 1 || len(rows) != 1 {
		t.Fatalf("FindAll = (%d, %+v, %v), want DB passthrough", total, rows, err)
	}
	if inner.calls != 1 {
		t.Errorf("inner calls = %d, want 1", inner.calls)
	}
}
