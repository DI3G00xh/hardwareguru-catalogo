package service_test

import (
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/service"
)

func TestCatalogoService_ListarComponentes(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCatalogoService(repo)

	todos, err := svc.ListarComponentes("")
	if err != nil {
		t.Fatalf("error listando todos: %v", err)
	}
	if len(todos) == 0 {
		t.Errorf("la lista de todos no debería estar vacía")
	}

	cpus, err := svc.ListarComponentes("cpu")
	if err != nil {
		t.Fatalf("error listando cpus: %v", err)
	}
	for _, c := range cpus {
		if c.Tipo != domain.TipoCPU {
			t.Errorf("se esperaba tipo cpu, obtenido %s", c.Tipo)
		}
	}

	_, err = svc.ListarComponentes("tipo_invalido_xyz")
	if err != domain.ErrTipoInvalido {
		t.Errorf("se esperaba ErrTipoInvalido, obtenido %v", err)
	}
}

func TestCatalogoService_ObtenerSpecs(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewCatalogoService(repo)

	comp, err := svc.ObtenerSpecs("cpu-9600x")
	if err != nil {
		t.Fatalf("error obteniendo specs: %v", err)
	}
	if comp.Nombre != "Ryzen 5 9600X" {
		t.Errorf("esperaba Ryzen 5 9600X, obtenido %s", comp.Nombre)
	}

	_, err = svc.ObtenerSpecs("inexistente-123")
	if err != domain.ErrComponenteNoEncontrado {
		t.Errorf("esperaba ErrComponenteNoEncontrado, obtenido %v", err)
	}

	_, err = svc.ObtenerSpecs("")
	if err != domain.ErrComponenteNoEncontrado {
		t.Errorf("esperaba ErrComponenteNoEncontrado con id vacío, obtenido %v", err)
	}
}
