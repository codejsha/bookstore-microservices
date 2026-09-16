package support

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestSQLiteDB(t *testing.T, name string) *sql.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db handle: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return sqlDB
}

func TestApplyReadinessPoolLimits_freshPool_capsFarBelowAppPool(t *testing.T) {
	sqlDB := sql.OpenDB(unreachableConnector{})
	t.Cleanup(func() { _ = sqlDB.Close() })

	applyReadinessPoolLimits(sqlDB)

	if got := sqlDB.Stats().MaxOpenConnections; got != readinessMaxOpenConns {
		t.Fatalf("MaxOpenConnections = %d, want %d", got, readinessMaxOpenConns)
	}
	if readinessMaxOpenConns >= maxOpenConns {
		t.Fatalf("readiness pool cap %d must stay below app pool cap %d", readinessMaxOpenConns, maxOpenConns)
	}
}

func TestReadinessDataSource_openPool_pingSucceeds(t *testing.T) {
	r := &ReadinessDataSource{db: newTestSQLiteDB(t, t.Name())}

	if err := r.Ping(context.Background()); err != nil {
		t.Fatalf("Ping() = %v, want nil", err)
	}
}

func TestReadinessDataSource_closedPool_pingFails(t *testing.T) {
	r := &ReadinessDataSource{db: newTestSQLiteDB(t, t.Name())}
	if err := r.Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}

	if err := r.Ping(context.Background()); err == nil {
		t.Fatal("Ping() = nil, want error from closed pool")
	}
}

func TestReadinessDataSource_unconfiguredPool_pingFailsAndCloseIsNoOp(t *testing.T) {
	var r *ReadinessDataSource

	if err := r.Ping(context.Background()); err == nil {
		t.Fatal("Ping() = nil, want error")
	}
	if err := r.Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}
}

func TestGinServer_readyRoute_readinessPoolHealthy_returns200(t *testing.T) {
	s := newTestGinServer(t)
	s.readiness = &ReadinessDataSource{db: newTestSQLiteDB(t, t.Name())}

	if code := getHealthPath(s, "/health/ready"); code != http.StatusOK {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusOK)
	}
}

func TestGinServer_readyRoute_readinessPoolClosed_returns503(t *testing.T) {
	s := newTestGinServer(t)
	s.readiness = &ReadinessDataSource{db: newTestSQLiteDB(t, t.Name())}
	if err := s.readiness.Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}

	if code := getHealthPath(s, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusServiceUnavailable)
	}
}

func TestGinServer_readyRoute_drainingWithHealthyReadinessPool_returns503(t *testing.T) {
	s := newTestGinServer(t)
	s.readiness = &ReadinessDataSource{db: newTestSQLiteDB(t, t.Name())}

	s.BeginDrain()

	if code := getHealthPath(s, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusServiceUnavailable)
	}
}

func TestGinServer_readyRoute_appPoolFullyHeld_returns200(t *testing.T) {
	appDB := newTestSQLiteDB(t, "app-"+t.Name())
	applyConnectionPoolLimits(appDB)
	held := make([]*sql.Conn, 0, maxOpenConns)
	for range maxOpenConns {
		conn, err := appDB.Conn(context.Background())
		if err != nil {
			t.Fatalf("hold app connection: %v", err)
		}
		held = append(held, conn)
	}
	t.Cleanup(func() {
		for _, conn := range held {
			_ = conn.Close()
		}
	})

	saturated, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := appDB.PingContext(saturated); err == nil {
		t.Fatal("app pool served another connection, want it saturated")
	}

	s := newTestGinServer(t)
	s.readiness = &ReadinessDataSource{db: newTestSQLiteDB(t, "readiness-"+t.Name())}

	if code := getHealthPath(s, "/health/ready"); code != http.StatusOK {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusOK)
	}
}
