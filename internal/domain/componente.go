package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error al leer JSONB de la base de datos")
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

