package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		wantPort string
		wantErr  bool
	}{
		{
			name: "success",
			env: map[string]string{
				"PORT":        "9000",
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "notes_app",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "password",
				"JWT_SECRET":  "secret",
			},
			wantPort: "9000",
		},
		{
			name: "default port",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "notes_app",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "password",
				"JWT_SECRET":  "secret",
			},
			wantPort: "8080",
		},
		{
			name: "missing db host",
			env: map[string]string{
				"DB_PORT":     "5432",
				"DB_NAME":     "notes_app",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "password",
				"JWT_SECRET":  "secret",
			},
			wantErr: true,
		},
		{
			name: "missing db password",
			env: map[string]string{
				"DB_HOST":    "localhost",
				"DB_PORT":    "5432",
				"DB_NAME":    "notes_app",
				"DB_USER":    "postgres",
				"JWT_SECRET": "secret",
			},
			wantErr: true,
		},
		{
			name: "missing jwt secret",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "notes_app",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "password",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.env["PORT"])
			t.Setenv("DB_HOST", tt.env["DB_HOST"])
			t.Setenv("DB_PORT", tt.env["DB_PORT"])
			t.Setenv("DB_NAME", tt.env["DB_NAME"])
			t.Setenv("DB_USER", tt.env["DB_USER"])
			t.Setenv("DB_PASSWORD", tt.env["DB_PASSWORD"])
			t.Setenv("JWT_SECRET", tt.env["JWT_SECRET"])

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Port != tt.wantPort {
				t.Fatalf(
					"expected port %q, got %q",
					tt.wantPort,
					cfg.Port,
				)
			}
			if tt.name == "success" {
				if cfg.DBHost != "localhost" {
					t.Fatalf("expected DBHost %q, got %q", "localhost", cfg.DBHost)
				}

				if cfg.DBPort != "5432" {
					t.Fatalf("expected DBPort %q, got %q", "5432", cfg.DBPort)
				}

				if cfg.DBName != "notes_app" {
					t.Fatalf("expected DBName %q, got %q", "notes_app", cfg.DBName)
				}

				if cfg.DBUser != "postgres" {
					t.Fatalf("expected DBUser %q, got %q", "postgres", cfg.DBUser)
				}

				if cfg.DBPassword != "password" {
					t.Fatalf("expected DBPassword %q, got %q", "password", cfg.DBPassword)
				}

				if cfg.JWTSecret != "secret" {
					t.Fatalf("expected JWTSecret %q, got %q", "secret", cfg.JWTSecret)
				}
			}
		})
	}
}
