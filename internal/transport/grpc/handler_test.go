package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	grpcTransport "github.com/DI3G00xh/hardwareguru-catalogo/internal/transport/grpc"
	hardwarepb "github.com/DI3G00xh/hardwareguru-catalogo/proto/hardware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock para CatalogoService
type mockCatalogoService struct {
	componentes []domain.Componente
	specsMap    map[string]*domain.Componente
	errListar   error
	errSpecs    error
}

func (m *mockCatalogoService) ListarComponentes(tipo string) ([]domain.Componente, error) {
	if m.errListar != nil {
		return nil, m.errListar
	}
	if tipo == "invalido" {
		return nil, domain.ErrTipoInvalido
	}
	var res []domain.Componente
	for _, c := range m.componentes {
		if tipo == "" || string(c.Tipo) == tipo {
			res = append(res, c)
		}
	}
	return res, nil
}

func (m *mockCatalogoService) ObtenerSpecs(id string) (*domain.Componente, error) {
	if m.errSpecs != nil {
		return nil, m.errSpecs
	}
	if c, ok := m.specsMap[id]; ok {
		return c, nil
	}
	return nil, domain.ErrComponenteNoEncontrado
}

// Mock para CompatibilidadService
type mockCompatibilidadService struct {
	resultado *domain.ResultadoCompatibilidad
	err       error
}

func (m *mockCompatibilidadService) VerificarCompatibilidad(build domain.BuildEspecificacion) (*domain.ResultadoCompatibilidad, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resultado, nil
}

func (m *mockCompatibilidadService) CalcularConsumoEstimado(cpuTDP, gpuTDP int) int {
	return 350
}

func setupTestFixtures() (*mockCatalogoService, *mockCompatibilidadService, *grpcTransport.Handler) {
	socketAM5 := "AM5"
	ramDDR5 := "DDR5"

	catalogoSvc := &mockCatalogoService{
		componentes: []domain.Componente{
			{
				ID:         "cpu-9600x",
				Tipo:       domain.TipoCPU,
				Nombre:     "Ryzen 5 9600X",
				Fabricante: "AMD",
				TDP:        65,
				Socket:     &socketAM5,
				Origen:     "seed",
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
					"slotsRam": 4,
				},
				Benchmarks: domain.JSONMap{
					"vrm": 95,
				},
			},
		},
		specsMap: map[string]*domain.Componente{
			"cpu-9600x": {
				ID:         "cpu-9600x",
				Tipo:       domain.TipoCPU,
				Nombre:     "Ryzen 5 9600X",
				Fabricante: "AMD",
				TDP:        65,
				Socket:     &socketAM5,
				Origen:     "seed",
				Especificaciones: domain.JSONMap{
					"cores": 6,
				},
			},
		},
	}

	compatibilidadSvc := &mockCompatibilidadService{
		resultado: &domain.ResultadoCompatibilidad{
			Compatible:       true,
			Errores:          []string{},
			ConsumoEstimadoW: 354,
		},
	}

	handler := grpcTransport.NewHandler(catalogoSvc, compatibilidadSvc)
	return catalogoSvc, compatibilidadSvc, handler
}

func TestHandler_ListarComponentes_Exito(t *testing.T) {
	_, _, handler := setupTestFixtures()
	ctx := context.Background()

	// Listar todos
	res, err := handler.ListarComponentes(ctx, &hardwarepb.ListarComponentesRequest{})
	if err != nil {
		t.Fatalf("error inesperado al listar: %v", err)
	}
	if len(res.GetComponentes()) != 2 {
		t.Errorf("se esperaban 2 componentes, obtenidos %d", len(res.GetComponentes()))
	}

	// Listar con nil request
	resNil, err := handler.ListarComponentes(ctx, nil)
	if err != nil {
		t.Fatalf("error inesperado con nil request: %v", err)
	}
	if len(resNil.GetComponentes()) != 2 {
		t.Errorf("se esperaban 2 componentes con nil req, obtenidos %d", len(resNil.GetComponentes()))
	}

	// Filtrar por tipo
	resCPU, err := handler.ListarComponentes(ctx, &hardwarepb.ListarComponentesRequest{Tipo: "cpu"})
	if err != nil {
		t.Fatalf("error inesperado al filtrar: %v", err)
	}
	if len(resCPU.GetComponentes()) != 1 || resCPU.GetComponentes()[0].GetId() != "cpu-9600x" {
		t.Errorf("esperaba solo cpu-9600x, obtenido: %v", resCPU.GetComponentes())
	}
}

func TestHandler_ListarComponentes_Errores(t *testing.T) {
	catalogoSvc, _, handler := setupTestFixtures()
	ctx := context.Background()

	// Tipo inválido
	_, err := handler.ListarComponentes(ctx, &hardwarepb.ListarComponentesRequest{Tipo: "invalido"})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("esperaba código InvalidArgument, obtenido %v", status.Code(err))
	}

	// Error interno
	catalogoSvc.errListar = errors.New("db error")
	_, err = handler.ListarComponentes(ctx, &hardwarepb.ListarComponentesRequest{})
	if status.Code(err) != codes.Internal {
		t.Errorf("esperaba código Internal, obtenido %v", status.Code(err))
	}
}

func TestHandler_ObtenerSpecs_Exito(t *testing.T) {
	_, _, handler := setupTestFixtures()
	ctx := context.Background()

	res, err := handler.ObtenerSpecs(ctx, &hardwarepb.ObtenerSpecsRequest{Id: "cpu-9600x"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.GetNombre() != "Ryzen 5 9600X" {
		t.Errorf("esperaba Ryzen 5 9600X, obtenido %s", res.GetNombre())
	}
	if res.GetEspecificacionesJson() == "" {
		t.Errorf("especificaciones JSON no deberían estar vacías")
	}
}

func TestHandler_ObtenerSpecs_Errores(t *testing.T) {
	catalogoSvc, _, handler := setupTestFixtures()
	ctx := context.Background()

	// Request nil
	_, err := handler.ObtenerSpecs(ctx, nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("esperaba InvalidArgument para nil request, obtenido %v", status.Code(err))
	}

	// ID vacío
	_, err = handler.ObtenerSpecs(ctx, &hardwarepb.ObtenerSpecsRequest{Id: "   "})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("esperaba InvalidArgument para ID vacío, obtenido %v", status.Code(err))
	}

	// No encontrado
	_, err = handler.ObtenerSpecs(ctx, &hardwarepb.ObtenerSpecsRequest{Id: "inexistente"})
	if status.Code(err) != codes.NotFound {
		t.Errorf("esperaba NotFound, obtenido %v", status.Code(err))
	}

	// Error interno
	catalogoSvc.errSpecs = errors.New("fallo en BD")
	_, err = handler.ObtenerSpecs(ctx, &hardwarepb.ObtenerSpecsRequest{Id: "cpu-9600x"})
	if status.Code(err) != codes.Internal {
		t.Errorf("esperaba Internal, obtenido %v", status.Code(err))
	}
}

func TestHandler_VerificarCompatibilidad_Exito(t *testing.T) {
	_, compatibilidadSvc, handler := setupTestFixtures()
	ctx := context.Background()

	req := &hardwarepb.VerificarCompatibilidadRequest{
		Build: &hardwarepb.BuildRequest{
			CpuId:         "cpu-9600x",
			MotherboardId: "mb-b650",
			RamId:         "ram-ddr5-32",
			GpuId:         "gpu-5060ti",
			PsuId:         "psu-750",
			RamCantidad:   2,
		},
	}

	res, err := handler.VerificarCompatibilidad(ctx, req)
	if err != nil {
		t.Fatalf("error inesperado en VerificarCompatibilidad: %v", err)
	}
	if !res.GetCompatible() {
		t.Errorf("esperaba Compatible=true, obtenido false")
	}
	if res.GetConsumoEstimadoW() != 354 {
		t.Errorf("esperaba 354W, obtenido %d", res.GetConsumoEstimadoW())
	}

	// Caso incompatible
	compatibilidadSvc.resultado = &domain.ResultadoCompatibilidad{
		Compatible:       false,
		Errores:          []string{"socket incompatible"},
		ConsumoEstimadoW: 400,
	}

	resIncomp, err := handler.VerificarCompatibilidad(ctx, req)
	if err != nil {
		t.Fatalf("error inesperado con build incompatible: %v", err)
	}
	if resIncomp.GetCompatible() {
		t.Errorf("esperaba Compatible=false")
	}
	if len(resIncomp.GetErrores()) != 1 {
		t.Errorf("esperaba 1 error, obtenido %v", resIncomp.GetErrores())
	}
}

func TestHandler_VerificarCompatibilidad_Errores(t *testing.T) {
	_, compatibilidadSvc, handler := setupTestFixtures()
	ctx := context.Background()

	// Nil request
	_, err := handler.VerificarCompatibilidad(ctx, nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("esperaba InvalidArgument para nil request, obtenido %v", status.Code(err))
	}

	// Nil build
	_, err = handler.VerificarCompatibilidad(ctx, &hardwarepb.VerificarCompatibilidadRequest{Build: nil})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("esperaba InvalidArgument para build nil, obtenido %v", status.Code(err))
	}

	// Error interno
	compatibilidadSvc.err = errors.New("fallo interno")
	req := &hardwarepb.VerificarCompatibilidadRequest{
		Build: &hardwarepb.BuildRequest{CpuId: "cpu-9600x"},
	}
	_, err = handler.VerificarCompatibilidad(ctx, req)
	if status.Code(err) != codes.Internal {
		t.Errorf("esperaba Internal, obtenido %v", status.Code(err))
	}
}

func TestMappers_CasosBorde(t *testing.T) {
	// ToComponenteResumen nil
	if grpcTransport.ToComponenteResumen(nil) != nil {
		t.Errorf("esperaba nil para entrada nil")
	}

	// ToComponenteDetalleResponse nil
	res, err := grpcTransport.ToComponenteDetalleResponse(nil)
	if err != nil || res != nil {
		t.Errorf("esperaba nil y sin error para nil")
	}

	// ToDomainBuildEspecificacion nil
	domBuild := grpcTransport.ToDomainBuildEspecificacion(nil)
	if domBuild.CPUId != "" {
		t.Errorf("esperaba CPUId vacío para nil")
	}

	// ToCompatibilidadResponse nil
	resCompat := grpcTransport.ToCompatibilidadResponse(nil)
	if resCompat.GetCompatible() {
		t.Errorf("esperaba Compatible=false para nil")
	}
}

