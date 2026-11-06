package repository

import (
	"fmt"
	"time"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/config"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error al conectar con PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error al obtener instancia sql.DB: %w", err)
	}

	// Configuración del pool de conexiones
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migración del esquema según el modelo de dominio
	if err := db.AutoMigrate(&domain.Componente{}); err != nil {
		return nil, fmt.Errorf("error al ejecutar AutoMigrate para Componente: %w", err)
	}

	return db, nil
}

