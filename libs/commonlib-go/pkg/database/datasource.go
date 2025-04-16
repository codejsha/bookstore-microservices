package database

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
)

func NewDataSource(
	dbCfg *config.DatabaseConfig,
) *DataSource {
	gormCfg := createConfig(dbCfg)
	db := setupDB(dbCfg, gormCfg)

	return &DataSource{
		dbCfg: dbCfg,
		db:    db,
	}
}

type DataSource struct {
	dbCfg *config.DatabaseConfig
	db    *gorm.DB
}

func (s *DataSource) DB() *gorm.DB {
	return s.db
}

func createConfig(dbCfg *config.DatabaseConfig) *gorm.Config {
	gormLoggerCfg := logger.Config{
		SlowThreshold: time.Millisecond,
		LogLevel:      setLoggingLevel(dbCfg.Gorm.Logging.Level),
	}
	gormWriter := logrus.NewWriter()
	gormLogger := logger.New(gormWriter, gormLoggerCfg)
	cfg := &gorm.Config{Logger: gormLogger}
	return cfg
}

func setupDB(dbCfg *config.DatabaseConfig, gormCfg *gorm.Config) *gorm.DB {
	dialect := mysql.Open(dbCfg.DataSource.Dsn)
	db, err := gorm.Open(dialect, gormCfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.Use(tracing.NewPlugin())
	if err != nil {
		log.Fatalf("failed to setup tracing: %v", err)
	}

	log.Printf("connected to database successfully")
	return db
}

func setLoggingLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	default:
		return logger.Info
	}
}
