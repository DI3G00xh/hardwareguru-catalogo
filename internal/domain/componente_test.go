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

