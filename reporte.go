/**************************************************************
	Melyza Alejandra Rodriguez Contreras
	201314821
	Laboratorio de Manejo e implementacion de Archivos
	Segundo Semestre 2020
	Proyecto No. 1
***************************************************************/
package main

/**************************************************************
	Importaciones
***************************************************************/
import (
	"fmt"
	"strings"
	"unsafe"
)

/**************************************************************
	REPORTE DISK
***************************************************************/
func reporteDisk(pathCrear string, id string) {
	id = strings.ToLower(id)
	s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
	if len(s) > 3 {
		if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
			var letra byte = s[2][0]
			inputFmt := id[3:len(id)] + ""
			idParticion := atributoSize(inputFmt)
			if idParticion > 0 {
				numResult, _, path := desmontarParticion2(letra, int64(idParticion))
				if numResult == 1 {
					//mandar a cambiar el estado de la particion en el mbr\
					pathEnviar := ""
					numeroEnviar := 0
					for index := 0; index < len(path); index++ {
						if path[index] == 0 {
							numeroEnviar = index
							index = len(path) + 1
							break
						}
					}
					pathEnviar = BytesToString(path[0:numeroEnviar])

					//path
					//nombre
					//formato

					separacion := strings.Split(pathCrear, ".")
					if len(separacion) > 1 {
						formato := strings.ToLower(strings.TrimSpace(separacion[1]))
						ruta1 := strings.Split(strings.TrimSpace(separacion[0]), "/")
						pathDestino := ""
						if len(ruta1) > 0 {
							for i := 0; i < len(ruta1)-1; i++ {
								pathDestino = pathDestino + "/" + ruta1[i]
							}
						}
						nombreDestino := ruta1[len(ruta1)-1]
						crearDirectorioSiNoExiste(pathDestino)

						/*fmt.Println("el destino es ")
						fmt.Println(pathDestino)
						fmt.Println("El formato es")
						fmt.Println(formato)
						fmt.Println("El nombre es")
						fmt.Println(nombreDestino)*/

						graficarDISCO(pathEnviar, pathDestino, nombreDestino, formato)
						fmt.Println("RESULTADO: Grafica del disco creada con exito")
					} else {
						fmt.Println("RESULTADO: Error al crear la grafica, destino incorrecto")
					}
				}
			} else {
				fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
		}
	}
}

/**************************************************************
	REPORTE MBR
***************************************************************/
func reporteMBR(pathCrear string, id string) {
	id = strings.ToLower(id)
	s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
	if len(s) > 3 {
		if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
			var letra byte = s[2][0]
			inputFmt := id[3:len(id)] + ""
			idParticion := atributoSize(inputFmt)
			if idParticion > 0 {
				numResult, _, path := desmontarParticion2(letra, int64(idParticion))
				if numResult == 1 {
					//mandar a cambiar el estado de la particion en el mbr\
					pathEnviar := ""
					numeroEnviar := 0
					for index := 0; index < len(path); index++ {
						if path[index] == 0 {
							numeroEnviar = index
							index = len(path) + 1
							break
						}
					}
					pathEnviar = BytesToString(path[0:numeroEnviar])

					//path
					//nombre
					//formato

					separacion := strings.Split(pathCrear, ".")
					if len(separacion) > 1 {
						formato := strings.ToLower(strings.TrimSpace(separacion[1]))
						ruta1 := strings.Split(strings.TrimSpace(separacion[0]), "/")
						pathDestino := ""
						if len(ruta1) > 0 {
							for i := 0; i < len(ruta1)-1; i++ {
								pathDestino = pathDestino + "/" + ruta1[i]
							}
						}
						nombreDestino := ruta1[len(ruta1)-1]
						crearDirectorioSiNoExiste(pathDestino)

						/*fmt.Println("el destino es ")
						fmt.Println(pathDestino)
						fmt.Println("El formato es")
						fmt.Println(formato)
						fmt.Println("El nombre es")
						fmt.Println(nombreDestino)*/

						graficarMBR(pathEnviar, pathDestino, nombreDestino, formato)
						fmt.Println("RESULTADO: Grafica del MBR creada con exito")
					} else {
						fmt.Println("RESULTADO: Error al crear la grafica, destino incorrecto")
					}
				}
			} else {
				fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
		}
	}
}

/**************************************************************
	REPORTE SB
***************************************************************/
func reporteSB(pathCrear string, id string) {
	id = strings.ToLower(id)
	s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
	if len(s) > 3 {
		if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
			var letra byte = s[2][0]
			inputFmt := id[3:len(id)] + ""
			idParticion := atributoSize(inputFmt)
			if idParticion > 0 {
				numResult, nombreParticion, path := desmontarParticion2(letra, int64(idParticion))
				if numResult == 1 {
					//mandar a cambiar el estado de la particion en el mbr\
					pathEnviar := ""
					numeroEnviar := 0
					for index := 0; index < len(path); index++ {
						if path[index] == 0 {
							numeroEnviar = index
							index = len(path) + 1
							break
						}
					}
					pathEnviar = BytesToString(path[0:numeroEnviar])
					var tipoParticion int
					var inicioParticion int64
					//var tamParticion int64
					tipoParticion, inicioParticion, _ = obtenerDatosParticion(pathEnviar, *nombreParticion)
					//path
					//nombre
					//formato
					if tipoParticion == 2 {
						inicioParticion = inicioParticion + int64(unsafe.Sizeof(ebr{}))
					}

					separacion := strings.Split(pathCrear, ".")
					if len(separacion) > 1 {
						formato := strings.ToLower(strings.TrimSpace(separacion[1]))
						ruta1 := strings.Split(strings.TrimSpace(separacion[0]), "/")
						pathDestino := ""
						if len(ruta1) > 0 {
							for i := 0; i < len(ruta1)-1; i++ {
								pathDestino = pathDestino + "/" + ruta1[i]
							}
						}
						nombreDestino := ruta1[len(ruta1)-1]
						crearDirectorioSiNoExiste(pathDestino)

						/*fmt.Println("el destino es ")
						fmt.Println(pathDestino)
						fmt.Println("El formato es")
						fmt.Println(formato)
						fmt.Println("El nombre es")
						fmt.Println(nombreDestino)*/

						graficarSB(pathEnviar, inicioParticion, pathDestino, nombreDestino, formato)
						fmt.Println("RESULTADO: Grafica del superbloque creada con exito")
					} else {
						fmt.Println("RESULTADO: Error al crear la grafica, destino incorrecto")
					}
				}
			} else {
				fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
		}
	}
}

/**************************************************************
	REPORTE BITMAP AVD
***************************************************************/
func reporteBMAVD(pathCrear string, id string) {
	id = strings.ToLower(id)
	s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
	if len(s) > 3 {
		if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
			var letra byte = s[2][0]
			inputFmt := id[3:len(id)] + ""
			idParticion := atributoSize(inputFmt)
			if idParticion > 0 {
				numResult, nombreParticion, path := desmontarParticion2(letra, int64(idParticion))
				if numResult == 1 {
					//mandar a cambiar el estado de la particion en el mbr\
					pathEnviar := ""
					numeroEnviar := 0
					for index := 0; index < len(path); index++ {
						if path[index] == 0 {
							numeroEnviar = index
							index = len(path) + 1
							break
						}
					}
					pathEnviar = BytesToString(path[0:numeroEnviar])
					var tipoParticion int
					var inicioParticion int64
					//var tamParticion int64
					tipoParticion, inicioParticion, _ = obtenerDatosParticion(pathEnviar, *nombreParticion)
					//path
					//nombre
					//formato
					if tipoParticion == 2 {
						inicioParticion = inicioParticion + int64(unsafe.Sizeof(ebr{}))
					}

					separacion := strings.Split(pathCrear, ".")
					if len(separacion) > 1 {
						formato := strings.ToLower(strings.TrimSpace(separacion[1]))
						ruta1 := strings.Split(strings.TrimSpace(separacion[0]), "/")
						pathDestino := ""
						if len(ruta1) > 0 {
							for i := 0; i < len(ruta1)-1; i++ {
								pathDestino = pathDestino + "/" + ruta1[i]
							}
						}
						nombreDestino := ruta1[len(ruta1)-1]
						crearDirectorioSiNoExiste(pathDestino)

						/*fmt.Println("el destino es ")
						fmt.Println(pathDestino)
						fmt.Println("El formato es")
						fmt.Println(formato)
						fmt.Println("El nombre es")
						fmt.Println(nombreDestino)*/

						if strings.Compare(strings.ToLower(strings.TrimSpace(formato)), "txt") == 0 {
							graficarBitMapDirectorio(pathEnviar, inicioParticion, pathDestino, nombreDestino, formato)
							fmt.Println("RESULTADO: Reporte de BITMAP DE ARBOL DE DIRECTORIO creada con exito")
						} else {
							fmt.Println("RESULTADO: El formato del reporte debe ser TXT")
						}
					} else {
						fmt.Println("RESULTADO: Error al crear la grafica, destino incorrecto")
					}
				}
			} else {
				fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
		}
	}
}

/**************************************************************
	REPORTE BITMAP DD
***************************************************************/
func reporteBMDD(pathCrear string, id string) {
	id = strings.ToLower(id)
	s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
	if len(s) > 3 {
		if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
			var letra byte = s[2][0]
			inputFmt := id[3:len(id)] + ""
			idParticion := atributoSize(inputFmt)
			if idParticion > 0 {
				numResult, nombreParticion, path := desmontarParticion2(letra, int64(idParticion))
				if numResult == 1 {
					//mandar a cambiar el estado de la particion en el mbr\
					pathEnviar := ""
					numeroEnviar := 0
					for index := 0; index < len(path); index++ {
						if path[index] == 0 {
							numeroEnviar = index
							index = len(path) + 1
							break
						}
					}
					pathEnviar = BytesToString(path[0:numeroEnviar])
					var tipoParticion int
					var inicioParticion int64
					//var tamParticion int64
					tipoParticion, inicioParticion, _ = obtenerDatosParticion(pathEnviar, *nombreParticion)
					//path
					//nombre
					//formato
					if tipoParticion == 2 {
						inicioParticion = inicioParticion + int64(unsafe.Sizeof(ebr{}))
					}

					separacion := strings.Split(pathCrear, ".")
					if len(separacion) > 1 {
						formato := strings.ToLower(strings.TrimSpace(separacion[1]))
						ruta1 := strings.Split(strings.TrimSpace(separacion[0]), "/")
						pathDestino := ""
						if len(ruta1) > 0 {
							for i := 0; i < len(ruta1)-1; i++ {
								pathDestino = pathDestino + "/" + ruta1[i]
							}
						}
						nombreDestino := ruta1[len(ruta1)-1]
						crearDirectorioSiNoExiste(pathDestino)

						/*fmt.Println("el destino es ")
						fmt.Println(pathDestino)
						fmt.Println("El formato es")
						fmt.Println(formato)
						fmt.Println("El nombre es")
						fmt.Println(nombreDestino)*/

						if strings.Compare(strings.ToLower(strings.TrimSpace(formato)), "txt") == 0 {
							graficarBitMapDetalle(pathEnviar, inicioParticion, pathDestino, nombreDestino, formato)
							fmt.Println("RESULTADO: Reporte de BITMAP DE DETALLE DE DIRECTORIO creada con exito")
						} else {
							fmt.Println("RESULTADO: El formato del reporte debe ser TXT")
						}
					} else {
						fmt.Println("RESULTADO: Error al crear la grafica, destino incorrecto")
					}
				}
			} else {
				fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
		}
	}
}

/**************************************************************
	REPORTE BITMAP INODO
***************************************************************/
func reporteBMINODO(pathCrear string, id string) {
	id = strings.ToLower(id)
	s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
	if len(s) > 3 {
		if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
			var letra byte = s[2][0]
			inputFmt := id[3:len(id)] + ""
			idParticion := atributoSize(inputFmt)
			if idParticion > 0 {
				numResult, nombreParticion, path := desmontarParticion2(letra, int64(idParticion))
				if numResult == 1 {
					//mandar a cambiar el estado de la particion en el mbr\
					pathEnviar := ""
					numeroEnviar := 0
					for index := 0; index < len(path); index++ {
						if path[index] == 0 {
							numeroEnviar = index
							index = len(path) + 1
							break
						}
					}
					pathEnviar = BytesToString(path[0:numeroEnviar])
					var tipoParticion int
					var inicioParticion int64
					//var tamParticion int64
					tipoParticion, inicioParticion, _ = obtenerDatosParticion(pathEnviar, *nombreParticion)
					//path
					//nombre
					//formato
					if tipoParticion == 2 {
						inicioParticion = inicioParticion + int64(unsafe.Sizeof(ebr{}))
					}

					separacion := strings.Split(pathCrear, ".")
					if len(separacion) > 1 {
						formato := strings.ToLower(strings.TrimSpace(separacion[1]))
						ruta1 := strings.Split(strings.TrimSpace(separacion[0]), "/")
						pathDestino := ""
						if len(ruta1) > 0 {
							for i := 0; i < len(ruta1)-1; i++ {
								pathDestino = pathDestino + "/" + ruta1[i]
							}
						}
						nombreDestino := ruta1[len(ruta1)-1]
						crearDirectorioSiNoExiste(pathDestino)

						/*fmt.Println("el destino es ")
						fmt.Println(pathDestino)
						fmt.Println("El formato es")
						fmt.Println(formato)
						fmt.Println("El nombre es")
						fmt.Println(nombreDestino)*/

						if strings.Compare(strings.ToLower(strings.TrimSpace(formato)), "txt") == 0 {
							graficarBitMapInodo(pathEnviar, inicioParticion, pathDestino, nombreDestino, formato)
							fmt.Println("RESULTADO: Reporte de BITMAP DE INODOS creada con exito")
						} else {
							fmt.Println("RESULTADO: El formato del reporte debe ser TXT")
						}
					} else {
						fmt.Println("RESULTADO: Error al crear la grafica, destino incorrecto")
					}
				}
			} else {
				fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
		}
	}
}

/**************************************************************
	REPORTE BITMAP DE BLOQUES
***************************************************************/
func reporteBMBLOQUE(pathCrear string, id string) {
	id = strings.ToLower(id)
	s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
	if len(s) > 3 {
		if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
			var letra byte = s[2][0]
			inputFmt := id[3:len(id)] + ""
			idParticion := atributoSize(inputFmt)
			if idParticion > 0 {
				numResult, nombreParticion, path := desmontarParticion2(letra, int64(idParticion))
				if numResult == 1 {
					//mandar a cambiar el estado de la particion en el mbr\
					pathEnviar := ""
					numeroEnviar := 0
					for index := 0; index < len(path); index++ {
						if path[index] == 0 {
							numeroEnviar = index
							index = len(path) + 1
							break
						}
					}
					pathEnviar = BytesToString(path[0:numeroEnviar])
					var tipoParticion int
					var inicioParticion int64
					//var tamParticion int64
					tipoParticion, inicioParticion, _ = obtenerDatosParticion(pathEnviar, *nombreParticion)
					//path
					//nombre
					//formato
					if tipoParticion == 2 {
						inicioParticion = inicioParticion + int64(unsafe.Sizeof(ebr{}))
					}

					separacion := strings.Split(pathCrear, ".")
					if len(separacion) > 1 {
						formato := strings.ToLower(strings.TrimSpace(separacion[1]))
						ruta1 := strings.Split(strings.TrimSpace(separacion[0]), "/")
						pathDestino := ""
						if len(ruta1) > 0 {
							for i := 0; i < len(ruta1)-1; i++ {
								pathDestino = pathDestino + "/" + ruta1[i]
							}
						}
						nombreDestino := ruta1[len(ruta1)-1]
						crearDirectorioSiNoExiste(pathDestino)

						/*fmt.Println("el destino es ")
						fmt.Println(pathDestino)
						fmt.Println("El formato es")
						fmt.Println(formato)
						fmt.Println("El nombre es")
						fmt.Println(nombreDestino)*/

						if strings.Compare(strings.ToLower(strings.TrimSpace(formato)), "txt") == 0 {
							graficarBitMapBloque(pathEnviar, inicioParticion, pathDestino, nombreDestino, formato)
							fmt.Println("RESULTADO: Reporte de BITMAP DE BLOQUES creada con exito")
						} else {
							fmt.Println("RESULTADO: El formato del reporte debe ser TXT")
						}
					} else {
						fmt.Println("RESULTADO: Error al crear la grafica, destino incorrecto")
					}
				}
			} else {
				fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
		}
	}
}
