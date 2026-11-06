package domain

type BuildEspecificacion struct {
	CPUId         string `json:"cpuId"`
	MotherboardId string `json:"motherboardId"`
	RAMId         string `json:"ramId"`
	GPUId         string `json:"gpuId,omitempty"`
	PSUId         string `json:"psuId"`
	RAMCantidad   int    `json:"ramCantidad,omitempty"`
}

type ResultadoCompatibilidad struct {
	Compatible       bool     `json:"compatible"`
	Errores          []string `json:"errores"`
	ConsumoEstimadoW int      `json:"consumoEstimadoW"`
}
