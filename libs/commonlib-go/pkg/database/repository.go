package database

import (
	"context"

	"gorm.io/gorm/clause"
)

func NewBaseRepository[T any](dataSource *DataSource) *BaseRepository[T] {
	return &BaseRepository[T]{
		dataSource: dataSource,
	}
}

type BaseRepository[T any] struct {
	dataSource *DataSource
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, clauses []clause.Expression) (int64, []T, error) {
	var items []T
	query := r.dataSource.db.WithContext(ctx)
	if len(clauses) > 0 {
		query = query.Clauses(clauses...)
	}
	err := query.Find(&items).Error
	return int64(len(items)), items, err
}

func (r *BaseRepository[T]) Find(ctx context.Context, id int64) (T, error) {
	var item T
	result := r.dataSource.db.WithContext(ctx).First(&item, id)
	if result.Error != nil {
		return item, result.Error
	}
	return item, nil
}

func (r *BaseRepository[T]) Save(ctx context.Context, item *T) error {
	return r.dataSource.db.WithContext(ctx).Save(item).Error
}

func (r *BaseRepository[T]) Create(ctx context.Context, item *T) error {
	return r.dataSource.db.WithContext(ctx).Create(item).Error
}

func (r *BaseRepository[T]) CreateInBatches(ctx context.Context, items []T, size int) error {
	return r.dataSource.db.WithContext(ctx).CreateInBatches(items, size).Error
}

func (r *BaseRepository[T]) Updates(ctx context.Context, item *T) error {
	return r.dataSource.db.WithContext(ctx).Updates(item).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id int64) error {
	return r.dataSource.db.WithContext(ctx).Delete(new(T), id).Error
}
