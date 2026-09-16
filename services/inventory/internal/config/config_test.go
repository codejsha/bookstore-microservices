package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyVaultTemporalCredentials_fileWithCredentials_setsAuthClient(t *testing.T) {
	path := filepath.Join(t.TempDir(), "temporal.properties")
	content := "temporal.auth.clientId=temporal-worker-inventory\ntemporal.auth.clientSecret=s3cret=x\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}
	cfg := &Config{Temporal: &TemporalConfig{Auth: &TemporalAuthConfig{Enabled: true, TokenURL: "http://keycloak"}}}

	applyVaultTemporalCredentials(cfg, path)

	auth := cfg.Temporal.Auth
	if auth.ClientID != "temporal-worker-inventory" || auth.ClientSecret != "s3cret=x" {
		t.Fatalf("auth = %+v, want client id and secret from file", auth)
	}
	if !auth.Enabled || auth.TokenURL != "http://keycloak" {
		t.Fatalf("auth = %+v, want enabled and token url preserved", auth)
	}
}

func TestApplyVaultTemporalCredentials_missingFile_leavesAuthUnset(t *testing.T) {
	cfg := &Config{Temporal: &TemporalConfig{Host: "localhost:7233"}}

	applyVaultTemporalCredentials(cfg, filepath.Join(t.TempDir(), "absent.properties"))

	if cfg.Temporal.Auth != nil {
		t.Fatalf("auth = %+v, want nil", cfg.Temporal.Auth)
	}
}
