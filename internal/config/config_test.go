package config_test

import (
	"os"
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/config"
)

func TestConfig_DefaultValues(t *testing.T) {
	// Limpiar variables de entorno para asegurar defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")

	cfg := config.Load()

	if cfg.Port != "50051" {
		t.Errorf("esperaba puerto por defecto 50051, obtenido %s", cfg.Port)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("esperaba db host localhost, obtenido %s", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("esperaba db port 5432, obtenido %s", cfg.DBPort)
	}
	if cfg.DBName != "hardwareguru_catalogo" {
		t.Errorf("esperaba db name hardwareguru_catalogo, obtenido %s", cfg.DBName)
	}
}

func TestConfig_CustomValuesAndDSN(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DB_HOST", "custom_host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "custom_user")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "custom_db")
	os.Setenv("DB_SSLMODE", "require")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
	}()

	cfg := config.Load()

	if cfg.Port != "9090" {
		t.Errorf("esperaba port 9090, obtenido %s", cfg.Port)
	}
	expectedDSN := "host=custom_host port=5433 user=custom_user password=secret dbname=custom_db sslmode=require"
	if cfg.DSN() != expectedDSN {
		t.Errorf("esperaba DSN '%s', obtenido '%s'", expectedDSN, cfg.DSN())
	}
}

