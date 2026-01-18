package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any existing env vars
	os.Unsetenv("TEXT_TO_SQL_PROXY_PORT")
	os.Unsetenv("TEXT_TO_SQL_PROXY_ALLOWED_ORIGIN")
	os.Unsetenv("TEXT_TO_SQL_PROXY_PROVIDER")
	os.Unsetenv("TEXT_TO_SQL_PROXY_DATABASE")
	os.Unsetenv("TEXT_TO_SQL_PROXY_TLS_CERT")
	os.Unsetenv("TEXT_TO_SQL_PROXY_TLS_KEY")

	cfg := Load()

	if cfg.Port != 4000 {
		t.Errorf("expected default port 4000, got %d", cfg.Port)
	}
	if cfg.AllowedOrigin != "https://sql-workbench.com" {
		t.Errorf("expected default origin https://sql-workbench.com, got %s", cfg.AllowedOrigin)
	}
	if cfg.Provider != "claude" {
		t.Errorf("expected default provider claude, got %s", cfg.Provider)
	}
	if cfg.Database != "DuckDB" {
		t.Errorf("expected default database DuckDB, got %s", cfg.Database)
	}
	if cfg.TLSCert != "" {
		t.Errorf("expected empty TLS cert, got %s", cfg.TLSCert)
	}
	if cfg.TLSKey != "" {
		t.Errorf("expected empty TLS key, got %s", cfg.TLSKey)
	}
	if cfg.TLSEnabled() {
		t.Error("expected TLS to be disabled by default")
	}
}

func TestLoad_CustomPort(t *testing.T) {
	os.Setenv("TEXT_TO_SQL_PROXY_PORT", "8080")
	defer os.Unsetenv("TEXT_TO_SQL_PROXY_PORT")

	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	tests := []struct {
		name     string
		portVal  string
		expected int
	}{
		{"non-numeric", "abc", 4000},
		{"zero", "0", 4000},
		{"negative", "-1", 4000},
		{"too high", "65536", 4000},
		{"empty", "", 4000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.portVal == "" {
				os.Unsetenv("TEXT_TO_SQL_PROXY_PORT")
			} else {
				os.Setenv("TEXT_TO_SQL_PROXY_PORT", tc.portVal)
				defer os.Unsetenv("TEXT_TO_SQL_PROXY_PORT")
			}

			cfg := Load()

			if cfg.Port != tc.expected {
				t.Errorf("expected port %d for %q, got %d", tc.expected, tc.portVal, cfg.Port)
			}
		})
	}
}

func TestLoad_CustomAllowedOrigin(t *testing.T) {
	os.Setenv("TEXT_TO_SQL_PROXY_ALLOWED_ORIGIN", "https://example.com")
	defer os.Unsetenv("TEXT_TO_SQL_PROXY_ALLOWED_ORIGIN")

	cfg := Load()

	if cfg.AllowedOrigin != "https://example.com" {
		t.Errorf("expected origin https://example.com, got %s", cfg.AllowedOrigin)
	}
}

func TestLoad_CustomProvider(t *testing.T) {
	os.Setenv("TEXT_TO_SQL_PROXY_PROVIDER", "gemini")
	defer os.Unsetenv("TEXT_TO_SQL_PROXY_PROVIDER")

	cfg := Load()

	if cfg.Provider != "gemini" {
		t.Errorf("expected provider gemini, got %s", cfg.Provider)
	}
}

func TestLoad_CustomDatabase(t *testing.T) {
	os.Setenv("TEXT_TO_SQL_PROXY_DATABASE", "PostgreSQL")
	defer os.Unsetenv("TEXT_TO_SQL_PROXY_DATABASE")

	cfg := Load()

	if cfg.Database != "PostgreSQL" {
		t.Errorf("expected database PostgreSQL, got %s", cfg.Database)
	}
}

func TestLoad_AllCustomValues(t *testing.T) {
	os.Setenv("TEXT_TO_SQL_PROXY_PORT", "3000")
	os.Setenv("TEXT_TO_SQL_PROXY_ALLOWED_ORIGIN", "https://myapp.com")
	os.Setenv("TEXT_TO_SQL_PROXY_PROVIDER", "codex")
	os.Setenv("TEXT_TO_SQL_PROXY_DATABASE", "MySQL")
	defer func() {
		os.Unsetenv("TEXT_TO_SQL_PROXY_PORT")
		os.Unsetenv("TEXT_TO_SQL_PROXY_ALLOWED_ORIGIN")
		os.Unsetenv("TEXT_TO_SQL_PROXY_PROVIDER")
		os.Unsetenv("TEXT_TO_SQL_PROXY_DATABASE")
	}()

	cfg := Load()

	if cfg.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Port)
	}
	if cfg.AllowedOrigin != "https://myapp.com" {
		t.Errorf("expected origin https://myapp.com, got %s", cfg.AllowedOrigin)
	}
	if cfg.Provider != "codex" {
		t.Errorf("expected provider codex, got %s", cfg.Provider)
	}
	if cfg.Database != "MySQL" {
		t.Errorf("expected database MySQL, got %s", cfg.Database)
	}
}

func TestLoad_TLSConfig(t *testing.T) {
	os.Setenv("TEXT_TO_SQL_PROXY_TLS_CERT", "/path/to/cert.pem")
	os.Setenv("TEXT_TO_SQL_PROXY_TLS_KEY", "/path/to/key.pem")
	defer func() {
		os.Unsetenv("TEXT_TO_SQL_PROXY_TLS_CERT")
		os.Unsetenv("TEXT_TO_SQL_PROXY_TLS_KEY")
	}()

	cfg := Load()

	if cfg.TLSCert != "/path/to/cert.pem" {
		t.Errorf("expected TLS cert /path/to/cert.pem, got %s", cfg.TLSCert)
	}
	if cfg.TLSKey != "/path/to/key.pem" {
		t.Errorf("expected TLS key /path/to/key.pem, got %s", cfg.TLSKey)
	}
	if !cfg.TLSEnabled() {
		t.Error("expected TLS to be enabled when both cert and key are set")
	}
}

func TestTLSEnabled_OnlyCert(t *testing.T) {
	cfg := Config{TLSCert: "/path/to/cert.pem", TLSKey: ""}
	if cfg.TLSEnabled() {
		t.Error("TLS should not be enabled with only cert")
	}
}

func TestTLSEnabled_OnlyKey(t *testing.T) {
	cfg := Config{TLSCert: "", TLSKey: "/path/to/key.pem"}
	if cfg.TLSEnabled() {
		t.Error("TLS should not be enabled with only key")
	}
}

func TestTLSEnabled_BothSet(t *testing.T) {
	cfg := Config{TLSCert: "/path/to/cert.pem", TLSKey: "/path/to/key.pem"}
	if !cfg.TLSEnabled() {
		t.Error("TLS should be enabled when both cert and key are set")
	}
}
