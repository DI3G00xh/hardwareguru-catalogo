package service

import (
	"strings"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/repository"
)

// CatalogoService define las operaciones de consulta del catálogo de componentes.
type CatalogoService interface {
	ListarComponentes(tipo string) ([]domain.Componente, error)
	ObtenerSpecs(id string) (*domain.Componente, error)
}

type catalogoService struct {
	repo repository.ComponenteRepository
}

// NewCatalogoService crea una nueva instancia del servicio de catálogo.
func NewCatalogoService(repo repository.ComponenteRepository) CatalogoService {
	return &catalogoService{repo: repo}
}

// ListarComponentes obtiene los componentes, con opción de filtrar por tipo (cpu, placa, ram, gpu, fuente).
func (s *catalogoService) ListarComponentes(tipo string) ([]domain.Componente, error) {
	tipoLimpio := strings.ToLower(strings.TrimSpace(tipo))
	if tipoLimpio != "" {
		switch domain.TipoComponente(tipoLimpio) {
		case domain.TipoCPU, domain.TipoPlaca, domain.TipoRAM, domain.TipoGPU, domain.TipoFuente:
			// Tipo válido
		default:
			return nil, domain.ErrTipoInvalido
		}
	}
	return s.repo.Listar(tipoLimpio)
}

// ObtenerSpecs recupera las especificaciones técnicas completas de un componente por su ID.
func (s *catalogoService) ObtenerSpecs(id string) (*domain.Componente, error) {
	idLimpio := strings.TrimSpace(id)
	if idLimpio == "" {
		return nil, domain.ErrComponenteNoEncontrado
	}
	return s.repo.ObtenerPorID(idLimpio)
}

