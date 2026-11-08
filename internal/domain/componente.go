package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type TipoComponente string

const (
	TipoCPU    TipoComponente = "cpu"
	TipoPlaca  TipoComponente = "placa"
	TipoRAM    TipoComponente = "ram"
	TipoGPU    TipoComponente = "gpu"
	TipoFuente TipoComponente = "fuente"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("error al leer JSONB de la base de datos: tipo no soportado")
	}

	if len(bytes) == 0 {
		*j = make(JSONMap)
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type Componente struct {
	ID               string         `gorm:"primaryKey;size:64" json:"id"`
	Tipo             TipoComponente `gorm:"size:16;not null;index" json:"tipo"`
	Nombre           string         `gorm:"size:128;not null" json:"nombre"`
	Fabricante       string         `gorm:"size:64;not null" json:"fabricante"`
	TDP              int            `gorm:"column:tdp;default:0" json:"tdp"`
	Socket           *string        `gorm:"size:32;index" json:"socket,omitempty"`
	TipoRAM          *string        `gorm:"column:tipo_ram;size:16" json:"tipoRam,omitempty"`
	PotenciaW        int            `gorm:"column:potencia_w;default:0" json:"potenciaW"`
	Origen           string         `gorm:"size:16;not null;default:'seed'" json:"origen"`
	Especificaciones JSONMap        `gorm:"type:jsonb;not null;default:'{}'" json:"especificaciones"`
	Benchmarks       JSONMap        `gorm:"type:jsonb;default:'{}'" json:"benchmarks,omitempty"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
}

// ObtenerSocket retorna el socket normalizado en mayúsculas y sin espacios.
func (c *Componente) ObtenerSocket() string {
	if c.Socket != nil && *c.Socket != "" {
		return strings.ToUpper(strings.TrimSpace(*c.Socket))
	}
	if val, ok := c.Especificaciones["socket"]; ok && val != nil {
		if s, ok := val.(string); ok {
			return strings.ToUpper(strings.TrimSpace(s))
		}
	}
	return ""
}

// ObtenerTipoRAM retorna la tecnología de memoria normalizada
func (c *Componente) ObtenerTipoRAM() string {
	if c.TipoRAM != nil && *c.TipoRAM != "" {
		return strings.ToUpper(strings.TrimSpace(*c.TipoRAM))
	}
	for _, key := range []string{"tipoRam", "tipo_ram", "tipoRAM"} {
		if val, ok := c.Especificaciones[key]; ok && val != nil {
			if s, ok := val.(string); ok {
				return strings.ToUpper(strings.TrimSpace(s))
			}
		}
	}
	return ""
}

// ObtenerSlotsRAM retorna la cantidad de ranuras de RAM disponibles en la placa madre.
func (c *Componente) ObtenerSlotsRAM() int {
	keys := []string{"slotsRam", "slots_ram", "slotsRAM", "ranurasRam", "ranuras_ram"}
	for _, key := range keys {
		if val, ok := c.Especificaciones[key]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				return int(v)
			case int:
				return v
			case int32:
				return int(v)
			case int64:
				return int(v)
			case string:
				if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
					return parsed
				}
			}
		}
	}
	return 0
}

// ObtenerModulosRAM retorna la cantidad de módulos que incluye el componente RAM (por defecto 1).
func (c *Componente) ObtenerModulosRAM() int {
	keys := []string{"modulos", "sticks", "cantidadModulos", "modulosPorKit"}
	for _, key := range keys {
		if val, ok := c.Especificaciones[key]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				if int(v) > 0 {
					return int(v)
				}
			case int:
				if v > 0 {
					return v
				}
			case int32:
				if v > 0 {
					return int(v)
				}
			case int64:
				if v > 0 {
					return int(v)
				}
			case string:
				if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && parsed > 0 {
					return parsed
				}
			}
		}
	}
	return 1
}

// ObtenerTDP retorna el TDP en Watts del componente.
func (c *Componente) ObtenerTDP() int {
	if c.TDP > 0 {
		return c.TDP
	}
	if val, ok := c.Especificaciones["tdp"]; ok && val != nil {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int32:
			return int(v)
		case int64:
			return int(v)
		case string:
			if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
				return parsed
			}
		}
	}
	return 0
}

// ObtenerPotenciaW retorna la potencia en Watts del componente
func (c *Componente) ObtenerPotenciaW() int {
	if c.PotenciaW > 0 {
		return c.PotenciaW
	}
	keys := []string{"potenciaW", "potencia_w", "potencia", "watts"}
	for _, key := range keys {
		if val, ok := c.Especificaciones[key]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				return int(v)
			case int:
				return v
			case int32:
				return int(v)
			case int64:
				return int(v)
			case string:
				if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
					return parsed
				}
			}
		}
	}
	return 0
}
