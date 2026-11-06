package domain_test

import (
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
)

func TestJSONMap_Value(t *testing.T) {
	// Nil map debe devolver "{}" serializado
	var nilMap domain.JSONMap
	val, err := nilMap.Value()
	if err != nil {
		t.Fatalf("error inesperado con nil map: %v", err)
	}
	if string(val.([]byte)) != "{}" {
		t.Errorf("esperaba '{}', obtenido %s", string(val.([]byte)))
	}

	// Mapa con datos
	m := domain.JSONMap{"socket": "AM5", "cores": 8}
	val, err = m.Value()
	if err != nil {
		t.Fatalf("error serializando map: %v", err)
	}
	if len(val.([]byte)) == 0 {
		t.Errorf("el JSON serializado no debería estar vacío")
	}
}

func TestJSONMap_Scan(t *testing.T) {
	var m domain.JSONMap

	// Escaneo de nil
	err := m.Scan(nil)
	if err != nil {
		t.Fatalf("error escaneando nil: %v", err)
	}
	if m == nil {
		t.Fatal("el mapa no debería ser nil tras escanear nil")
	}

	// Escaneo de JSON válido
	rawJSON := []byte(`{"socket":"LGA1700","watts":125}`)
	err = m.Scan(rawJSON)
	if err != nil {
		t.Fatalf("error escaneando bytes válidos: %v", err)
	}
	if m["socket"] != "LGA1700" {
		t.Errorf("esperaba socket LGA1700, obtenido %v", m["socket"])
	}

	// Escaneo de tipo inválido
	err = m.Scan("invalid-string-not-bytes")
	if err == nil {
		t.Errorf("se esperaba error al escanear un tipo no []byte")
	}
}

func TestComponente_Creation(t *testing.T) {
	socket := "AM5"
	tipoRAM := "DDR5"
	comp := domain.Componente{
		ID:         "cpu-ryzen-5-7600",
		Tipo:       domain.TipoCPU,
		Nombre:     "AMD Ryzen 5 7600",
		Fabricante: "AMD",
		TDP:        65,
		Socket:     &socket,
		TipoRAM:    &tipoRAM,
		PotenciaW:  0,
		Origen:     "seed",
		Especificaciones: domain.JSONMap{
			"cores":   6,
			"threads": 12,
		},
	}

	if comp.ID != "cpu-ryzen-5-7600" {
		t.Errorf("esperaba ID cpu-ryzen-5-7600, obtenido %s", comp.ID)
	}
	if comp.Tipo != domain.TipoCPU {
		t.Errorf("esperaba tipo cpu, obtenido %s", comp.Tipo)
	}
	if *comp.Socket != "AM5" {
		t.Errorf("esperaba socket AM5, obtenido %s", *comp.Socket)
	}
}

func TestComponente_Helpers(t *testing.T) {
	socket := " am5 "
	tipoRAM := " ddr5 "
	comp := domain.Componente{
		ID:        "mb-b650",
		Tipo:      domain.TipoPlaca,
		Socket:    &socket,
		TipoRAM:   &tipoRAM,
		TDP:       0,
		PotenciaW: 0,
		Especificaciones: domain.JSONMap{
			"slotsRam": 4,
		},
	}

	if s := comp.ObtenerSocket(); s != "AM5" {
		t.Errorf("esperaba socket 'AM5', obtenido '%s'", s)
	}
	if ram := comp.ObtenerTipoRAM(); ram != "DDR5" {
		t.Errorf("esperaba tipo RAM 'DDR5', obtenido '%s'", ram)
	}
	if slots := comp.ObtenerSlotsRAM(); slots != 4 {
		t.Errorf("esperaba 4 slots de RAM, obtenido %d", slots)
	}

	// Caso de extracción desde Especificaciones JSON
	compFromSpecs := domain.Componente{
		ID:   "mb-legacy",
		Tipo: domain.TipoPlaca,
		Especificaciones: domain.JSONMap{
			"socket":   "lga1700",
			"tipoRam":  "ddr4",
			"slotsRam": "2",
			"tdp":      125,
		},
	}
	if s := compFromSpecs.ObtenerSocket(); s != "LGA1700" {
		t.Errorf("esperaba socket 'LGA1700', obtenido '%s'", s)
	}
	if ram := compFromSpecs.ObtenerTipoRAM(); ram != "DDR4" {
		t.Errorf("esperaba 'DDR4', obtenido '%s'", ram)
	}
	if slots := compFromSpecs.ObtenerSlotsRAM(); slots != 2 {
		t.Errorf("esperaba 2 slots, obtenido %d", slots)
	}
	if tdp := compFromSpecs.ObtenerTDP(); tdp != 125 {
		t.Errorf("esperaba 125 TDP, obtenido %d", tdp)
	}

	psu := domain.Componente{
		ID:        "psu-750",
		Tipo:      domain.TipoFuente,
		PotenciaW: 750,
	}
	if psu.ObtenerPotenciaW() != 750 {
		t.Errorf("esperaba 750W, obtenido %d", psu.ObtenerPotenciaW())
	}

	// Pruebas de tipos en Especificaciones (int32, int64, float64, string)
	compTipos := domain.Componente{
		ID:   "comp-tipos",
		Tipo: domain.TipoPlaca,
		Especificaciones: domain.JSONMap{
			"socket":   "am5",
			"tipo_ram": "ddr5",
			"slotsRam": int64(4),
			"modulos":  int32(2),
			"tdp":      float64(105),
			"potencia": "850",
		},
	}
	if compTipos.ObtenerSocket() != "AM5" {
		t.Errorf("esperaba AM5, obtenido %s", compTipos.ObtenerSocket())
	}
	if compTipos.ObtenerTipoRAM() != "DDR5" {
		t.Errorf("esperaba DDR5, obtenido %s", compTipos.ObtenerTipoRAM())
	}
	if compTipos.ObtenerSlotsRAM() != 4 {
		t.Errorf("esperaba 4 slots, obtenido %d", compTipos.ObtenerSlotsRAM())
	}
	if compTipos.ObtenerModulosRAM() != 2 {
		t.Errorf("esperaba 2 modulos, obtenido %d", compTipos.ObtenerModulosRAM())
	}
	if compTipos.ObtenerTDP() != 105 {
		t.Errorf("esperaba 105 TDP, obtenido %d", compTipos.ObtenerTDP())
	}
	if compTipos.ObtenerPotenciaW() != 850 {
		t.Errorf("esperaba 850W, obtenido %d", compTipos.ObtenerPotenciaW())
	}

	// Fallbacks con valores por defecto
	compVacio := domain.Componente{
		ID:               "comp-vacio",
		Especificaciones: domain.JSONMap{},
	}
	if compVacio.ObtenerSocket() != "" {
		t.Errorf("esperaba socket vacío")
	}
	if compVacio.ObtenerTipoRAM() != "" {
		t.Errorf("esperaba tipo RAM vacío")
	}
	if compVacio.ObtenerSlotsRAM() != 0 {
		t.Errorf("esperaba 0 slots")
	}
	if compVacio.ObtenerModulosRAM() != 1 {
		t.Errorf("esperaba 1 módulo por defecto")
	}
	if compVacio.ObtenerTDP() != 0 {
		t.Errorf("esperaba 0 TDP")
	}
	if compVacio.ObtenerPotenciaW() != 0 {
		t.Errorf("esperaba 0 PotenciaW")
	}
}
