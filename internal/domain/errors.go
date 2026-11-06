package domain

import "errors"

var (
	ErrComponenteNoEncontrado = errors.New("componente no encontrado")
	ErrTipoInvalido           = errors.New("tipo de componente inválido")
	ErrBuildInvalida          = errors.New("la build no contiene todos los componentes requeridos")
)

