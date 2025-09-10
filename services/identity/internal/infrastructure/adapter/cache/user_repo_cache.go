package cache

import (
	"context"
	"time"

	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/identity/internal/domain/model/option"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

func userKey(uid string) string { return "identity:user:uid:" + uid }

type Store interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type cachingUserRepo struct {
	inner repo.UserRepo
	store Store
}

var _ repo.UserRepo = (*cachingUserRepo)(nil)

func NewCachingUserRepo(inner repo.UserRepo, store Store) repo.UserRepo {
	return &cachingUserRepo{inner: inner, store: store}
}

func (r *cachingUserRepo) FindByUid(ctx context.Context, uid string) (*repo.UserResult, error) {
	var cached repo.UserResult
	if err := r.store.Get(ctx, userKey(uid), &cached); err == nil {
		return &cached, nil
	}

	user, err := r.inner.FindByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	_ = r.store.Set(ctx, userKey(uid), user, constant.UserCacheTTL)
	return user, nil
}

func (r *cachingUserRepo) FindAll(
	ctx context.Context,
	opt option.UserQueryOption,
) (int64, []*repo.UserResult, error) {
	return r.inner.FindAll(ctx, opt)
}

func (r *cachingUserRepo) FetchByEmail(ctx context.Context, email string) ([]*repo.UserResult, error) {
	return r.inner.FetchByEmail(ctx, email)
}

func (r *cachingUserRepo) Upsert(ctx context.Context, profile repo.UserUpsert) (*repo.UserResult, error) {
	user, err := r.inner.Upsert(ctx, profile)
	if err != nil {
		return nil, err
	}
	r.invalidate(ctx, profile.IdpUid)
	return user, nil
}

func (r *cachingUserRepo) UpdateProfile(
	ctx context.Context,
	idpUid string,
	patch repo.UserProfileUpdate,
) error {
	if err := r.inner.UpdateProfile(ctx, idpUid, patch); err != nil {
		return err
	}
	r.invalidate(ctx, idpUid)
	return nil
}

func (r *cachingUserRepo) UpdateStatus(ctx context.Context, idpUid, status string) error {
	if err := r.inner.UpdateStatus(ctx, idpUid, status); err != nil {
		return err
	}
	r.invalidate(ctx, idpUid)
	return nil
}

func (r *cachingUserRepo) UpdateRoles(ctx context.Context, idpUid string, roles []string) error {
	if err := r.inner.UpdateRoles(ctx, idpUid, roles); err != nil {
		return err
	}
	r.invalidate(ctx, idpUid)
	return nil
}

func (r *cachingUserRepo) SoftDelete(ctx context.Context, idpUid string) error {
	if err := r.inner.SoftDelete(ctx, idpUid); err != nil {
		return err
	}
	r.invalidate(ctx, idpUid)
	return nil
}

func (r *cachingUserRepo) invalidate(ctx context.Context, uid string) {
	_ = r.store.Del(ctx, userKey(uid))
}
