package domain_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
)

func TestSeeds_ComponentsIntegrity(t *testing.T) {
	// Ruta relativa desde internal/domain hacia db/seeds/components.json
	seedPath := filepath.Join("..", "..", "db", "seeds", "components.json")
	data, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("no se pudo leer el archivo de semillas en %s: %v", seedPath, err)
	}

	var componentes []domain.Componente
	if err := json.Unmarshal(data, &componentes); err != nil {
		t.Fatalf("error deserializando components.json: %v", err)
	}

	if len(componentes) < 25 {
		t.Fatalf("se esperaban al menos 25 componentes, pero se encontraron %d", len(componentes))
	}

	ids := make(map[string]bool)
	tiposValidos := map[domain.TipoComponente]bool{
		domain.TipoCPU:    true,
		domain.TipoPlaca:  true,
		domain.TipoRAM:    true,
		domain.TipoGPU:    true,
		domain.TipoFuente: true,
	}

	for _, c := range componentes {
		if c.ID == "" {
			t.Errorf("componente sin ID: %+v", c)
		}
		if ids[c.ID] {
			t.Errorf("ID duplicado encontrado en semillas: %s", c.ID)
		}
		ids[c.ID] = true

		if !tiposValidos[c.Tipo] {
			t.Errorf("tipo inválido para componente %s: %s", c.ID, c.Tipo)
		}
		if c.Nombre == "" {
			t.Errorf("componente %s no tiene nombre", c.ID)
		}
		if c.Fabricante == "" {
			t.Errorf("componente %s no tiene fabricante", c.ID)
		}
		if c.Especificaciones == nil {
			t.Errorf("componente %s no tiene especificaciones", c.ID)
		}
	}

	// Verificar piezas clave 
	piezasClave := []string{
		"cpu-9600x",
		"cpu-14600k",
		"gpu-5060ti",
		"mb-b650",
		"ram-ddr5-16",
	}

	for _, id := range piezasClave {
		if !ids[id] {
			t.Errorf("pieza clave requerida no encontrada en semillas: %s", id)
		}
	}
}

