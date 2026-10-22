package config

import (
	"strings"
	"testing"
)

func setValidEnv(t *testing.T) {
	t.Helper()

	t.Setenv("JWT_SECRET", strings.Repeat("j", minJWTSecretLength))
	t.Setenv("DATABASE_URL", "postgres://postgres:secret@localhost:5432/kyokusu")
}

func TestLoadValidConfig(t *testing.T) {
	t.Chdir(t.TempDir())
	setValidEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.CookieSecure {
		t.Error("CookieSecure = true by default, want false")
	}
}

func TestLoadCookieSecureFromEnv(t *testing.T) {
	t.Chdir(t.TempDir())
	setValidEnv(t)
	t.Setenv("COOKIE_SECURE", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.CookieSecure {
		t.Error("CookieSecure = false, want true")
	}
}

func TestLoadRejectsInsecureDefaults(t *testing.T) {
	tests := []struct {
		name    string
		envKey  string
		envVal  string
		wantErr string
	}{
		{
			name:    "default jwt secret",
			envKey:  "JWT_SECRET",
			envVal:  defaultJWTSecret,
			wantErr: "JWT_SECRET",
		},
		{
			name:    "empty jwt secret",
			envKey:  "JWT_SECRET",
			envVal:  "",
			wantErr: "JWT_SECRET",
		},
		{
			name:    "short jwt secret",
			envKey:  "JWT_SECRET",
			envVal:  strings.Repeat("j", minJWTSecretLength-1),
			wantErr: "JWT_SECRET",
		},
		{
			name:    "default database url",
			envKey:  "DATABASE_URL",
			envVal:  defaultDatabaseURL,
			wantErr: "DATABASE_URL",
		},
		{
			name:    "empty database url",
			envKey:  "DATABASE_URL",
			envVal:  "",
			wantErr: "DATABASE_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			setValidEnv(t)
			t.Setenv(tt.envKey, tt.envVal)

			_, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want validation error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Load() error = %q, want it to mention %q", err, tt.wantErr)
			}
		})
	}
}
