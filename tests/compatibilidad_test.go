package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/repository"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/service"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupSeededDB(t *testing.T) repository.ComponenteRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("error al abrir base de datos en memoria: %v", err)
	}

	if err := db.AutoMigrate(&domain.Componente{}); err != nil {
		t.Fatalf("error en AutoMigrate: %v", err)
	}

	seedPath := filepath.Join("..", "db", "seeds", "components.json")
	data, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("no se pudo leer archivo de semillas en %s: %v", seedPath, err)
	}

	var componentes []domain.Componente
	if err := json.Unmarshal(data, &componentes); err != nil {
		t.Fatalf("error al parsear JSON de semillas: %v", err)
	}

	for _, c := range componentes {
		if err := db.Create(&c).Error; err != nil {
			t.Fatalf("error al insertar semilla %s: %v", c.ID, err)
		}
	}

	return repository.NewComponenteRepository(db)
}

func TestCompatibilidad_AM5_DDR5_Exito(t *testing.T) {
	repo := setupSeededDB(t)
	svc := service.NewCompatibilidadService(repo)

	build := domain.BuildEspecificacion{
		CPUId:         "cpu-9600x",
		MotherboardId: "mb-b650",
		RAMId:         "ram-ddr5-32",
		GPUId:         "gpu-5060ti",
		PSUId:         "psu-650",
		RAMCantidad:   2,
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado en VerificarCompatibilidad: %v", err)
	}

	if !res.Compatible {
		t.Fatalf("la build AM5 debería ser completamente compatible, errores: %v", res.Errores)
	}

	if res.ConsumoEstimadoW != 354 {
		t.Errorf("esperaba consumo de 354W, obtenido %d", res.ConsumoEstimadoW)
	}
}

func TestCompatibilidad_Intel_LGA1700_Exito(t *testing.T) {
	repo := setupSeededDB(t)
	svc := service.NewCompatibilidadService(repo)

	build := domain.BuildEspecificacion{
		CPUId:         "cpu-14700k",
		MotherboardId: "mb-z790",
		RAMId:         "ram-ddr5-64",
		GPUId:         "gpu-5070",
		PSUId:         "psu-850",
		RAMCantidad:   2,
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado en VerificarCompatibilidad: %v", err)
	}

	if !res.Compatible {
		t.Fatalf("la build Intel debería ser completamente compatible, errores: %v", res.Errores)
	}

	if res.ConsumoEstimadoW != 510 {
		t.Errorf("esperaba consumo de 510W, obtenido %d", res.ConsumoEstimadoW)
	}
}

func TestCompatibilidad_ChoqueSockets(t *testing.T) {
	repo := setupSeededDB(t)
	svc := service.NewCompatibilidadService(repo)

	build := domain.BuildEspecificacion{
		CPUId:         "cpu-14600k",
		MotherboardId: "mb-b650",
		RAMId:         "ram-ddr5-32",
		GPUId:         "gpu-5060ti",
		PSUId:         "psu-650",
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if res.Compatible {
		t.Fatalf("la build no debe ser compatible debido a sockets incompatibles (LGA1700 vs AM5)")
	}

	encontrado := false
	for _, e := range res.Errores {
		if strings.Contains(strings.ToLower(e), "socket") {
			encontrado = true
		}
	}
	if !encontrado {
		t.Errorf("se esperaba error explícito de socket, obtenido: %v", res.Errores)
	}
}

func TestCompatibilidad_TecnologiaRAM_DDR4_vs_DDR5(t *testing.T) {
	repo := setupSeededDB(t)
	svc := service.NewCompatibilidadService(repo)

	build := domain.BuildEspecificacion{
		CPUId:         "cpu-14600k",
		MotherboardId: "mb-b660",
		RAMId:         "ram-ddr5-32",
		GPUId:         "gpu-5060ti",
		PSUId:         "psu-650",
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if res.Compatible {
		t.Fatalf("la build no debe ser compatible por tecnologías de RAM incompatibles (DDR4 vs DDR5)")
	}

	encontrado := false
	for _, e := range res.Errores {
		if strings.Contains(strings.ToLower(e), "ram") || strings.Contains(strings.ToLower(e), "tecnología") {
			encontrado = true
		}
	}
	if !encontrado {
		t.Errorf("se esperaba error de tecnología de RAM, obtenido: %v", res.Errores)
	}
}

func TestCompatibilidad_RanurasRAM_Insuficientes(t *testing.T) {
	repo := setupSeededDB(t)
	svc := service.NewCompatibilidadService(repo)

	// Placa ASRock B840M (AM5, DDR5, solo 2 slots de RAM) con 4 módulos solicitados
	build := domain.BuildEspecificacion{
		CPUId:         "cpu-9600x",
		MotherboardId: "mb-b840",
		RAMId:         "ram-ddr5-16",
		GPUId:         "gpu-4060",
		PSUId:         "psu-500",
		RAMCantidad:   4,
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if res.Compatible {
		t.Fatalf("la build no debe ser compatible al requerir más ranuras de las disponibles")
	}

	encontrado := false
	for _, e := range res.Errores {
		if strings.Contains(strings.ToLower(e), "ranuras") {
			encontrado = true
		}
	}
	if !encontrado {
		t.Errorf("se esperaba error de ranuras de RAM insuficientes, obtenido: %v", res.Errores)
	}
}

func TestCompatibilidad_FuenteInsuficiente(t *testing.T) {
	repo := setupSeededDB(t)
	svc := service.NewCompatibilidadService(repo)

	// CPU Core i9-14900K (TDP 253W) + RX 7800 XT (TDP 263W) + 50W base = 566W * 1.2 = 680W
	// Fuente EVGA 450W (450W < 680W -> Insuficiente)
	build := domain.BuildEspecificacion{
		CPUId:         "cpu-14900k",
		MotherboardId: "mb-z790",
		RAMId:         "ram-ddr5-32",
		GPUId:         "gpu-7800xt",
		PSUId:         "psu-450",
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if res.Compatible {
		t.Fatalf("la build no debe ser compatible con fuente de 450W para un consumo estimado de 680W")
	}

	if res.ConsumoEstimadoW != 680 {
		t.Errorf("esperaba consumo estimado de 680W, obtenido %d", res.ConsumoEstimadoW)
	}

	encontrado := false
	for _, e := range res.Errores {
		if strings.Contains(strings.ToLower(e), "fuente") || strings.Contains(strings.ToLower(e), "potencia") {
			encontrado = true
		}
	}
	if !encontrado {
		t.Errorf("se esperaba error de potencia de fuente insuficiente, obtenido: %v", res.Errores)
	}
}
