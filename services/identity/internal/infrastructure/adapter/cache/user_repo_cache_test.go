package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

type fakeStore struct {
	entry   *repo.UserResult
	getErr  error
	setKey  string
	setTTL  time.Duration
	sets    int
	delKeys []string
}

func (f *fakeStore) Get(_ context.Context, _ string, dest interface{}) error {
	if f.getErr != nil {
		return f.getErr
	}
	*dest.(*repo.UserResult) = *f.entry
	return nil
}

func (f *fakeStore) Set(_ context.Context, key string, _ interface{}, ttl time.Duration) error {
	f.sets++
	f.setKey = key
	f.setTTL = ttl
	return nil
}

func (f *fakeStore) Del(_ context.Context, keys ...string) error {
	f.delKeys = append(f.delKeys, keys...)
	return nil
}

type fakeRepo struct {
	user      *repo.UserResult
	findErr   error
	finds     int
	writeErr  error
	writes    int
	lastIdpId string
}

func (f *fakeRepo) FindByUid(context.Context, string) (*repo.UserResult, error) {
	f.finds++
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.user, nil
}
func (f *fakeRepo) FindAll(context.Context, option.UserQueryOption) (int64, []*repo.UserResult, error) {
	return 0, nil, nil
}
func (f *fakeRepo) FetchByEmail(context.Context, string) ([]*repo.UserResult, error) {
	return nil, nil
}
func (f *fakeRepo) Upsert(_ context.Context, p repo.UserUpsert) (*repo.UserResult, error) {
	f.writes++
	f.lastIdpId = p.IdpUid
	if f.writeErr != nil {
		return nil, f.writeErr
	}
	return f.user, nil
}
func (f *fakeRepo) UpdateProfile(_ context.Context, idpUid string, _ repo.UserProfileUpdate) error {
	f.writes++
	f.lastIdpId = idpUid
	return f.writeErr
}
func (f *fakeRepo) UpdateStatus(_ context.Context, idpUid, _ string) error {
	f.writes++
	f.lastIdpId = idpUid
	return f.writeErr
}
func (f *fakeRepo) UpdateRoles(_ context.Context, idpUid string, _ []string) error {
	f.writes++
	f.lastIdpId = idpUid
	return f.writeErr
}
func (f *fakeRepo) SoftDelete(_ context.Context, idpUid string) error {
	f.writes++
	f.lastIdpId = idpUid
	return f.writeErr
}

func TestFindByUid_WhenEntryCached_SkipsRepo(t *testing.T) {
	store := &fakeStore{entry: &repo.UserResult{Email: "cached@example.com"}}
	inner := &fakeRepo{}
	r := NewCachingUserRepo(inner, store)

	user, err := r.FindByUid(context.Background(), "u1")
	if err != nil {
		t.Fatalf("FindByUid: %v", err)
	}
	if user.Email != "cached@example.com" {
		t.Fatalf("email = %q, want the cached entry", user.Email)
	}
	if inner.finds != 0 {
		t.Fatalf("inner finds = %d, want 0 on a cache hit", inner.finds)
	}
}

func TestFindByUid_WhenEntryMissing_ReadsRepoAndPopulatesCache(t *testing.T) {
	store := &fakeStore{getErr: errors.New("miss")}
	inner := &fakeRepo{user: &repo.UserResult{Email: "db@example.com"}}
	r := NewCachingUserRepo(inner, store)

	user, err := r.FindByUid(context.Background(), "u1")
	if err != nil {
		t.Fatalf("FindByUid: %v", err)
	}
	if user.Email != "db@example.com" {
		t.Fatalf("email = %q, want the DB row", user.Email)
	}
	if inner.finds != 1 {
		t.Fatalf("inner finds = %d, want 1", inner.finds)
	}
	if store.sets != 1 || store.setKey != userKey("u1") {
		t.Fatalf("sets = %d, key = %q, want 1 write to %q", store.sets, store.setKey, userKey("u1"))
	}
	if store.setTTL != constant.UserCacheTTL {
		t.Fatalf("ttl = %v, want %v", store.setTTL, constant.UserCacheTTL)
	}
}

func TestFindByUid_WhenRepoFails_DoesNotCacheNegativeEntry(t *testing.T) {
	store := &fakeStore{getErr: errors.New("miss")}
	inner := &fakeRepo{findErr: errors.New("record not found")}
	r := NewCachingUserRepo(inner, store)

	if _, err := r.FindByUid(context.Background(), "u1"); err == nil {
		t.Fatal("expected the repo error to surface")
	}
	if store.sets != 0 {
		t.Fatalf("sets = %d, want 0 (a negative entry would outlive the user's creation)", store.sets)
	}
}

// Every write method is driven through the same subtest, named after the method under test.
func TestWrites_WhenWriteSucceeds_InvalidateCachedUser(t *testing.T) {
	ctx := context.Background()
	cases := map[string]func(repo.UserRepo) error{
		"Upsert": func(r repo.UserRepo) error {
			_, err := r.Upsert(ctx, repo.UserUpsert{IdpUid: "u1"})
			return err
		},
		"UpdateProfile": func(r repo.UserRepo) error {
			return r.UpdateProfile(ctx, "u1", repo.UserProfileUpdate{})
		},
		"UpdateStatus": func(r repo.UserRepo) error { return r.UpdateStatus(ctx, "u1", "ACTIVE") },
		"UpdateRoles":  func(r repo.UserRepo) error { return r.UpdateRoles(ctx, "u1", []string{"USER"}) },
		"SoftDelete":   func(r repo.UserRepo) error { return r.SoftDelete(ctx, "u1") },
	}

	for name, write := range cases {
		t.Run(name, func(t *testing.T) {
			store := &fakeStore{}
			inner := &fakeRepo{user: &repo.UserResult{}}
			if err := write(NewCachingUserRepo(inner, store)); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if len(store.delKeys) != 1 || store.delKeys[0] != userKey("u1") {
				t.Fatalf("delKeys = %v, want [%q]", store.delKeys, userKey("u1"))
			}
		})
	}
}

func TestWrites_WhenWriteFails_LeaveCachedUserInPlace(t *testing.T) {
	store := &fakeStore{}
	inner := &fakeRepo{writeErr: errors.New("db down")}
	r := NewCachingUserRepo(inner, store)

	if err := r.UpdateRoles(context.Background(), "u1", []string{"ADMIN"}); err == nil {
		t.Fatal("expected the repo error to surface")
	}
	if len(store.delKeys) != 0 {
		t.Fatalf("delKeys = %v, want none — the row never changed", store.delKeys)
	}
}
