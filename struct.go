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

type superbloque struct {
	NombreHD               [100]byte
	ArbolVirtualCount      int64
	DetalleDirectorioCount int64
	InodosCount            int64
	BloquesCount           int64
	ArbolVirtualFree       int64
	DetalleDirectorioFree  int64
	InodosFree             int64
	BloquesFree            int64
	DateCreacion           [16]byte
	DateUltimoMontaje      [16]byte
	MontajesCount          int64
	InicioBMAV             int64
	InicioAV               int64
	InicioBMDD             int64
	InicioDD               int64
	InicioBMInodos         int64
	InicioInodos           int64
	InicioBMBloques        int64
	InicioBloques          int64
	InicioLog              int64
	TamAV                  int64
	TamDD                  int64
	TamInodo               int64
	TamBloque              int64
	PrimerLibreAV          int64
	PrimerLibreDD          int64
	PrimerLibreInodo       int64
	PrimerLibreBloque      int64
	MagicNum               int64
}

type avd struct {
	AVDFechaCreacion            [16]byte
	AVDNombreDirectorio         [16]byte
	AVDApArraySubdirectorios    [6]int64
	AVDApDetalleDirectorio      int64
	AVDApArbolVirtualDirectorio int64
	AVDProper                   [16]byte
}

type dd struct {
	DDArrayFiles          [5]file
	DDApDetalleDirectorio int64
}

type inodo struct {
	ICountInodo            int64
	ISizeArchivo           int64
	ICountBloquesAsignados int64
	IArrayBloques          [4]int64
	IApIndirecto           int64
	IIdProper              [16]byte
}

type bloque struct {
	DBData [25]byte
}

type bitacora struct {
	LogTipoOperacion int64
	LogTipo          int64
	LogNombre        [16]byte
	LogContenido     int64
	LogFecha         [16]byte
}

type file struct {
	FileNombre           [16]byte
	FileApInodo          int64
	FileDateCreacion     [16]byte
	FileDateModificacion [16]byte
}
