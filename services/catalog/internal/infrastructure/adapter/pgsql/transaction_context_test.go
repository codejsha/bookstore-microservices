package pgsql

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
)

func newExhaustedPoolDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	held, err := sqlDB.Begin()
	if err != nil {
		t.Fatalf("hold the only pooled connection: %v", err)
	}
	t.Cleanup(func() {
		_ = held.Rollback()
		_ = sqlDB.Close()
	})
	return db
}

func assertStopsAtRequestDeadline(t *testing.T, run func(ctx context.Context) error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- run(ctx) }()
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want context deadline exceeded", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("transaction kept waiting for a pooled connection past the request deadline")
	}
}

func TestWorkRepository_poolExhausted_createStopsAtRequestDeadline(t *testing.T) {
	db := newExhaustedPoolDB(t)
	r := &workRepository{q: dao.Use(db), db: db}

	assertStopsAtRequestDeadline(t, func(ctx context.Context) error {
		_, err := r.Create(ctx, command.WorkCreateCommand{})
		return err
	})
}

func TestWorkRepository_poolExhausted_updateStopsAtRequestDeadline(t *testing.T) {
	db := newExhaustedPoolDB(t)
	r := &workRepository{q: dao.Use(db), db: db}

	assertStopsAtRequestDeadline(t, func(ctx context.Context) error {
		return r.Update(ctx, 1, command.WorkUpdateCommand{})
	})
}
