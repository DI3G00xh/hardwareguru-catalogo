package repository_test

import (
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("error al abrir base de datos en memoria: %v", err)
	}

	if err := db.AutoMigrate(&domain.Componente{}); err != nil {
		t.Fatalf("error en AutoMigrate de prueba: %v", err)
	}

	socketAM5 := "AM5"
	socketLGA := "LGA1700"
	ramDDR5 := "DDR5"
	ramDDR4 := "DDR4"

	componentes := []domain.Componente{
		{
			ID:         "cpu-9600x",
			Tipo:       domain.TipoCPU,
			Nombre:     "Ryzen 5 9600X",
			Fabricante: "AMD",
			TDP:        65,
			Socket:     &socketAM5,
			Origen:     "seed",
			Especificaciones: domain.JSONMap{
				"socket": "AM5",
				"tdp":    65,
			},
		},
		{
			ID:         "cpu-14600k",
			Tipo:       domain.TipoCPU,
			Nombre:     "Core i5-14600K",
			Fabricante: "Intel",
			TDP:        125,
			Socket:     &socketLGA,
			Origen:     "seed",
			Especificaciones: domain.JSONMap{
				"socket": "LGA1700",
				"tdp":    125,
			},
		},
		{
			ID:         "mb-b650",
			Tipo:       domain.TipoPlaca,
			Nombre:     "B650 Gaming",
			Fabricante: "ASUS",
			Socket:     &socketAM5,
			TipoRAM:    &ramDDR5,
			Origen:     "seed",
			Especificaciones: domain.JSONMap{
				"socket":   "AM5",
				"tipoRam":  "DDR5",
				"slotsRam": 4,
			},
		},
		{
			ID:         "ram-ddr5-32",
			Tipo:       domain.TipoRAM,
			Nombre:     "DDR5 32 GB",
			Fabricante: "Corsair",
			TipoRAM:    &ramDDR5,
			Origen:     "seed",
			Especificaciones: domain.JSONMap{
				"tipoRam":     "DDR5",
				"capacidadGB": 32,
			},
		},
		{
			ID:         "ram-ddr4-16",
			Tipo:       domain.TipoRAM,
			Nombre:     "DDR4 16 GB",
			Fabricante: "Kingston",
			TipoRAM:    &ramDDR4,
			Origen:     "seed",
		},
		{
			ID:         "psu-650",
			Tipo:       domain.TipoFuente,
			Nombre:     "650W Gold",
			Fabricante: "Corsair",
			PotenciaW:  650,
			Origen:     "seed",
		},
	}

	for _, c := range componentes {
		if err := db.Create(&c).Error; err != nil {
			t.Fatalf("error insertando componente de prueba %s: %v", c.ID, err)
		}
	}

	return db
}

func TestComponenteRepository_Listar(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewComponenteRepository(db)

	todos, err := repo.Listar("")
	if err != nil {
		t.Fatalf("error listando todos: %v", err)
	}
	if len(todos) != 6 {
		t.Errorf("esperaba 6 componentes, obtenido %d", len(todos))
	}

	cpus, err := repo.Listar("cpu")
	if err != nil {
		t.Fatalf("error listando cpus: %v", err)
	}
	if len(cpus) != 2 {
		t.Errorf("esperaba 2 CPUs, obtenido %d", len(cpus))
	}
	for _, c := range cpus {
		if c.Tipo != domain.TipoCPU {
			t.Errorf("componente no es CPU: %v", c)
		}
	}

	inexistentes, err := repo.Listar("refrigeracion")
	if err != nil {
		t.Fatalf("error listando tipo inexistente: %v", err)
	}
	if len(inexistentes) != 0 {
		t.Errorf("esperaba 0 componentes para tipo inexistente, obtenido %d", len(inexistentes))
	}
}

func TestComponenteRepository_ObtenerPorID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewComponenteRepository(db)

	c, err := repo.ObtenerPorID("cpu-9600x")
	if err != nil {
		t.Fatalf("error obteniendo componente existente: %v", err)
	}
	if c.Nombre != "Ryzen 5 9600X" {
		t.Errorf("esperaba Ryzen 5 9600X, obtenido %s", c.Nombre)
	}

	_, err = repo.ObtenerPorID("cpu-inexistente")
	if err != domain.ErrComponenteNoEncontrado {
		t.Errorf("esperaba ErrComponenteNoEncontrado, obtenido %v", err)
	}

	_, err = repo.ObtenerPorID("   ")
	if err != domain.ErrComponenteNoEncontrado {
		t.Errorf("esperaba ErrComponenteNoEncontrado con id vacío, obtenido %v", err)
	}
}

func TestComponenteRepository_ObtenerPorIDs(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewComponenteRepository(db)

	vacio, err := repo.ObtenerPorIDs([]string{})
	if err != nil {
		t.Fatalf("error con lista vacía: %v", err)
	}
	if len(vacio) != 0 {
		t.Errorf("esperaba slice vacío, obtenido %d", len(vacio))
	}

	ids := []string{"cpu-9600x", "mb-b650", "ram-ddr5-32"}
	encontrados, err := repo.ObtenerPorIDs(ids)
	if err != nil {
		t.Fatalf("error obteniendo por IDs: %v", err)
	}
	if len(encontrados) != 3 {
		t.Errorf("esperaba 3 componentes, obtenido %d", len(encontrados))
	}
}
