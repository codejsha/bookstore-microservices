package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/logging"
)

const (
	readinessMaxOpenConns = 2
	readinessMaxIdleConns = 1
)

type ReadinessDataSource struct {
	db *sql.DB
}

func NewReadinessDataSource(
	dbCfg *config.DatabaseConfig,
	logHelper *logging.LogHelper,
) (*ReadinessDataSource, error) {
	dataSource, err := database.NewVaultAwareDataSource(dbCfg, logHelper)
	if err != nil {
		return nil, fmt.Errorf("readiness datasource: %w", err)
	}
	sqlDB, err := dataSource.DB().DB()
	if err != nil {
		return nil, fmt.Errorf("readiness datasource: %w", err)
	}
	applyReadinessPoolLimits(sqlDB)
	return &ReadinessDataSource{db: sqlDB}, nil
}

func applyReadinessPoolLimits(sqlDB *sql.DB) {
	sqlDB.SetMaxOpenConns(readinessMaxOpenConns)
	sqlDB.SetMaxIdleConns(readinessMaxIdleConns)
}

func (r *ReadinessDataSource) Ping(ctx context.Context) error {
	if r == nil || r.db == nil {
		return errors.New("readiness datasource is not configured")
	}
	return r.db.PingContext(ctx)
}

func (r *ReadinessDataSource) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}
