package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
)

func NewMigrationManager(dbCfg *config.DatabaseConfig) *MigrationManager {
	return &MigrationManager{
		dbCfg: dbCfg,
	}
}

type MigrationManager struct {
	dbCfg *config.DatabaseConfig
}

func (m MigrationManager) Migrate() {
	m.migrateDB(m.dbCfg.DataSource.Dsn)
	log.Printf("database migration applied successfully")
}

func (m MigrationManager) migrateDB(dsn string) {
	db, err := sql.Open("mysql", dsn+"&multiStatements=true")
	if err != nil {
		_ = fmt.Errorf("failed to open database: %v", err)
		panic(err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		_ = fmt.Errorf("failed to create driver: %v", err)
		panic(err)
	}
	inst, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		_ = fmt.Errorf("failed to create migration instance: %v", err)
		panic(err)
	}

	err = inst.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		_ = fmt.Errorf("failed to apply migrations: %v", err)
		panic(err)
	}
}
