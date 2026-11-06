package repository

import (
	"errors"
	"strings"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"gorm.io/gorm"
)

// ComponenteRepository define el contrato para el acceso a datos de componentes.
type ComponenteRepository interface {
	Listar(tipo string) ([]domain.Componente, error)
	ObtenerPorID(id string) (*domain.Componente, error)
	ObtenerPorIDs(ids []string) ([]domain.Componente, error)
}

type gormComponenteRepository struct {
	db *gorm.DB
}

// NewComponenteRepository crea una nueva instancia del repositorio de componentes respaldado por GORM.
func NewComponenteRepository(db *gorm.DB) ComponenteRepository {
	return &gormComponenteRepository{db: db}
}

// Listar obtiene todos los componentes o filtra por tipo si se especifica.
func (r *gormComponenteRepository) Listar(tipo string) ([]domain.Componente, error) {
	var componentes []domain.Componente
	query := r.db.Model(&domain.Componente{})

	tipoLimpio := strings.ToLower(strings.TrimSpace(tipo))
	if tipoLimpio != "" {
		query = query.Where("tipo = ?", tipoLimpio)
	}

	if err := query.Order("nombre ASC").Find(&componentes).Error; err != nil {
		return nil, err
	}

	if componentes == nil {
		componentes = []domain.Componente{}
	}

	return componentes, nil
}

// ObtenerPorID busca un componente individual por su clave primaria.
func (r *gormComponenteRepository) ObtenerPorID(id string) (*domain.Componente, error) {
	var comp domain.Componente
	idLimpio := strings.TrimSpace(id)
	if idLimpio == "" {
		return nil, domain.ErrComponenteNoEncontrado
	}

	err := r.db.First(&comp, "id = ?", idLimpio).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrComponenteNoEncontrado
		}
		return nil, err
	}

	return &comp, nil
}

// ObtenerPorIDs recupera múltiples componentes cuyos IDs coincidan con los provistos.
func (r *gormComponenteRepository) ObtenerPorIDs(ids []string) ([]domain.Componente, error) {
	if len(ids) == 0 {
		return []domain.Componente{}, nil
	}

	idsFiltrados := make([]string, 0, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			idsFiltrados = append(idsFiltrados, trimmed)
		}
	}

	if len(idsFiltrados) == 0 {
		return []domain.Componente{}, nil
	}

	var componentes []domain.Componente
	if err := r.db.Where("id IN ?", idsFiltrados).Find(&componentes).Error; err != nil {
		return nil, err
	}

	if componentes == nil {
		componentes = []domain.Componente{}
	}

	return componentes, nil
}

