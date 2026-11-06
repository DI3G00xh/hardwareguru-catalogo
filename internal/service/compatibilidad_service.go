package service

import (
	"fmt"
	"math"
	"strings"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/repository"
)

const (
	MargenConsumoBaseWatts  = 50
	FactorSeguridadPotencia = 1.2
)

// TecnologiasRAMSoportadas define los estándares de memoria RAM contemplados en la validación
var TecnologiasRAMSoportadas = map[string]bool{
	"DDR3": true,
	"DDR4": true,
	"DDR5": true,
}

// CompatibilidadService define el contrato para validar la compatibilidad técnica entre piezas de una build.
type CompatibilidadService interface {
	VerificarCompatibilidad(build domain.BuildEspecificacion) (*domain.ResultadoCompatibilidad, error)
	CalcularConsumoEstimado(cpuTDP, gpuTDP int) int
}

type compatibilidadService struct {
	repo repository.ComponenteRepository
}

// NewCompatibilidadService crea una instancia del servicio de validación de compatibilidad.
func NewCompatibilidadService(repo repository.ComponenteRepository) CompatibilidadService {
	return &compatibilidadService{repo: repo}
}

// CalcularConsumoEstimado aplica la fórmula oficial del catálogo:
// Consumo Estimado = Ceil((TDP_cpu + TDP_gpu + 50W) * 1.2)
func (s *compatibilidadService) CalcularConsumoEstimado(cpuTDP, gpuTDP int) int {
	consumoSinMargen := cpuTDP + gpuTDP + MargenConsumoBaseWatts
	return int(math.Ceil(float64(consumoSinMargen) * FactorSeguridadPotencia))
}

// VerificarCompatibilidad evalúa integralmente las reglas de hardware
func (s *compatibilidadService) VerificarCompatibilidad(build domain.BuildEspecificacion) (*domain.ResultadoCompatibilidad, error) {
	var errores []string

	cpuID := strings.TrimSpace(build.CPUId)
	mbID := strings.TrimSpace(build.MotherboardId)
	ramID := strings.TrimSpace(build.RAMId)
	psuID := strings.TrimSpace(build.PSUId)
	gpuID := strings.TrimSpace(build.GPUId)

	if cpuID == "" {
		errores = append(errores, "Debe seleccionar un procesador (cpu_id)")
	}
	if mbID == "" {
		errores = append(errores, "Debe seleccionar una placa madre (motherboard_id)")
	}
	if ramID == "" {
		errores = append(errores, "Debe seleccionar una memoria RAM (ram_id)")
	}
	if psuID == "" {
		errores = append(errores, "Debe seleccionar una fuente de poder (psu_id)")
	}

	if len(errores) > 0 {
		return &domain.ResultadoCompatibilidad{
			Compatible:       false,
			Errores:          errores,
			ConsumoEstimadoW: 0,
		}, nil
	}

	idsParaBuscar := []string{cpuID, mbID, ramID, psuID}
	if gpuID != "" {
		idsParaBuscar = append(idsParaBuscar, gpuID)
	}

	componentes, err := s.repo.ObtenerPorIDs(idsParaBuscar)
	if err != nil {
		return nil, fmt.Errorf("error al consultar componentes para verificar compatibilidad: %w", err)
	}

	compMap := make(map[string]domain.Componente, len(componentes))
	for _, c := range componentes {
		compMap[c.ID] = c
	}

	cpu, okCPU := compMap[cpuID]
	if !okCPU {
		errores = append(errores, fmt.Sprintf("Procesador con ID '%s' no encontrado en el catálogo", cpuID))
	}
	mb, okMB := compMap[mbID]
	if !okMB {
		errores = append(errores, fmt.Sprintf("Placa madre con ID '%s' no encontrada en el catálogo", mbID))
	}
	ram, okRAM := compMap[ramID]
	if !okRAM {
		errores = append(errores, fmt.Sprintf("Memoria RAM con ID '%s' no encontrada en el catálogo", ramID))
	}
	psu, okPSU := compMap[psuID]
	if !okPSU {
		errores = append(errores, fmt.Sprintf("Fuente de poder con ID '%s' no encontrada en el catálogo", psuID))
	}

	var gpu domain.Componente
	var okGPU bool
	if gpuID != "" {
		gpu, okGPU = compMap[gpuID]
		if !okGPU {
			errores = append(errores, fmt.Sprintf("Tarjeta gráfica con ID '%s' no encontrada en el catálogo", gpuID))
		}
	}

	if len(errores) > 0 {
		return &domain.ResultadoCompatibilidad{
			Compatible:       false,
			Errores:          errores,
			ConsumoEstimadoW: 0,
		}, nil
	}

	if cpu.Tipo != domain.TipoCPU {
		errores = append(errores, fmt.Sprintf("El componente '%s' no es un procesador (tipo actual: '%s')", cpu.ID, cpu.Tipo))
	}
	if mb.Tipo != domain.TipoPlaca {
		errores = append(errores, fmt.Sprintf("El componente '%s' no es una placa madre (tipo actual: '%s')", mb.ID, mb.Tipo))
	}
	if ram.Tipo != domain.TipoRAM {
		errores = append(errores, fmt.Sprintf("El componente '%s' no es una memoria RAM (tipo actual: '%s')", ram.ID, ram.Tipo))
	}
	if psu.Tipo != domain.TipoFuente {
		errores = append(errores, fmt.Sprintf("El componente '%s' no es una fuente de poder (tipo actual: '%s')", psu.ID, psu.Tipo))
	}
	if gpuID != "" && gpu.Tipo != domain.TipoGPU {
		errores = append(errores, fmt.Sprintf("El componente '%s' no es una tarjeta gráfica (tipo actual: '%s')", gpu.ID, gpu.Tipo))
	}

	if len(errores) > 0 {
		return &domain.ResultadoCompatibilidad{
			Compatible:       false,
			Errores:          errores,
			ConsumoEstimadoW: 0,
		}, nil
	}

	// Regla 1: Coincidencia exacta del socket entre procesador y placa madre
	cpuSocket := cpu.ObtenerSocket()
	mbSocket := mb.ObtenerSocket()
	if cpuSocket == "" || mbSocket == "" {
		errores = append(errores, "No se pudo determinar el socket del procesador o de la placa madre")
	} else if !strings.EqualFold(cpuSocket, mbSocket) {
		errores = append(errores, fmt.Sprintf("Incompatibilidad de socket: el procesador '%s' requiere socket '%s', pero la placa madre '%s' tiene socket '%s'", cpu.Nombre, cpuSocket, mb.Nombre, mbSocket))
	}

	// Regla 2: Validar compatibilidad de la RAM con la placa según soporte para tecnología DDR3, DDR4, DDR5
	mbRAMTech := mb.ObtenerTipoRAM()
	ramTech := ram.ObtenerTipoRAM()
	if mbRAMTech == "" || ramTech == "" {
		errores = append(errores, "No se pudo determinar la tecnología de memoria RAM de la placa madre o de la memoria")
	} else {
		if !TecnologiasRAMSoportadas[ramTech] {
			errores = append(errores, fmt.Sprintf("Tecnología de RAM no reconocida o fuera de estándar (DDR3/DDR4/DDR5): '%s'", ramTech))
		}
		if !TecnologiasRAMSoportadas[mbRAMTech] {
			errores = append(errores, fmt.Sprintf("Tecnología de RAM de la placa madre no reconocida o fuera de estándar (DDR3/DDR4/DDR5): '%s'", mbRAMTech))
		}
		if !strings.EqualFold(mbRAMTech, ramTech) {
			errores = append(errores, fmt.Sprintf("Incompatibilidad de memoria RAM: la placa madre '%s' soporta tecnología '%s', pero la memoria seleccionada '%s' es '%s'", mb.Nombre, mbRAMTech, ram.Nombre, ramTech))
		}
	}

	// Regla 3: Validar si la cantidad de ranuras disponibles en la placa es compatible con las RAM solicitadas
	slotsDisponibles := mb.ObtenerSlotsRAM()
	cantidadRAMSolicitada := build.RAMCantidad
	if cantidadRAMSolicitada <= 0 {
		cantidadRAMSolicitada = 1
	}
	modulosPorUnidad := ram.ObtenerModulosRAM()
	modulosTotales := cantidadRAMSolicitada * modulosPorUnidad

	if slotsDisponibles > 0 && modulosTotales > slotsDisponibles {
		errores = append(errores, fmt.Sprintf("Ranuras de memoria insuficientes: la placa madre '%s' dispone de %d ranuras de RAM, pero la configuración seleccionada requiere %d módulos", mb.Nombre, slotsDisponibles, modulosTotales))
	}

	// Regla 4: Calcular el consumo para determinar si la fuente dispone energía suficiente
	cpuTDP := cpu.ObtenerTDP()
	gpuTDP := 0
	if gpuID != "" && okGPU {
		gpuTDP = gpu.ObtenerTDP()
	}
	consumoEstimado := s.CalcularConsumoEstimado(cpuTDP, gpuTDP)

	potenciaFuente := psu.ObtenerPotenciaW()
	if potenciaFuente < consumoEstimado {
		errores = append(errores, fmt.Sprintf("Potencia de fuente insuficiente: el consumo estimado para la build es de %dW (incluyendo margen de seguridad del 20%%), pero la fuente '%s' solo entrega %dW", consumoEstimado, psu.Nombre, potenciaFuente))
	}

	return &domain.ResultadoCompatibilidad{
		Compatible:       len(errores) == 0,
		Errores:          errores,
		ConsumoEstimadoW: consumoEstimado,
	}, nil
}
