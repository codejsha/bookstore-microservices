package cache

import (
	"context"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/option"
)

type cachingWorkRepo struct {
	inner repo.WorkRepo
	cache *WorkListCache
}

var _ repo.WorkRepo = (*cachingWorkRepo)(nil)

func NewCachingWorkRepo(inner repo.WorkRepo, cache *WorkListCache) repo.WorkRepo {
	if !cache.enabled() {
		return inner
	}
	return &cachingWorkRepo{inner: inner, cache: cache}
}

func (r *cachingWorkRepo) FindAll(ctx context.Context, opt option.WorkQueryOption) (int64, []*repo.WorkResult, error) {
	if total, results, ok := r.cache.Lookup(ctx, opt); ok {
		return total, results, nil
	}
	total, results, err := r.inner.FindAll(ctx, opt)
	if err != nil {
		return 0, nil, err
	}
	r.cache.Save(ctx, opt, total, results)
	return total, results, nil
}

func (r *cachingWorkRepo) FindOne(ctx context.Context, id int64) (*repo.WorkResult, error) {
	return r.inner.FindOne(ctx, id)
}

func (r *cachingWorkRepo) FindByUid(ctx context.Context, uid string) (*repo.WorkResult, error) {
	return r.inner.FindByUid(ctx, uid)
}

func (r *cachingWorkRepo) Create(ctx context.Context, cmd command.WorkCreateCommand) (int64, error) {
	return r.inner.Create(ctx, cmd)
}

func (r *cachingWorkRepo) Update(ctx context.Context, id int64, cmd command.WorkUpdateCommand) error {
	return r.inner.Update(ctx, id, cmd)
}
