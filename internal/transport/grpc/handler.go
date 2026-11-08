package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/service"
	hardwarepb "github.com/DI3G00xh/hardwareguru-catalogo/proto/hardware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler implementa el servidor gRPC de CatalogoService.
type Handler struct {
	hardwarepb.UnimplementedCatalogoServiceServer
	catalogoService       service.CatalogoService
	compatibilidadService service.CompatibilidadService
}

func NewHandler(catalogoSvc service.CatalogoService, compatibilidadSvc service.CompatibilidadService) *Handler {
	return &Handler{
		catalogoService:       catalogoSvc,
		compatibilidadService: compatibilidadSvc,
	}
}

// ListarComponentes filtra los componentes por tipo o retorna todo el catálogo.
func (h *Handler) ListarComponentes(ctx context.Context, req *hardwarepb.ListarComponentesRequest) (*hardwarepb.ListarComponentesResponse, error) {
	tipo := ""
	if req != nil {
		tipo = req.GetTipo()
	}

	componentes, err := h.catalogoService.ListarComponentes(tipo)
	if err != nil {
		if errors.Is(err, domain.ErrTipoInvalido) {
			return nil, status.Errorf(codes.InvalidArgument, "tipo de componente inválido: '%s'", tipo)
		}
		return nil, status.Errorf(codes.Internal, "error interno al listar componentes: %v", err)
	}

	return &hardwarepb.ListarComponentesResponse{
		Componentes: ToComponenteResumenList(componentes),
	}, nil
}

// ObtenerSpecs retorna las especificaciones detalladas y benchmarks de un componente en formato JSON.
func (h *Handler) ObtenerSpecs(ctx context.Context, req *hardwarepb.ObtenerSpecsRequest) (*hardwarepb.ComponenteDetalleResponse, error) {
	if req == nil || strings.TrimSpace(req.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "el id del componente es requerido")
	}

	id := strings.TrimSpace(req.GetId())
	comp, err := h.catalogoService.ObtenerSpecs(id)
	if err != nil {
		if errors.Is(err, domain.ErrComponenteNoEncontrado) {
			return nil, status.Errorf(codes.NotFound, "componente no encontrado con ID: '%s'", id)
		}
		return nil, status.Errorf(codes.Internal, "error interno al obtener especificaciones: %v", err)
	}

	return ToComponenteDetalleResponse(comp)
}

// VerificarCompatibilidad valida la viabilidad técnica y consumo de una configuración de hardware armada.
func (h *Handler) VerificarCompatibilidad(ctx context.Context, req *hardwarepb.VerificarCompatibilidadRequest) (*hardwarepb.CompatibilidadResponse, error) {
	if req == nil || req.GetBuild() == nil {
		return nil, status.Error(codes.InvalidArgument, "la configuración de la build es requerida")
	}

	buildDominio := ToDomainBuildEspecificacion(req.GetBuild())
	resultado, err := h.compatibilidadService.VerificarCompatibilidad(buildDominio)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error al verificar la compatibilidad: %v", err)
	}

	return ToCompatibilidadResponse(resultado), nil
}

