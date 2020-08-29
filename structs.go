/**************************************************************
	Melyza Alejandra Rodriguez Contreras
	201314821
	Laboratorio de Manejo e implementacion de Archivos
	Segundo Semestre 2020
	Proyecto No. 1
***************************************************************/
package main

/**************************************************************
	Definicion de structs
***************************************************************/
type mbr struct {
	Tamano int64
	Fecha  [16]byte
	Firma  int8
	Tabla  [4]particion
}

type particion struct {
	Status byte
	Type   byte
	Fit    byte
	Start  int64
	Size   int64
	Name   [16]byte
}

type ebr struct {
	Status byte
	Fit    byte
	Start  int64
	Size   int64
	Next   int64
	Name   [16]byte
}

type discoMontado struct {
	Path   [100]byte
	ID     byte
	Estado byte
	lista  [100]particionMontada
}

type particionMontada struct {
	ID            [4]byte
	nombre        [16]byte
	EstadoFormato byte
	EstadoMount   byte
}
