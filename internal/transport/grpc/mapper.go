package grpc

import (
	"encoding/json"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	hardwarepb "github.com/DI3G00xh/hardwareguru-catalogo/proto/hardware"
)

// ToComponenteResumen convierte una entidad de dominio Componente a su representación Protobuf resumida.
func ToComponenteResumen(c *domain.Componente) *hardwarepb.ComponenteResumen {
	if c == nil {
		return nil
	}

	socket := ""
	if c.Socket != nil {
		socket = *c.Socket
	}

	tipoRAM := ""
	if c.TipoRAM != nil {
		tipoRAM = *c.TipoRAM
	}

	return &hardwarepb.ComponenteResumen{
		Id:         c.ID,
		Tipo:       string(c.Tipo),
		Nombre:     c.Nombre,
		Fabricante: c.Fabricante,
		Tdp:        int32(c.TDP),
		Socket:     socket,
		TipoRam:    tipoRAM,
		PotenciaW:  int32(c.PotenciaW),
		Origen:     c.Origen,
	}
}

// ToComponenteResumenList convierte una lista de componentes de dominio a una lista Protobuf.
func ToComponenteResumenList(componentes []domain.Componente) []*hardwarepb.ComponenteResumen {
	res := make([]*hardwarepb.ComponenteResumen, len(componentes))
	for i := range componentes {
		res[i] = ToComponenteResumen(&componentes[i])
	}
	return res
}

// ToComponenteDetalleResponse convierte una entidad de dominio Componente a Protobuf incluyendo los JSON serializados.
func ToComponenteDetalleResponse(c *domain.Componente) (*hardwarepb.ComponenteDetalleResponse, error) {
	if c == nil {
		return nil, nil
	}

	socket := ""
	if c.Socket != nil {
		socket = *c.Socket
	}

	tipoRAM := ""
	if c.TipoRAM != nil {
		tipoRAM = *c.TipoRAM
	}

	specsJSON, err := json.Marshal(c.Especificaciones)
	if err != nil {
		specsJSON = []byte("{}")
	}

	benchmarksJSON, err := json.Marshal(c.Benchmarks)
	if err != nil {
		benchmarksJSON = []byte("{}")
	}

	return &hardwarepb.ComponenteDetalleResponse{
		Id:                   c.ID,
		Tipo:                 string(c.Tipo),
		Nombre:               c.Nombre,
		Fabricante:           c.Fabricante,
		Tdp:                  int32(c.TDP),
		Socket:               socket,
		TipoRam:              tipoRAM,
		PotenciaW:            int32(c.PotenciaW),
		Origen:               c.Origen,
		EspecificacionesJson: string(specsJSON),
		BenchmarksJson:       string(benchmarksJSON),
	}, nil
}

// ToDomainBuildEspecificacion transforma el mensaje Protobuf BuildRequest a la estructura interna de dominio.
func ToDomainBuildEspecificacion(b *hardwarepb.BuildRequest) domain.BuildEspecificacion {
	if b == nil {
		return domain.BuildEspecificacion{}
	}

	ramCantidad := int(b.GetRamCantidad())
	if ramCantidad <= 0 {
		ramCantidad = 1
	}

	return domain.BuildEspecificacion{
		CPUId:         b.GetCpuId(),
		MotherboardId: b.GetMotherboardId(),
		RAMId:         b.GetRamId(),
		GPUId:         b.GetGpuId(),
		PSUId:         b.GetPsuId(),
		RAMCantidad:   ramCantidad,
	}
}

// ToCompatibilidadResponse transforma el resultado de dominio a su representación Protobuf.
func ToCompatibilidadResponse(res *domain.ResultadoCompatibilidad) *hardwarepb.CompatibilidadResponse {
	if res == nil {
		return &hardwarepb.CompatibilidadResponse{
			Compatible:       false,
			Errores:          []string{"resultado de compatibilidad no disponible"},
			ConsumoEstimadoW: 0,
		}
	}

	errores := res.Errores
	if errores == nil {
		errores = []string{}
	}

	return &hardwarepb.CompatibilidadResponse{
		Compatible:       res.Compatible,
		Errores:          errores,
		ConsumoEstimadoW: int32(res.ConsumoEstimadoW),
	}
}

