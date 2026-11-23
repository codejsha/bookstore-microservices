package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const testConfigYaml = `app:
  logging:
    level: info
server:
  port: 8080
  mode: release
database:
  datasource:
    driver: postgres
    host: db.local
    port: 5432
    db_name: test_db
    params:
      - sslmode=disable
      - connect_timeout=5
telemetry:
  collector:
    traceUrl: collector.local:4317
  resource:
    environment: prod
  sampling:
    ratio: 0.1
grpc:
  server:
    host: 0.0.0.0
    port: 9090
  reflection_enabled: true
keycloak:
  url: http://keycloak.local:8080
  realm: bookstore
  client_id: ${TEST_CFG_USERNAME}
  client_secret: ${TEST_CFG_SECRET}
`

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "application.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}
	return path
}

func TestConfigFilePath_envUnset_usesMountedDefault(t *testing.T) {
	t.Setenv(configFileEnv, "")
	if got := configFilePath(); got != "/app/config/application.yaml" {
		t.Errorf("configFilePath() = %q, want %q", got, "/app/config/application.yaml")
	}
}

func TestConfigFilePath_envSet_usesEnvPath(t *testing.T) {
	t.Setenv(configFileEnv, "/tmp/custom.yaml")
	if got := configFilePath(); got != "/tmp/custom.yaml" {
		t.Errorf("configFilePath() = %q, want %q", got, "/tmp/custom.yaml")
	}
}

func TestReadConfigFile_nestedYaml_decodesSections(t *testing.T) {
	t.Setenv("TEST_CFG_USERNAME", "search-user")
	t.Setenv("TEST_CFG_SECRET", "s3cret")
	cfg, err := readConfigFile(writeConfigFile(t, testConfigYaml))
	if err != nil {
		t.Fatalf("readConfigFile: %v", err)
	}
	if cfg.Server.Port != "8080" || cfg.Server.Mode != "release" {
		t.Errorf("Server = %+v, want port 8080 mode release", cfg.Server)
	}
	ds := cfg.Database.DataSource
	if ds.Host != "db.local" || ds.Port != 5432 || ds.DbName != "test_db" {
		t.Errorf("DataSource = %+v, want db.local:5432/test_db", ds)
	}
	if len(ds.Params) != 2 || ds.Params[0] != "sslmode=disable" || ds.Params[1] != "connect_timeout=5" {
		t.Errorf("DataSource.Params = %v, want [sslmode=disable connect_timeout=5]", ds.Params)
	}
	if cfg.Telemetry.Collector.TraceUrl != "collector.local:4317" {
		t.Errorf("Telemetry.Collector.TraceUrl = %q, want %q", cfg.Telemetry.Collector.TraceUrl, "collector.local:4317")
	}
	if cfg.Telemetry.Sampling.Ratio != 0.1 {
		t.Errorf("Telemetry.Sampling.Ratio = %v, want 0.1", cfg.Telemetry.Sampling.Ratio)
	}
	if cfg.App.Logging.Level != "info" {
		t.Errorf("App.Logging.Level = %q, want %q", cfg.App.Logging.Level, "info")
	}
	if cfg.Grpc.Server == nil || cfg.Grpc.Server.Port != "9090" || !cfg.Grpc.ReflectionEnabled {
		t.Errorf("Grpc = %+v, want server port 9090 with reflection", cfg.Grpc)
	}
}

func TestReadConfigFile_envPlaceholders_expandsFromEnvironment(t *testing.T) {
	secret := "p@ss: #w$rd 'x\""
	t.Setenv("TEST_CFG_USERNAME", "search-user")
	t.Setenv("TEST_CFG_SECRET", secret)
	cfg, err := readConfigFile(writeConfigFile(t, testConfigYaml))
	if err != nil {
		t.Fatalf("readConfigFile: %v", err)
	}
	if cfg.Keycloak.ClientId != "search-user" {
		t.Errorf("Keycloak.ClientId = %q, want %q", cfg.Keycloak.ClientId, "search-user")
	}
	if cfg.Keycloak.ClientSecret != secret {
		t.Errorf("Keycloak.ClientSecret = %q, want %q", cfg.Keycloak.ClientSecret, secret)
	}
}

func TestReadConfigFile_absentFile_failsWithNotExist(t *testing.T) {
	_, err := readConfigFile(filepath.Join(t.TempDir(), "missing.yaml"))
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("readConfigFile error = %v, want fs.ErrNotExist", err)
	}
}

func TestReadConfigFile_scalarForList_fails(t *testing.T) {
	content := "app:\n  logging:\n    level: info\ndatabase:\n  datasource:\n    params: sslmode=disable\n"
	if _, err := readConfigFile(writeConfigFile(t, content)); err == nil {
		t.Error("readConfigFile error = nil, want scalar-for-list rejection")
	}
}

func TestReadConfigFile_missingAppSection_fails(t *testing.T) {
	if _, err := readConfigFile(writeConfigFile(t, "server:\n  port: 8080\n")); err == nil {
		t.Error("readConfigFile error = nil, want missing app section error")
	}
}

func TestReadConfigFile_malformedYaml_fails(t *testing.T) {
	if _, err := readConfigFile(writeConfigFile(t, "app: [unclosed\n")); err == nil {
		t.Error("readConfigFile error = nil, want parse error")
	}
}
