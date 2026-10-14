package pgsql

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/dao"
	"github.com/codejsha/bookstore-microservices/catalog/generated/infrastructure/port/entity"
	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/repo"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/model/command"
)

const (
	testWorkUid      = "0190a000-0000-7000-8000-000000000001"
	testOtherWorkUid = "0190a000-0000-7000-8000-000000000002"
	testPublisherUid = "0190a000-0000-7000-8000-000000000003"
)

const editionTableDDL = `CREATE TABLE edition (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid TEXT NOT NULL,
	title TEXT NOT NULL,
	isbn10 TEXT, isbn13 TEXT, number_of_pages INTEGER, publish_date TEXT,
	cover_ids TEXT, cover_uids TEXT, languages TEXT, physical_format TEXT, description TEXT,
	work_id INTEGER NOT NULL,
	work_uid TEXT NOT NULL,
	publisher_id INTEGER,
	publisher_uid TEXT,
	ol_key TEXT,
	created_at DATETIME NOT NULL, updated_at DATETIME, deleted_at DATETIME,
	actor INTEGER NOT NULL, version INTEGER NOT NULL
)`

func newEditionTableDB(t *testing.T) *gorm.DB {
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
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.Exec(editionTableDDL).Error; err != nil {
		t.Fatalf("create edition table: %v", err)
	}
	return db
}

func loadEdition(t *testing.T, db *gorm.DB, id int64) entity.EditionEntity {
	t.Helper()
	var got entity.EditionEntity
	if err := db.Unscoped().Where("id = ?", id).Take(&got).Error; err != nil {
		t.Fatalf("load edition %d: %v", id, err)
	}
	return got
}

func TestNewEditionEntity_PublisherGiven_PopulatesWorkAndPublisherUids(t *testing.T) {
	now := time.Date(2026, 10, 7, 0, 0, 0, 0, time.UTC)
	got := newEditionEntity(
		command.EditionCreateCommand{Title: "1st"},
		repo.EditionWorkRef{Id: 50, Uid: testWorkUid},
		&repo.EditionPublisherRef{Id: 7, Uid: testPublisherUid},
		"e-uid", now,
	)
	if got.WorkId != 50 || got.WorkUid != testWorkUid {
		t.Errorf("work = %d/%q, want 50/%q", got.WorkId, got.WorkUid, testWorkUid)
	}
	if got.PublisherId == nil || *got.PublisherId != 7 || got.PublisherUid == nil || *got.PublisherUid != testPublisherUid {
		t.Errorf("publisher = %v/%v, want 7/%q", got.PublisherId, got.PublisherUid, testPublisherUid)
	}
	if got.Uid != "e-uid" || got.Title != "1st" || !got.CreatedAt.Equal(now) || got.Version != 1 {
		t.Errorf("entity = %+v", got)
	}
}

func TestNewEditionEntity_PublisherOmitted_LeavesPublisherColumnsNil(t *testing.T) {
	got := newEditionEntity(
		command.EditionCreateCommand{Title: "1st"},
		repo.EditionWorkRef{Id: 50, Uid: testWorkUid},
		nil, "e-uid", time.Now(),
	)
	if got.WorkUid != testWorkUid {
		t.Errorf("WorkUid = %q, want %q", got.WorkUid, testWorkUid)
	}
	if got.PublisherId != nil || got.PublisherUid != nil {
		t.Errorf("publisher = %v/%v, want nil/nil", got.PublisherId, got.PublisherUid)
	}
}

func TestEditionRepository_createWithRefs_persistsWorkAndPublisherUids(t *testing.T) {
	db := newEditionTableDB(t)
	r := &editionRepository{q: dao.Use(db)}

	id, err := r.Create(context.Background(),
		command.EditionCreateCommand{Title: "1st"},
		repo.EditionWorkRef{Id: 50, Uid: testWorkUid},
		&repo.EditionPublisherRef{Id: 7, Uid: testPublisherUid},
	)
	if err != nil {
		t.Fatalf("Create err: %v", err)
	}

	got := loadEdition(t, db, id)
	if got.WorkId != 50 || got.WorkUid != testWorkUid {
		t.Errorf("work = %d/%q, want 50/%q", got.WorkId, got.WorkUid, testWorkUid)
	}
	if got.PublisherId == nil || *got.PublisherId != 7 || got.PublisherUid == nil || *got.PublisherUid != testPublisherUid {
		t.Errorf("publisher = %v/%v, want 7/%q", got.PublisherId, got.PublisherUid, testPublisherUid)
	}
}

func TestEditionRepository_updateWithRefs_refreshesWorkAndPublisherUids(t *testing.T) {
	db := newEditionTableDB(t)
	r := &editionRepository{q: dao.Use(db)}
	id, err := r.Create(context.Background(),
		command.EditionCreateCommand{Title: "1st"},
		repo.EditionWorkRef{Id: 50, Uid: testWorkUid},
		nil,
	)
	if err != nil {
		t.Fatalf("Create err: %v", err)
	}

	err = r.Update(context.Background(), id, command.EditionUpdateCommand{},
		&repo.EditionWorkRef{Id: 60, Uid: testOtherWorkUid},
		&repo.EditionPublisherRef{Id: 7, Uid: testPublisherUid},
	)
	if err != nil {
		t.Fatalf("Update err: %v", err)
	}

	got := loadEdition(t, db, id)
	if got.WorkId != 60 || got.WorkUid != testOtherWorkUid {
		t.Errorf("work = %d/%q, want 60/%q", got.WorkId, got.WorkUid, testOtherWorkUid)
	}
	if got.PublisherId == nil || *got.PublisherId != 7 || got.PublisherUid == nil || *got.PublisherUid != testPublisherUid {
		t.Errorf("publisher = %v/%v, want 7/%q", got.PublisherId, got.PublisherUid, testPublisherUid)
	}
}

func TestEditionRepository_updateWithoutRefs_keepsWorkAndPublisherUids(t *testing.T) {
	db := newEditionTableDB(t)
	r := &editionRepository{q: dao.Use(db)}
	id, err := r.Create(context.Background(),
		command.EditionCreateCommand{Title: "1st"},
		repo.EditionWorkRef{Id: 50, Uid: testWorkUid},
		&repo.EditionPublisherRef{Id: 7, Uid: testPublisherUid},
	)
	if err != nil {
		t.Fatalf("Create err: %v", err)
	}
	title := "2nd"

	if err := r.Update(context.Background(), id, command.EditionUpdateCommand{Title: &title}, nil, nil); err != nil {
		t.Fatalf("Update err: %v", err)
	}

	got := loadEdition(t, db, id)
	if got.Title != "2nd" || got.WorkUid != testWorkUid || got.PublisherUid == nil || *got.PublisherUid != testPublisherUid {
		t.Errorf("edition = %+v", got)
	}
}
