package service_test

import (
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/service"
)

// mockComponenteRepository simula el acceso a datos en memoria para pruebas de servicio
type mockComponenteRepository struct {
	componentes map[string]domain.Componente
	err         error
}

func newMockRepository() *mockComponenteRepository {
	socketAM5 := "AM5"
	socketLGA := "LGA1700"
	ramDDR5 := "DDR5"
	ramDDR4 := "DDR4"
	ramDDR3 := "DDR3"

	return &mockComponenteRepository{
		componentes: map[string]domain.Componente{
			"cpu-9600x": {
				ID:         "cpu-9600x",
				Tipo:       domain.TipoCPU,
				Nombre:     "Ryzen 5 9600X",
				Fabricante: "AMD",
				TDP:        65,
				Socket:     &socketAM5,
				Especificaciones: domain.JSONMap{
					"socket": "AM5",
					"tdp":    65,
				},
			},
			"cpu-14900k": {
				ID:         "cpu-14900k",
				Tipo:       domain.TipoCPU,
				Nombre:     "Core i9-14900K",
				Fabricante: "Intel",
				TDP:        253,
				Socket:     &socketLGA,
				Especificaciones: domain.JSONMap{
					"socket": "LGA1700",
					"tdp":    253,
				},
			},
			"mb-b650": {
				ID:         "mb-b650",
				Tipo:       domain.TipoPlaca,
				Nombre:     "B650 Gaming",
				Fabricante: "ASUS",
				Socket:     &socketAM5,
				TipoRAM:    &ramDDR5,
				Especificaciones: domain.JSONMap{
					"socket":   "AM5",
					"tipoRam":  "DDR5",
					"slotsRam": 4,
				},
			},
			"mb-b840": {
				ID:         "mb-b840",
				Tipo:       domain.TipoPlaca,
				Nombre:     "B840M Micro-ATX",
				Fabricante: "ASRock",
				Socket:     &socketAM5,
				TipoRAM:    &ramDDR5,
				Especificaciones: domain.JSONMap{
					"socket":   "AM5",
					"tipoRam":  "DDR5",
					"slotsRam": 2,
				},
			},
			"mb-b660-ddr4": {
				ID:         "mb-b660-ddr4",
				Tipo:       domain.TipoPlaca,
				Nombre:     "B660M DDR4",
				Fabricante: "ASUS",
				Socket:     &socketLGA,
				TipoRAM:    &ramDDR4,
				Especificaciones: domain.JSONMap{
					"socket":   "LGA1700",
					"tipoRam":  "DDR4",
					"slotsRam": 4,
				},
			},
			"mb-h81-ddr3": {
				ID:         "mb-h81-ddr3",
				Tipo:       domain.TipoPlaca,
				Nombre:     "H81M DDR3",
				Fabricante: "Gigabyte",
				Socket:     &socketLGA,
				TipoRAM:    &ramDDR3,
				Especificaciones: domain.JSONMap{
					"socket":   "LGA1700",
					"tipoRam":  "DDR3",
					"slotsRam": 2,
				},
			},
			"ram-ddr5-32": {
				ID:         "ram-ddr5-32",
				Tipo:       domain.TipoRAM,
				Nombre:     "Corsair DDR5 32GB",
				Fabricante: "Corsair",
				TipoRAM:    &ramDDR5,
				Especificaciones: domain.JSONMap{
					"tipoRam":     "DDR5",
					"capacidadGB": 32,
				},
			},
			"ram-ddr4-16": {
				ID:         "ram-ddr4-16",
				Tipo:       domain.TipoRAM,
				Nombre:     "Kingston DDR4 16GB",
				Fabricante: "Kingston",
				TipoRAM:    &ramDDR4,
				Especificaciones: domain.JSONMap{
					"tipoRam":     "DDR4",
					"capacidadGB": 16,
				},
			},
			"ram-ddr3-8": {
				ID:         "ram-ddr3-8",
				Tipo:       domain.TipoRAM,
				Nombre:     "Kingston DDR3 8GB",
				Fabricante: "Kingston",
				TipoRAM:    &ramDDR3,
				Especificaciones: domain.JSONMap{
					"tipoRam":     "DDR3",
					"capacidadGB": 8,
				},
			},
			"gpu-5060ti": {
				ID:         "gpu-5060ti",
				Tipo:       domain.TipoGPU,
				Nombre:     "RTX 5060 Ti",
				Fabricante: "NVIDIA",
				TDP:        180,
				Especificaciones: domain.JSONMap{
					"tdp": 180,
				},
			},
			"gpu-7800xt": {
				ID:         "gpu-7800xt",
				Tipo:       domain.TipoGPU,
				Nombre:     "RX 7800 XT",
				Fabricante: "AMD",
				TDP:        263,
				Especificaciones: domain.JSONMap{
					"tdp": 263,
				},
			},
			"psu-450": {
				ID:         "psu-450",
				Tipo:       domain.TipoFuente,
				Nombre:     "EVGA 450W Bronze",
				Fabricante: "EVGA",
				PotenciaW:  450,
			},
			"psu-500": {
				ID:         "psu-500",
				Tipo:       domain.TipoFuente,
				Nombre:     "EVGA 500W Bronze",
				Fabricante: "EVGA",
				PotenciaW:  500,
			},
			"psu-750": {
				ID:         "psu-750",
				Tipo:       domain.TipoFuente,
				Nombre:     "Corsair 750W Gold",
				Fabricante: "Corsair",
				PotenciaW:  750,
			},
		},
	}
}

func (m *mockComponenteRepository) Listar(tipo string) ([]domain.Componente, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []domain.Componente
	for _, c := range m.componentes {
		if tipo == "" || string(c.Tipo) == tipo {
			res = append(res, c)
		}
	}
	return res, nil
}

func (m *mockComponenteRepository) ObtenerPorID(id string) (*domain.Componente, error) {
	if m.err != nil {
		return nil, m.err
	}
	if c, ok := m.componentes[id]; ok {
		return &c, nil
	}
	return nil, domain.ErrComponenteNoEncontrado
}

func (m *mockComponenteRepository) ObtenerPorIDs(ids []string) ([]domain.Componente, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []domain.Componente
	for _, id := range ids {
		if c, ok := m.componentes[id]; ok {
			res = append(res, c)
		}
	}
	return res, nil
}

func TestCompatibilidadService_BuildCompletamenteCompatible(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	build := domain.BuildEspecificacion{
		CPUId:         "cpu-9600x",
		MotherboardId: "mb-b650",
		RAMId:         "ram-ddr5-32",
		GPUId:         "gpu-5060ti",
		PSUId:         "psu-750",
		RAMCantidad:   2,
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado en VerificarCompatibilidad: %v", err)
	}

	if !res.Compatible {
		t.Fatalf("la build debería ser compatible, errores: %v", res.Errores)
	}
	if len(res.Errores) != 0 {
		t.Errorf("no deberían haber errores, obtenido: %v", res.Errores)
	}

	if res.ConsumoEstimadoW != 354 {
		t.Errorf("esperaba consumo estimado 354W, obtenido %d", res.ConsumoEstimadoW)
	}
}

func TestCompatibilidadService_Regla1_SocketIncompatible(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	build := domain.BuildEspecificacion{
		CPUId:         "cpu-14900k",
		MotherboardId: "mb-b650",
		RAMId:         "ram-ddr5-32",
		GPUId:         "gpu-5060ti",
		PSUId:         "psu-750",
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Compatible {
		t.Errorf("la build debería ser incompatible debido al socket")
	}

	encontradoErrorSocket := false
	for _, e := range res.Errores {
		if testing.Verbose() {
			t.Logf("Error detectado: %s", e)
		}
		if containsString(e, "socket") {
			encontradoErrorSocket = true
		}
	}
	if !encontradoErrorSocket {
		t.Errorf("se esperaba error de socket en la lista: %v", res.Errores)
	}
}

func TestCompatibilidadService_Regla2_TecnologiaRAM(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	buildDDR4EnPlacaDDR5 := domain.BuildEspecificacion{
		CPUId:         "cpu-9600x",
		MotherboardId: "mb-b650",
		RAMId:         "ram-ddr4-16",
		GPUId:         "gpu-5060ti",
		PSUId:         "psu-750",
	}

	res, err := svc.VerificarCompatibilidad(buildDDR4EnPlacaDDR5)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Compatible {
		t.Errorf("debería ser incompatible por tecnologia RAM")
	}

	encontradoErrorRAM := false
	for _, e := range res.Errores {
		if containsString(e, "RAM") || containsString(e, "tecnología") {
			encontradoErrorRAM = true
		}
	}
	if !encontradoErrorRAM {
		t.Errorf("se esperaba error de tecnología de RAM, obtenido: %v", res.Errores)
	}

	buildDDR3Valida := domain.BuildEspecificacion{
		CPUId:         "cpu-14900k",
		MotherboardId: "mb-h81-ddr3",
		RAMId:         "ram-ddr3-8",
		PSUId:         "psu-750",
	}
	resDDR3, err := svc.VerificarCompatibilidad(buildDDR3Valida)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	for _, e := range resDDR3.Errores {
		if containsString(e, "RAM") {
			t.Errorf("no debería haber error de RAM para DDR3 en placa DDR3, obtenido: %s", e)
		}
	}
}

func TestCompatibilidadService_Regla3_RanurasRAMInsuficientes(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	buildRanurasExcedidas := domain.BuildEspecificacion{
		CPUId:         "cpu-9600x",
		MotherboardId: "mb-b840",
		RAMId:         "ram-ddr5-32",
		GPUId:         "gpu-5060ti",
		PSUId:         "psu-750",
		RAMCantidad:   4,
	}

	res, err := svc.VerificarCompatibilidad(buildRanurasExcedidas)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Compatible {
		t.Errorf("debería ser incompatible por ranuras de RAM insuficientes")
	}

	encontradoErrorRanuras := false
	for _, e := range res.Errores {
		if containsString(e, "ranuras") || containsString(e, "Ranuras") {
			encontradoErrorRanuras = true
		}
	}
	if !encontradoErrorRanuras {
		t.Errorf("se esperaba error de ranuras de RAM, obtenido: %v", res.Errores)
	}

	buildRanurasOK := buildRanurasExcedidas
	buildRanurasOK.RAMCantidad = 2
	resOK, err := svc.VerificarCompatibilidad(buildRanurasOK)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !resOK.Compatible {
		t.Errorf("debería ser compatible con 2 módulos en 2 slots, errores: %v", resOK.Errores)
	}
}

func TestCompatibilidadService_Regla4_ConsumoYFuenteInsuficiente(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	buildFuenteChica := domain.BuildEspecificacion{
		CPUId:         "cpu-14900k",
		MotherboardId: "mb-b660-ddr4",
		RAMId:         "ram-ddr4-16",
		GPUId:         "gpu-7800xt",
		PSUId:         "psu-500",
	}

	res, err := svc.VerificarCompatibilidad(buildFuenteChica)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if res.ConsumoEstimadoW != 680 {
		t.Errorf("esperaba consumo estimado 680W, obtenido %d", res.ConsumoEstimadoW)
	}

	if res.Compatible {
		t.Errorf("la build debería ser incompatible por fuente insuficiente")
	}

	encontradoErrorFuente := false
	for _, e := range res.Errores {
		if containsString(e, "fuente") || containsString(e, "Potencia") {
			encontradoErrorFuente = true
		}
	}
	if !encontradoErrorFuente {
		t.Errorf("se esperaba error de potencia de fuente, obtenido: %v", res.Errores)
	}

	buildFuenteOK := buildFuenteChica
	buildFuenteOK.PSUId = "psu-750"
	resOK, err := svc.VerificarCompatibilidad(buildFuenteOK)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !resOK.Compatible {
		t.Errorf("con fuente de 750W debería ser compatible, errores: %v", resOK.Errores)
	}
}

func TestCompatibilidadService_BuildSinGPU(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	buildSinGPU := domain.BuildEspecificacion{
		CPUId:         "cpu-9600x",
		MotherboardId: "mb-b650",
		RAMId:         "ram-ddr5-32",
		PSUId:         "psu-450",
		GPUId:         "",
	}

	res, err := svc.VerificarCompatibilidad(buildSinGPU)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if !res.Compatible {
		t.Fatalf("build de oficina sin GPU debería ser compatible, errores: %v", res.Errores)
	}
	if res.ConsumoEstimadoW != 138 {
		t.Errorf("esperaba 138W de consumo, obtenido %d", res.ConsumoEstimadoW)
	}
}

func TestCompatibilidadService_MultiplesErroresSimultaneos(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	buildMultiError := domain.BuildEspecificacion{
		CPUId:         "cpu-14900k",
		MotherboardId: "mb-b840",
		RAMId:         "ram-ddr4-16",
		GPUId:         "gpu-7800xt",
		PSUId:         "psu-500",
		RAMCantidad:   4,
	}

	res, err := svc.VerificarCompatibilidad(buildMultiError)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if res.Compatible {
		t.Fatalf("la build debería ser completamente incompatible")
	}

	if len(res.Errores) < 4 {
		t.Errorf("se esperaban al menos 4 errores específicos, obtenidos %d: %v", len(res.Errores), res.Errores)
	}
}

func TestCompatibilidadService_CamposObligatoriosFaltantes(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	res, err := svc.VerificarCompatibilidad(domain.BuildEspecificacion{})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Compatible {
		t.Errorf("solicitud vacía no debe ser compatible")
	}
	if len(res.Errores) != 4 {
		t.Errorf("se esperaban 4 errores por campos obligatorios faltantes, obtenido %d", len(res.Errores))
	}
}

func TestCompatibilidadService_ComponenteNoEncontrado(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCompatibilidadService(repo)

	build := domain.BuildEspecificacion{
		CPUId:         "cpu-fantasma",
		MotherboardId: "mb-b650",
		RAMId:         "ram-ddr5-32",
		PSUId:         "psu-750",
	}

	res, err := svc.VerificarCompatibilidad(build)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Compatible {
		t.Errorf("no debe ser compatible si un componente no existe")
	}
	if len(res.Errores) == 0 || !containsString(res.Errores[0], "no encontrado") {
		t.Errorf("se esperaba error de no encontrado, obtenido: %v", res.Errores)
	}
}

func containsString(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
