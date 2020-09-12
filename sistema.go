/**************************************************************
	Melyza Alejandra Rodriguez Contreras
	201314821
	Laboratorio de Manejo e implementacion de Archivos
	Segundo Semestre 2020
	Proyecto No. 1
***************************************************************/
package main

/**************************************************************
	Imports
***************************************************************/
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
	"math"
)

/**************************************************************
	Comando MOUNT PARTICION
***************************************************************/
func montarParticion(path string, nombre string) int {
	//leyendo el archivo
	mbrLeido := leerDisco(path)
	var tmpNombre [16]byte
	var tmpPath [100]byte
	copy(tmpNombre[:], nombre)
	copy(tmpPath[:], path)

	if mbrLeido != nil {
		//buscando en tabla principal*****************************************
		for i := 0; i < len(mbrLeido.Tabla); i++ {
			if tmpNombre == mbrLeido.Tabla[i].Name {
				if mbrLeido.Tabla[i].Type == 'e' { //es una particion extendida
					fmt.Println("RESULTADO: No se puede montar una particion extendida")
					return 0
				} else if mbrLeido.Tabla[i].Status == 1 {
					fmt.Println("RESULTADO: La particion ya se encuentra montada")
					return 0
				}
				//si sale es porque es una particion primaria
				nueva := particionMontada{}
				nueva.nombre = mbrLeido.Tabla[i].Name
				nuevoDisco := discoMontado{}
				nuevoDisco.Path = tmpPath
				nuevoDisco.Estado = 0
				nuevoDisco.ID = 0
				if subirParticion(&nueva, &nuevoDisco) == 1 {
					mbrLeido.Tabla[i].Status = 1
					reescribir(mbrLeido, path)
					//graficarMBR(path)
					return 1
				}
				return 0
			}
		}
		//buscando en las particiones logicas*********************************
		for i := 0; i < len(mbrLeido.Tabla); i++ {
			if mbrLeido.Tabla[i].Type == 'e' {
				//empiezo a verificar los nombres de las logicas
				//leyendo el archivo
				ebrTemp := ebr{}
				file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
				defer file.Close()
				if err != nil {
					log.Fatal(err)
				}
				file.Seek(mbrLeido.Tabla[i].Start, 0)
				data := readNextBytes(file, unsafe.Sizeof(ebr{}))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
				if err != nil {
					log.Fatal("binary.Read failed", err)
				}
				limite := mbrLeido.Tabla[i].Start + int64(unsafe.Sizeof(ebr{})) + mbrLeido.Tabla[i].Size
				if &ebrTemp != nil {
					for i := ebrTemp.Start; i < limite; i++ {
						ebrLeido := ebr{} //el ebr que lee en cada iteracion del for
						file.Seek(i, 0)
						data := readNextBytes(file, unsafe.Sizeof(ebr{}))
						buffer := bytes.NewBuffer(data)
						err = binary.Read(buffer, binary.BigEndian, &ebrLeido)
						if err != nil {
							log.Fatal("binary.Read failed", err)
						}
						if &ebrLeido != nil {
							if ebrLeido.Name == tmpNombre {
								nueva := particionMontada{}
								nueva.nombre = ebrLeido.Name
								nuevoDisco := discoMontado{}
								nuevoDisco.Path = tmpPath
								nuevoDisco.Estado = 0
								nuevoDisco.ID = 0
								if subirParticion(&nueva, &nuevoDisco) == 1 {
									ebrLeido.Status = 1
									escribirEbr(path, ebrLeido.Start, &ebrLeido)
									return 1
								}
								return 0
							}

							if ebrLeido.Next == -1 {
								i = limite + 1
								fmt.Println("RESULTADO: No se encuentra la particion especificada")
								return 0
							}
							i = ebrLeido.Next - 1
						}

					}
				}
			}
		}
		//si llega hasta aqui es porque no se encontro la particion
		fmt.Println("RESULTADO: La particion especificada no se encuentra en el disco")
		return 0
	}
	return 0
}

func subirParticion(nueva *particionMontada, nuevoDisco *discoMontado) int {
	//verificar si el coso esta vacio
	if discosMontados[0].ID == 0 {
		nuevoDisco.Estado = 1
		nuevoDisco.ID = 0 + 97 //es el ascii de la letra
		nueva.EstadoFormato = 0
		nueva.EstadoMount = 1
		nueva.ID[0] = 118
		nueva.ID[1] = 100
		nueva.ID[2] = nuevoDisco.ID
		nueva.ID[3] = 1
		nuevoDisco.lista[0] = *nueva
		discosMontados[0] = *nuevoDisco
		/*
			Mandar a escribir estatus de particion en el disco
		*/
		return 1
	}
	//Aqui se va a ver si el disco ya esta creado
	for i := 0; i < len(discosMontados); i++ {
		discoActual := discosMontados[i]
		if discoActual.ID != 0 {
			if discoActual.Path == nuevoDisco.Path {
				//El disco ya se encuentra creado
				//en busca del nombre a ver si esta repetido
				for j := 0; j < len(discoActual.lista); j++ {
					particionActual := discoActual.lista[j]
					if particionActual.nombre == nueva.nombre {
						fmt.Println("RESULTADO: La particion ya se encuentra montada")
						return 0
					}
				}
				//En este punto no esta repetido en nombre
				var idAnterior byte = 0
				for j := 0; j < len(discoActual.lista); j++ {
					particionActual := discoActual.lista[j]
					if particionActual.ID[3] == 0 {
						//aqui va a crear la nueva particion
						discoActual.Estado = 1
						particionActual.EstadoFormato = 0
						particionActual.EstadoMount = 1
						nueva.ID[0] = 118
						nueva.ID[1] = 100
						nueva.ID[2] = discoActual.ID
						nueva.ID[3] = idAnterior + 1
						discoActual.lista[j] = *nueva
						discosMontados[i] = discoActual
						return 1
					}
					idAnterior = particionActual.ID[3]
				}
				fmt.Println("RESULTADO: No se ha podido montar la particion")
				return 0
			}
		}
	}
	//En este punto no esta creado el disco pero existe al menos uno ya creado
	var letraAnterior byte
	for i := 0; i < len(discosMontados); i++ {
		discoActual := discosMontados[i]
		if discoActual.ID == 0 {
			nuevoDisco.Estado = 1
			nuevoDisco.ID = letraAnterior + 1 //es el ascii de la letra
			nueva.EstadoFormato = 0
			nueva.EstadoMount = 1
			nueva.ID[0] = 118
			nueva.ID[1] = 100
			nueva.ID[2] = nuevoDisco.ID
			nueva.ID[3] = 1
			nuevoDisco.lista[0] = *nueva
			discosMontados[i] = *nuevoDisco
			/*
				Mandar a escribir estatus de particion en el disco
			*/
			return 1
		}
		letraAnterior = discoActual.ID
	}

	fmt.Println("RESULTADO: No se ha podido montar la particion (FIN)")
	return 0
}

func retornarLetra(numero byte) byte {
	numero = 97 - numero
	switch numero {
	case 0:
		return 'a'
	case 1:
		return 'b'
	case 2:
		return 'c'
	case 3:
		return 'd'
	case 4:
		return 'e'
	case 5:
		return 'f'
	case 6:
		return 'g'
	case 7:
		return 'h'
	case 8:
		return 'i'
	case 9:
		return 'j'
	case 10:
		return 'k'
	case 11:
		return 'l'
	case 12:
		return 'm'
	case 13:
		return 'n'
	case 14:
		return 'o'
	case 15:
		return 'p'
	case 16:
		return 'q'
	case 17:
		return 'r'
	case 18:
		return 's'
	case 19:
		return 't'
	case 20:
		return 'u'
	case 21:
		return 'v'
	case 22:
		return 'w'
	case 23:
		return 'x'
	case 24:
		return 'y'
	case 25:
		return 'z'
	default:
		return 'A'
	}
}

/**************************************************************
	Comando MOUNT IMPRIMIR
	Se debe imprimir
	-id->
	-path->
	-name->
***************************************************************/
func imprimirMOUNT() {
	fmt.Println()
	fmt.Println("********************************PARTICIONES MONTADAS********************************")
	fmt.Println()
	for i := 0; i < len(discosMontados); i++ {
		discoActual := discosMontados[i]
		if discoActual.ID != 0 {
			//el disco esta lleno
			for j := 0; j < len(discoActual.lista); j++ {
				if discoActual.lista[j].ID[3] != 0 {
					fmt.Print("id->")
					bArray2 := discoActual.lista[j].ID
					bArray3 := []byte{bArray2[0], bArray2[1], bArray2[2], bArray2[3]}
					str2 := BytesToString(bArray3)
					fmt.Print(str2)
					fmt.Print(int64(bArray2[3]))
					fmt.Print(" path->")
					str3 := string(discoActual.Path[:])
					fmt.Print(str3)
					fmt.Print(" -name->")
					str4 := string(discoActual.lista[j].nombre[:])
					fmt.Println(str4)
				}
			}
		}
	}
	fmt.Println()
	fmt.Println("************************************************************************************")
	fmt.Println()
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

/**************************************************************
	Comando UNMOUNT
***************************************************************/
func desmontar(lista []string) {
	for i := 0; i < len(lista); i++ {
		s := strings.Split(strings.ToLower(strings.TrimSpace(lista[i])), "")
		if len(s) > 3 {
			if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
				var letra byte = s[2][0]
				inputFmt := lista[i][3:len(lista[i])] + ""
				idParticion := atributoSize(inputFmt)
				if idParticion > 0 {
					numResult, nombre, path := desmontarParticion(letra, int64(idParticion))
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
						var pa string = pathEnviar
						actualizarEstado(strings.TrimSpace(pa), nombre)
						fmt.Println("RESULTADO: Se ha desmontado la particion con exito ")
					}
				} else {
					fmt.Println("RESULTADO: Error en el id de la particion a desmontar")
				}
			} else {
				fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
			}
		}
	}
}

func desmontarParticion(letra byte, numero int64) (int, *[16]byte, *[100]byte) {
	for i := 0; i < len(discosMontados); i++ {
		//buscar la letra
		if discosMontados[i].ID != 0 { //solo va a verificar los discos que esten creados
			discoActual := discosMontados[i]
			if discoActual.ID == letra {
				//se encontro el disco
				for j := 0; j < len(discoActual.lista); j++ {
					//buscando el numero
					if discoActual.lista[j].ID[0] != 0 { //solo va a buscar en los que contengan
						if int64(discoActual.lista[j].ID[3]) == numero {
							//encontro la particion a desmontar
							particionMontadaVacia := particionMontada{}
							ret := discoActual.lista[j].nombre
							ret2 := discoActual.Path
							discoActual.lista[j] = particionMontadaVacia
							discosMontados[i] = discoActual
							return 1, &ret, &ret2
						}
					}
				}
			}
		}
	}
	fmt.Println("RESULTADO: No se ha encontrado la particion especificada")
	return 0, nil, nil
}

func desmontarParticion2(letra byte, numero int64) (int, *[16]byte, *[100]byte) {
	for i := 0; i < len(discosMontados); i++ {
		//buscar la letra
		if discosMontados[i].ID != 0 { //solo va a verificar los discos que esten creados
			discoActual := discosMontados[i]
			if discoActual.ID == letra {
				//se encontro el disco
				for j := 0; j < len(discoActual.lista); j++ {
					//buscando el numero
					if discoActual.lista[j].ID[0] != 0 { //solo va a buscar en los que contengan
						if int64(discoActual.lista[j].ID[3]) == numero {
							//encontro la particion a desmontar
							ret := discoActual.lista[j].nombre
							ret2 := discoActual.Path
							return 1, &ret, &ret2
						}
					}
				}
			}
		}
	}
	return 0, nil, nil
}

func actualizarEstado(path string, nombre *[16]byte) {
	//leyendo el archivo
	m := mbr{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		file.Seek(0, 0)
		data := readNextBytes(file, unsafe.Sizeof(mbr{}))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
	}
	var mbrLeido *mbr = &m

	//mbrLeido := leerDisco(strings.ReplaceAll(path, "\"", ""))
	if mbrLeido != nil {

		//buscando en tabla principal*****************************************
		for i := 0; i < len(mbrLeido.Tabla); i++ {
			if *nombre == mbrLeido.Tabla[i].Name {
				mbrLeido.Tabla[i].Status = 0
				reescribir(mbrLeido, path)
				return
			}
		}
		//buscando en las particiones logicas*********************************
		for i := 0; i < len(mbrLeido.Tabla); i++ {
			if mbrLeido.Tabla[i].Type == 'e' {
				//empiezo a verificar los nombres de las logicas
				//leyendo el archivo
				ebrTemp := ebr{}
				file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
				defer file.Close()
				if err != nil {
					log.Fatal(err)
				}
				file.Seek(mbrLeido.Tabla[i].Start, 0)
				data := readNextBytes(file, unsafe.Sizeof(ebr{}))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
				if err != nil {
					log.Fatal("binary.Read failed", err)
				}
				limite := mbrLeido.Tabla[i].Start + int64(unsafe.Sizeof(ebr{})) + mbrLeido.Tabla[i].Size
				if &ebrTemp != nil {
					for i := ebrTemp.Start; i < limite; i++ {
						ebrLeido := ebr{} //el ebr que lee en cada iteracion del for
						file.Seek(i, 0)
						data := readNextBytes(file, unsafe.Sizeof(ebr{}))
						buffer := bytes.NewBuffer(data)
						err = binary.Read(buffer, binary.BigEndian, &ebrLeido)
						if err != nil {
							log.Fatal("binary.Read failed", err)
						}
						if &ebrLeido != nil {
							if ebrLeido.Name == *nombre {
								ebrLeido.Status = 0
								escribirEbr(path, ebrLeido.Start, &ebrLeido)
								return
							}

							if ebrLeido.Next == -1 {
								i = limite + 1
								fmt.Println("RESULTADO: No se encuentra la particion especificada")
								return
							}
							i = ebrLeido.Next - 1
						}

					}
				}
			}
		}
		return
	}
	return
}

/**************************************************************
	FORMATEO
***************************************************************/
func formatear(idFormatear string, tipoFormato string) {
	var tamParticion int64
	var inicioParticion int64
	var tipoParticion int
	if strings.Compare(idFormatear, "") != 0 {
		s := strings.Split(strings.ToLower(strings.TrimSpace(idFormatear)), "")
		if len(s) > 3 {
			if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
				var letra byte = s[2][0]
				inputFmt := idFormatear[3:len(idFormatear)] + ""
				idParticion := atributoSize(inputFmt)
				if idParticion > 0 {
					numResult, nombre, path := desmontarParticion2(letra, int64(idParticion))
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
						//fmt.Println(pathEnviar)
						//fmt.Println(nombre)
						tipoParticion, inicioParticion, tamParticion = obtenerDatosParticion(pathEnviar, *nombre)
						if tipoParticion == 0 {
							fmt.Println("RESULTADO: No se encuentra la particion")
						} else if tipoParticion == 3 {
							fmt.Println("RESULTADO: No se puede formatear una particion extendida")
						} else if tipoParticion == 1 {
							estadoActual := getEstadoFormato(letra, int64(idParticion))
							if estadoActual == 5 {
								//hay un problema
								fmt.Println("RESULTADO: Existe un problema de formateo de la particion")
							} else if estadoActual == 0 {
								//no esta formateada
								resultado := realizarFormato(pathEnviar, inicioParticion, tamParticion, tipoFormato)
								if resultado == 1 {
									actualizarEstadoFormato(letra, int64(idParticion))
									fmt.Println("RESULTADO: Particion formateada con exito")
								}
							} else if estadoActual == 1 {
								//ya esta formateada pedir confirmacion
								fmt.Println("*****  ATENCION! LA PARTICION SE ENCUENTRA FORMATEADA   *****")
								fmt.Print("       Desea formatear 1) SI 2) NO : ")
								lector := bufio.NewReader(os.Stdin)
								comando, _ := lector.ReadString('\n')
								if strings.Compare(strings.TrimSpace(comando), "1") == 0 {
									//si formatear
									resultado := realizarFormato(pathEnviar, inicioParticion, tamParticion, tipoFormato)
									if resultado == 1 {
										actualizarEstadoFormato(letra, int64(idParticion))
										fmt.Println("RESULTADO: Particion formateada con exito")
									}
								} else if strings.Compare(strings.TrimSpace(comando), "2") == 0 {
									//no formatear
									fmt.Println("RESULTADO: La particion no ha sido formateada")
								} else {
									//error en opcion ingresada
									fmt.Println("RESULTADO: Se ha ingresado una opcion incorrecta, la particion no se ha formateado")
								}

							}
						} else if tipoParticion == 2 {
							//es una extendida
							estadoActual := getEstadoFormato(letra, int64(idParticion))
							if estadoActual == 5 {
								//hay un problema
								fmt.Println("RESULTADO: Existe un problema de formateo de la particion")
							} else if estadoActual == 0 {
								//no esta formateada
								resultado := realizarFormato(pathEnviar, inicioParticion+int64(unsafe.Sizeof(ebr{})), tamParticion, tipoFormato)
								if resultado == 1 {
									actualizarEstadoFormato(letra, int64(idParticion))
									fmt.Println("RESULTADO: Particion formateada con exito")
								}
							} else if estadoActual == 1 {
								//ya esta formateada pedir confirmacion
								fmt.Println("*****  ATENCION! LA PARTICION SE ENCUENTRA FORMATEADA   *****")
								fmt.Print("       Desea formatear 1) SI 2) NO : ")
								lector := bufio.NewReader(os.Stdin)
								comando, _ := lector.ReadString('\n')
								if strings.Compare(strings.TrimSpace(comando), "1") == 0 {
									//si formatear
									resultado := realizarFormato(pathEnviar, inicioParticion, tamParticion, tipoFormato)
									if resultado == 1 {
										actualizarEstadoFormato(letra, int64(idParticion))
										fmt.Print("RESULTADO: Particion formateada con exito")
									}
								} else if strings.Compare(strings.TrimSpace(comando), "2") == 0 {
									//no formatear
									fmt.Println("RESULTADO: La particion no ha sido formateada")
								} else {
									//error en opcion ingresada
									fmt.Println("RESULTADO: Se ha ingresado una opcion incorrecta, la particion no se ha formateado")
								}

							}
						}
					}
				} else {
					fmt.Println("RESULTADO: Error en el id de la particion a formatear")
				}
			} else {
				fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto")
			}
		}
	}

}

func getEstadoFormato(letra byte, numero int64) int {
	for i := 0; i < len(discosMontados); i++ {
		//buscar la letra
		if discosMontados[i].ID != 0 { //solo va a verificar los discos que esten creados
			discoActual := discosMontados[i]
			if discoActual.ID == letra {
				//se encontro el disco
				for j := 0; j < len(discoActual.lista); j++ {
					//buscando el numero
					if discoActual.lista[j].ID[0] != 0 { //solo va a buscar en los que contengan
						if int64(discoActual.lista[j].ID[3]) == numero {
							//encontro la particion a desmontar
							return int(discoActual.lista[j].EstadoFormato)
						}
					}
				}
				break
			}
		}
	}
	return 5
}

func actualizarEstadoFormato(letra byte, numero int64) {
	for i := 0; i < len(discosMontados); i++ {
		//buscar la letra
		if discosMontados[i].ID != 0 { //solo va a verificar los discos que esten creados
			discoActual := discosMontados[i]
			if discoActual.ID == letra {
				//se encontro el disco
				for j := 0; j < len(discoActual.lista); j++ {
					//buscando el numero
					if discoActual.lista[j].ID[0] != 0 { //solo va a buscar en los que contengan
						if int64(discoActual.lista[j].ID[3]) == numero {
							//encontro la particion a desmontar
							discosMontados[i].lista[j].EstadoFormato = 1
							return
						}
					}
				}
			}
		}
	}
	fmt.Println("RESULTADO: No se ha encontrado la particion especificada")
}

func obtenerDatosParticion(path string, nombre [16]byte) (int, int64, int64) {
	if strings.Compare(path, "") != 0 {
		s := leerDisco(path)
		if s != nil {
			//buscando en las particiones principales
			for i := 0; i < len(s.Tabla); i++ {
				if nombre == s.Tabla[i].Name {
					if s.Tabla[i].Type == 'p' {
						//es una primaria
						return 1, s.Tabla[i].Start, s.Tabla[i].Size
					}
					//es una extendida
					return 3, 0, 0
				}
			}
			//si logra salir del for es porque no la encontro en las principales
			//buscar entre las logicas
			for i := 0; i < len(s.Tabla); i++ {
				if s.Tabla[i].Type == 'e' {
					//encontro la particion extendida
					ebrTemp := ebr{}
					file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
					defer file.Close()
					if err != nil {
						log.Fatal(err)
					}
					file.Seek(s.Tabla[i].Start, 0)
					data := readNextBytes(file, unsafe.Sizeof(ebr{}))
					buffer := bytes.NewBuffer(data)
					err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
					if err != nil {
						log.Fatal("binary.Read failed", err)
					}
					limite := s.Tabla[i].Start + int64(unsafe.Sizeof(ebr{})) + s.Tabla[i].Size
					if &ebrTemp != nil {
						//fmt.Println(ebrTemp.Start)
						for j := s.Tabla[i].Start; j < limite; j++ {
							ebrLeido := ebr{}
							file.Seek(j, 0)
							data1 := readNextBytes(file, unsafe.Sizeof(ebr{}))
							buffer1 := bytes.NewBuffer(data1)
							err = binary.Read(buffer1, binary.BigEndian, &ebrLeido)
							if err != nil {
								log.Fatal("binary.Read failed", err)
							}
							if &ebrLeido != nil {
								/*fmt.Println("Ebr leido")
								fmt.Println(ebrLeido.Name)
								fmt.Println("Name que viene")
								fmt.Println(nombre)
								fmt.Println("ebr leido next")
								fmt.Println(ebrLeido.Next)*/
								if ebrLeido.Next != -1 && ebrLeido.Size == 0 {
									//Aqui no valuo porque es el primer EBR solo lo salto
									j = ebrLeido.Next - 1
								} else if ebrLeido.Next == -1 && ebrLeido.Size == 0 {
									//la particion esta vacia
									j = limite + 1
									return 0, 0, 0
								} else if ebrLeido.Next == -1 && ebrLeido.Size > 0 { //lego al utimo ebr
									//valuo con el limite porque es el ultimo ebr
									if ebrLeido.Name == nombre {
										return 2, ebrLeido.Start, ebrLeido.Size
									}
									return 0, 0, 0
								} else if ebrLeido.Next != -1 && ebrLeido.Size > 0 { //esta en los ebr antes del ultimo
									//verificar pero con el next
									//fmt.Println(ebrLeido.Name)
									//fmt.Println(nombre)
									if ebrLeido.Name == nombre {

										return 2, ebrLeido.Start, ebrLeido.Size
									}
									j = ebrLeido.Next - 1
								}
							}
						}
					}
				}
			}
		}
	}
	return 0, 0, 0
}

func realizarFormato(path string, inicioParticion int64, tamParticion int64, tipo string) int {

	//datos para formula
	var tamAVD int64 = int64(unsafe.Sizeof(avd{}))
	var tamDD int64 = int64(unsafe.Sizeof(dd{}))
	var tamInodo int64 = int64(unsafe.Sizeof(inodo{}))
	var tamBloque int64 = int64(unsafe.Sizeof(bloque{}))
	var tamBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
	var tamSuperBloque int64 = int64(unsafe.Sizeof(superbloque{}))

	//formula
	var nEstructuras = (tamParticion - (2 * tamSuperBloque)) / (27 + tamAVD + tamDD + (5*tamInodo + (20 * tamBloque) + tamBitacora))

	if nEstructuras <= 0 {
		fmt.Print("RESULTADO: Espacio insuficiente para formatear la particion")
		return 0
	}
	//cantidad de tipos de estructura
	cantidadAVD := nEstructuras
	cantidadDD := nEstructuras
	cantidadInodos := 5 * nEstructuras
	cantidadBloques := 20 * nEstructuras
	cantidadBitacoras := nEstructuras

	inicioSuperBloque := inicioParticion
	iniciobitmapAVD := inicioParticion + tamSuperBloque
	inicioAVD := iniciobitmapAVD + cantidadAVD
	iniciobitmapDD := inicioAVD + (tamAVD * cantidadAVD)
	inicioDD := iniciobitmapDD + cantidadDD
	iniciobitMapInodo := inicioDD + (tamDD * cantidadDD)
	inicioinodos := iniciobitMapInodo + cantidadInodos
	iniciobitmapBloque := inicioinodos + (tamInodo * cantidadInodos)
	iniciobloques := iniciobitmapBloque + cantidadBloques
	iniciobitacora := iniciobloques + (tamBloque * cantidadBloques)
	iniciocopiaSB := iniciobitacora + (tamBitacora * cantidadBitacoras)
	//finalparticion := iniciocopiaSB + tamSuperBloque

	/*
		fmt.Print("No. de estructuras: ")
		fmt.Println(nEstructuras)

		fmt.Print("Inicio de SuperBloque: ")
		fmt.Println(inicioSuperBloque)

		fmt.Print("Inicio de bitmap AVD: ")
		fmt.Println(iniciobitmapAVD)

		fmt.Print("Inicio de AVD: ")
		fmt.Println(inicioAVD)

		fmt.Print("Inicio de bitmap DD: ")
		fmt.Println(iniciobitmapDD)

		fmt.Print("Inicio de DD: ")
		fmt.Println(inicioDD)

		fmt.Print("Inicio de bitmap Inodos: ")
		fmt.Println(iniciobitMapInodo)

		fmt.Print("Inicio de Inodos: ")
		fmt.Println(inicioinodos)

		fmt.Print("Inicio de bitmap Bloques: ")
		fmt.Println(iniciobitmapBloque)

		fmt.Print("Inicio de bloques: ")
		fmt.Println(iniciobloques)

		fmt.Print("Inicio de bitacora: ")
		fmt.Println(iniciobitacora)

		fmt.Print("Inicio copia de SB: ")
		fmt.Println(iniciocopiaSB)

		fmt.Print("Final de particion: ")
		fmt.Println(finalparticion)
	*/
	//abriendo el archivo
	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	//formando el superbloque
	nuevoSB := superbloque{}
	/*Colocando nombre al superbloque*/
	arregloNombre := strings.Split(path, "/")
	nombreEncontrado := ""
	for i := 0; i < len(arregloNombre); i++ {
		if strings.Contains(strings.ToLower(arregloNombre[i]), ".dsk") {
			nombreEncontrado = arregloNombre[i]
			break
		}
	}
	nombreEncontrado = strings.TrimSpace(nombreEncontrado)
	if strings.Compare(nombreEncontrado, "") != 0 {
		if len(arregloNombre) < 101 {
			copy(nuevoSB.NombreHD[:], nombreEncontrado)
		}
	}
	nuevoSB.ArbolVirtualCount = cantidadAVD
	nuevoSB.DetalleDirectorioCount = cantidadDD
	nuevoSB.InodosCount = cantidadInodos
	nuevoSB.BloquesCount = cantidadBloques

	nuevoSB.ArbolVirtualFree = cantidadAVD - 1     //porque se crea la carpeta raiz
	nuevoSB.DetalleDirectorioFree = cantidadDD - 1 //porque se crea la carpeta raiz
	nuevoSB.InodosFree = cantidadInodos
	nuevoSB.BloquesFree = cantidadBloques
	nuevoSB.DateCreacion = getFechaHora()
	nuevoSB.DateUltimoMontaje = nuevoSB.DateCreacion
	nuevoSB.MontajesCount = 1
	nuevoSB.InicioBMAV = iniciobitmapAVD
	nuevoSB.InicioAV = inicioAVD
	nuevoSB.InicioBMDD = iniciobitmapDD
	nuevoSB.InicioDD = inicioDD
	nuevoSB.InicioBMInodos = iniciobitMapInodo
	nuevoSB.InicioInodos = inicioinodos
	nuevoSB.InicioBMBloques = iniciobitmapBloque
	nuevoSB.InicioBloques = iniciobloques
	nuevoSB.InicioLog = iniciobitacora
	nuevoSB.TamAV = tamAVD
	nuevoSB.TamDD = tamDD
	nuevoSB.TamInodo = tamInodo
	nuevoSB.TamBloque = tamBloque
	nuevoSB.PrimerLibreAV =2 //porque se crea la carpeta raiz
	nuevoSB.PrimerLibreDD =2 //porque se crea la carpeta raiz
	nuevoSB.PrimerLibreInodo = 1
	nuevoSB.PrimerLibreBloque = 1
	nuevoSB.MagicNum = 201314821

	pos := inicioParticion

	/*Escribiendo el superbloque*/
	pos = inicioSuperBloque
	file.Seek(pos, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, nuevoSB)
	writeNextBytes(file, binario.Bytes())

	/*Escribiendo bitmap de avd*/
	pos = iniciobitmapAVD
	for i := 0; i < int(cantidadAVD); i++ {
		file.Seek(pos, 0)
		var vacio byte = 0
		if i == 0 {
			vacio = 1
		}
		s := &vacio
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())
		pos = pos + int64(unsafe.Sizeof(vacio))
	}
	/*Escribiendo bitmap de dd*/
	pos = iniciobitmapDD
	for i := 0; i < int(cantidadDD); i++ {
		file.Seek(pos, 0)
		var vacio byte = 0
		if i == 0 {
			vacio = 1
		}
		s := &vacio
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())
		pos = pos + int64(unsafe.Sizeof(vacio))
	}

	/*Escribiendo bitmap de inodos*/
	pos = iniciobitMapInodo
	for i := 0; i < int(cantidadInodos); i++ {
		file.Seek(pos, 0)
		var vacio byte = 0
		s := &vacio
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())
		pos = pos + int64(unsafe.Sizeof(vacio))
	}

	/*Escribiendo bitmap de bloques*/
	pos = iniciobitmapBloque
	for i := 0; i < int(cantidadBloques); i++ {
		file.Seek(pos, 0)
		var vacio byte = 0
		s := &vacio
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())
		pos = pos + int64(unsafe.Sizeof(vacio))
	}

	/*Escribiendo bitacora*/
	pos = iniciobitacora
	for i := 0; i < int(cantidadBitacoras); i++ {
		file.Seek(pos, 0)
		bitacoraVacio := bitacora{}
		s := &bitacoraVacio
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		writeNextBytes(file, binario2.Bytes())
		pos = pos + int64(unsafe.Sizeof(bitacora{}))
	}

	/*Escribiendo copia del SB*/
	pos = iniciocopiaSB
	file.Seek(pos, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, &nuevoSB)
	writeNextBytes(file, binario2.Bytes())

	if strings.Compare(tipo, "full") == 0 || strings.Compare(tipo, "fast") == 0 {
		/*Escribiendo avd*/
		pos = inicioAVD
		for i := 0; i < int(cantidadAVD); i++ {
			file.Seek(pos, 0)
			avdVacio := avd{}
			if i == 0 {
				avdVacio.AVDApArbolVirtualDirectorio = 0
				avdVacio.AVDApDetalleDirectorio = inicioDD
				avdVacio.AVDFechaCreacion = nuevoSB.DateCreacion
				avdVacio.AVDNombreDirectorio[0] = 47
			}
			s := &avdVacio
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s)
			writeNextBytes(file, binario2.Bytes())
			pos = pos + int64(unsafe.Sizeof(avd{}))
		}

		/*Escribiendo dd*/
		pos = inicioDD
		for i := 0; i < int(cantidadDD); i++ {
			file.Seek(pos, 0)
			ddVacio := dd{}
			s := &ddVacio
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s)
			writeNextBytes(file, binario2.Bytes())
			pos = pos + int64(unsafe.Sizeof(dd{}))
		}

		/*Escribiendo inodos*/
		pos = inicioinodos
		for i := 0; i < int(cantidadInodos); i++ {
			file.Seek(pos, 0)
			inodoVacio := inodo{}
			s := &inodoVacio
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s)
			writeNextBytes(file, binario2.Bytes())
			pos = pos + int64(unsafe.Sizeof(inodo{}))
		}

		/*Escribiendo bloques*/
		pos = iniciobloques
		for i := 0; i < int(cantidadBloques); i++ {
			file.Seek(pos, 0)
			bloqueVacio := bloque{}
			s := &bloqueVacio
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s)
			writeNextBytes(file, binario2.Bytes())
			pos = pos + int64(unsafe.Sizeof(bloque{}))
		}
	}
	/*Crear la carpeta raiz*/
	//graficarSB(path, inicioParticion)
	return 1
}

func getFechaHora() [16]byte {
	var retFecha [16]byte
	fechahora := time.Now()
	fechahoraArreglo := strings.Split(fechahora.String(), "")
	fechahoraCadena := ""
	for i := 0; i < 16; i++ {
		fechahoraCadena = fechahoraCadena + fechahoraArreglo[i]
	}
	copy(retFecha[:], fechahoraCadena)
	return retFecha
}

/**************************************************************
	CREACION DE DIRECTORIOS
***************************************************************/
func crearDirectorio(id string, pathCadena string, atributoP int) {
	var tipoParticion int
	var inicioParticion int64
	if strings.Compare(id, "") != 0 {
		s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
		if len(s) > 3 {
			if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
				var letra byte = s[2][0]
				inputFmt := id[3:len(id)] + ""
				idParticion := atributoSize(inputFmt)
				if idParticion > 0 {
					numResult, nombre, path := desmontarParticion2(letra, int64(idParticion))
					//fmt.Println("El nombre que retorna desmontar particion")
					//fmt.Println(nombre)
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
						//Aqui ya tengo el path del disco y el nombre de la particion en la que se va a crear
						//fmt.Println(pathEnviar)
						//fmt.Println(nombre)

						/************************************************************************/
						tipoParticion, inicioParticion, _ = obtenerDatosParticion(pathEnviar, *nombre)
						if tipoParticion == 0 {
							fmt.Println("RESULTADO: No se encuentra la particion")
						} else if tipoParticion == 3 {
							fmt.Println("RESULTADO: No se puede formatear una particion extendida")
						} else if tipoParticion == 1 {
							estadoActual := getEstadoFormato(letra, int64(idParticion))
							if estadoActual == 5 {
								//hay un problema
								fmt.Println("RESULTADO: Existe un problema de formateo de la particion")
							} else if estadoActual == 0 {
								//no esta formateada
								fmt.Println("RESULTADO: La particion en la que intenta crear el directorio no se encuentra formateada")
								return
							} else if estadoActual == 1 {
								//ya esta formateada
								//mandar a crear la carpeta
								crearCarpeta(pathEnviar, inicioParticion, pathCadena, atributoP)
							}
						} else if tipoParticion == 2 {
							//es una extendida
							estadoActual := getEstadoFormato(letra, int64(idParticion))
							if estadoActual == 5 {
								//hay un problema
								fmt.Println("RESULTADO: Existe un problema de formateo de la particion")
							} else if estadoActual == 0 {
								//no esta formateada
								fmt.Println("RESULTADO: La particion en la que intenta crear el directorio no se encuentra formateada")
								return
							} else if estadoActual == 1 {
								//ya esta formateada
								crearCarpeta(pathEnviar, inicioParticion+int64(unsafe.Sizeof(ebr{})), pathCadena, atributoP)
							}
						}
						/**************************************************/
					} else {
						fmt.Println("RESULTADO: Error en el id de la particion en la que se desea crear el directorio")
					}
				} else {
					fmt.Println("RESULTADO: Error en el id de la particion en la que se desea crear el directorio")
				}
			} else {
				fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto, no se creara el directorio ")
			}
		}
	}

}

func crearCarpeta(pathDisco string, inicioSuperBloque int64, pathCrear string, atributoP int) {
	file, err := os.OpenFile(strings.ReplaceAll(pathDisco, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	sbLeido := superbloque{}
	file.Seek(inicioSuperBloque, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		return
	}

	if sbLeido.MagicNum == 201314821 {
		apRaiz := sbLeido.InicioAV
		dirActual := apRaiz
		//aqui ya tengo mi superbloque leido
		s := strings.Split(pathCrear, "/")
		if len(s) > 0 {
			for i := 0 ; i< len(s) ; i++{
				if strings.Compare(s[i] , "")!= 0{
					if i == len(s) - 1{
						index1, index2 := crearcarpetaRecursivo(pathDisco, s[i], inicioSuperBloque, dirActual, atributoP, 1)
						if index1 == 1{
							dirActual = index2
						}else{
							return
						}
					}else{
						index1, index2 := crearcarpetaRecursivo(pathDisco, s[i], inicioSuperBloque, dirActual, atributoP, 0)
						if index1 == 1{
							dirActual = index2
						}else{
							return
						}
					}
				}
			}
		}

	}
}

//el int retorna si el directorio existe o fue creado
//el int64 retorna la estructura nueva o ya existente
func crearcarpetaRecursivo(pathDisco string, carpeta string, inicioSuperBloque int64, avdActual int64, atributoP int, fin int) (int, int64){
	file, err := os.OpenFile(strings.ReplaceAll(pathDisco, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return 0,0
	}
	sbLeido := superbloque{}
	file.Seek(inicioSuperBloque, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		return 0,0
	}

	if sbLeido.MagicNum == 201314821 {
		//moverse a avdActual y leerlo
		avdLeido := avd{}
		file.Seek(avdActual, 0)
		data := readNextBytes(file, unsafe.Sizeof(avd{}))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &avdLeido)
		if err != nil {
			log.Fatal("binary.Read failed", err)
			return 0,0
		}
		//revisar entre el arreglo actual  
		var convertido [16]byte
		copy(convertido[:], strings.ToLower(strings.TrimSpace(carpeta)))

		if avdLeido.AVDNombreDirectorio == convertido{
			return 1, avdActual
		}
		contadorVacio := 0
		for i := 0; i < len(avdLeido.AVDApArraySubdirectorios); i++ {
			
			if avdLeido.AVDApArraySubdirectorios[i] != 0{
				avdNuevo := avd{}
				file.Seek(avdLeido.AVDApArraySubdirectorios[i], 0)
				data := readNextBytes(file, unsafe.Sizeof(avd{}))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &avdNuevo)
				if err != nil {
					log.Fatal("binary.Read failed", err)
					return 0,0
				}
				if avdNuevo.AVDNombreDirectorio == convertido {
					if fin == 1{
						fmt.Println("RESULTADO: El directorio ya se encuentra creado")
					}
					return 1, avdLeido.AVDApArraySubdirectorios[i]
				}

			}else {
				contadorVacio = contadorVacio + 1
			}	
		}

		if atributoP == 0 && fin == 0 && avdLeido.AVDApArbolVirtualDirectorio == 0{
			fmt.Println("RESULTADO: No se ha encontrado el directorio " + carpeta)
			return 0,0
		}else if avdLeido.AVDApArbolVirtualDirectorio != 0{
			return crearcarpetaRecursivo(pathDisco, carpeta, inicioSuperBloque, avdLeido.AVDApArbolVirtualDirectorio, atributoP, fin)
		}
		
		//si sale es porque no encontro el nombre del directorio
		if contadorVacio > 0 {
			//todavia hay directorios disponibles en el array actual 

			/*Verificar si hay disponibles*/

			if sbLeido.ArbolVirtualFree > 0 && sbLeido.DetalleDirectorioFree > 0 {
				//si lo hay
				ingresoAVD := avd{}
				ingresoAVD.AVDFechaCreacion = getFechaHora()
				ingresoAVD.AVDNombreDirectorio = convertido
				
				ingresoDD := dd{} //no lleva parametros 


				//2 escribir el dd  nuevo y modificar el bitmap
				posactual := sbLeido.InicioDD + (sbLeido.PrimerLibreDD - 1)*int64(unsafe.Sizeof(dd{})) //posicion para escribir el avd
				ingresoAVD.AVDApDetalleDirectorio = posactual
				file.Seek(posactual, 0 )
				var binario2 bytes.Buffer
				binary.Write(&binario2, binary.BigEndian, ingresoDD)
				writeNextBytes(file, binario2.Bytes())
				
				posactual = sbLeido.InicioBMDD + sbLeido.PrimerLibreDD -1
				file.Seek(posactual, 0)
				var cambioBitmap2 byte = 1
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, cambioBitmap2)
				writeNextBytes(file, binario3.Bytes())

				//1 escribir el avd nuevo y modificar el bitmap
				posactual = sbLeido.InicioAV + ((sbLeido.PrimerLibreAV-1)*int64(unsafe.Sizeof(avd{})))
				retorno := posactual //posicion para escribir el avd
				file.Seek(posactual, 0 )
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, ingresoAVD)
				writeNextBytes(file, binario.Bytes())
				
				posactual = sbLeido.InicioBMAV + sbLeido.PrimerLibreAV -1
				file.Seek(posactual, 0)
				var cambioBitmap byte = 1
				var binario1 bytes.Buffer
				binary.Write(&binario1, binary.BigEndian, cambioBitmap)
				writeNextBytes(file, binario1.Bytes())
				//meter el arbol actual a su directorio padre 
				for i:=0 ; i < len(avdLeido.AVDApArraySubdirectorios); i++{
					if avdLeido.AVDApArraySubdirectorios[i] == 0{
						avdLeido.AVDApArraySubdirectorios[i] = retorno
						break
					}
				}
				posactual = avdActual
				file.Seek(avdActual, 0)
				var binario10 bytes.Buffer

				binary.Write(&binario10, binary.BigEndian, avdLeido)
				writeNextBytes(file, binario10.Bytes())
				//3 modificar el superbloque
				sbLeido.ArbolVirtualCount = sbLeido.ArbolVirtualCount + 1			//se suma 1 directorio
				sbLeido.DetalleDirectorioCount = sbLeido.DetalleDirectorioCount + 1	//se suma un detalle 
				sbLeido.ArbolVirtualFree = sbLeido.ArbolVirtualFree -1 				//un avd libre menos
				sbLeido.DetalleDirectorioFree = sbLeido.DetalleDirectorioFree -1 	//un dd  libre menos 
				sbLeido.PrimerLibreAV = sbLeido.PrimerLibreAV + 1					//avanza uno en el bitmap de avd
				sbLeido.PrimerLibreDD = sbLeido.PrimerLibreDD + 1					//avanza uno en el bitmap de dd
				//4 escribir el superbloque
				posactual = inicioSuperBloque
				file.Seek(posactual, 0)
				var binario5 bytes.Buffer
				binary.Write(&binario5, binary.BigEndian, sbLeido)
				writeNextBytes(file, binario5.Bytes())
				fmt.Println("RESULTADO: Se ha creado el directorio " + carpeta)
				return 1, retorno

			}else{
				//no hay espacio 
				fmt.Println("RESULTADO: No hay espacio disponible para crear el directorio")
				return 0,0
			}
		}

		//si no se debe crear un nuevo avd con el nombre del actual para meter el subdirectorio 
		//lo repito aqui mismo 
		if sbLeido.ArbolVirtualFree > 0 && sbLeido.DetalleDirectorioFree > 0 {
				//si lo hay
				ingresoAVD := avd{}
				ingresoAVD.AVDFechaCreacion = avdLeido.AVDFechaCreacion
				ingresoAVD.AVDNombreDirectorio = avdLeido.AVDNombreDirectorio
				
				ingresoDD := dd{} //no lleva parametros 


				//2 escribir el dd  nuevo y modificar el bitmap
				posactual := sbLeido.InicioDD + (sbLeido.PrimerLibreDD - 1)*int64(unsafe.Sizeof(dd{})) //posicion para escribir el avd
				ingresoAVD.AVDApDetalleDirectorio = posactual
				file.Seek(posactual, 0 )
				var binario2 bytes.Buffer
				binary.Write(&binario2, binary.BigEndian, ingresoDD)
				writeNextBytes(file, binario2.Bytes())
				
				posactual = sbLeido.InicioBMDD + sbLeido.PrimerLibreDD -1
				file.Seek(posactual, 0)
				var cambioBitmap2 byte = 1
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, cambioBitmap2)
				writeNextBytes(file, binario3.Bytes())

				//1 escribir el avd nuevo y modificar el bitmap
				posactual = sbLeido.InicioAV + ((sbLeido.PrimerLibreAV-1)*int64(unsafe.Sizeof(avd{})))
				avdLeido.AVDApArbolVirtualDirectorio = posactual
				retorno2 := posactual //posicion para escribir el avd
				file.Seek(posactual, 0 )
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, ingresoAVD)
				writeNextBytes(file, binario.Bytes())
				
				posactual = sbLeido.InicioBMAV + sbLeido.PrimerLibreAV -1
				file.Seek(posactual, 0)
				var cambioBitmap byte = 1
				var binario1 bytes.Buffer
				binary.Write(&binario1, binary.BigEndian, cambioBitmap)
				writeNextBytes(file, binario1.Bytes())
				//meter el arbol actual a su directorio padre 
				for i:=0 ; i < len(avdLeido.AVDApArraySubdirectorios); i++{
					if avdLeido.AVDApArraySubdirectorios[i] == 0{
						avdLeido.AVDApArraySubdirectorios[i] = retorno2
						i = 100
					}
				}
				posactual = avdActual
				file.Seek(posactual, 0)
				var binario10 bytes.Buffer
				binary.Write(&binario10, binary.BigEndian, avdLeido)
				writeNextBytes(file, binario10.Bytes())
				//3 modificar el superbloque
				sbLeido.ArbolVirtualCount = sbLeido.ArbolVirtualCount + 1			//se suma 1 directorio
				sbLeido.DetalleDirectorioCount = sbLeido.DetalleDirectorioCount + 1	//se suma un detalle 
				sbLeido.ArbolVirtualFree = sbLeido.ArbolVirtualFree -1 				//un avd libre menos
				sbLeido.DetalleDirectorioFree = sbLeido.DetalleDirectorioFree -1 	//un dd  libre menos 
				sbLeido.PrimerLibreAV = sbLeido.PrimerLibreAV + 1					//avanza uno en el bitmap de avd
				sbLeido.PrimerLibreDD = sbLeido.PrimerLibreDD + 1					//avanza uno en el bitmap de dd
				//4 escribir el superbloque
				posactual = inicioSuperBloque
				file.Seek(posactual, 0)
				var binario5 bytes.Buffer
				binary.Write(&binario5, binary.BigEndian, sbLeido)
				writeNextBytes(file, binario5.Bytes())
				return crearcarpetaRecursivo(pathDisco, carpeta, inicioSuperBloque, retorno2, atributoP, fin)
			}else{
				//no hay espacio 
				fmt.Println("RESULTADO: No hay espacio disponible para crear el directorio")
				return 0,0
			}

		//*****************************************************************************************

	}
	return 0,0
}

/**************************************************************
	CREACION DE ARCHIVOS
***************************************************************/
func crearFile(id string, pathCadena string, atributoP int, tamano int64, contenido string){
	s := strings.Split(pathCadena, "/")
	soloPath := "/"
	nameArchivo:=""
	if len(s) > 0{
		for i:=0 ; i< len(s); i++{
			if i != len(s) - 1{
				soloPath= soloPath + "/" + s[i]
			}else if i == len(s) -1{
				nameArchivo = strings.TrimSpace(s[i])
				break
			}
		}
	}

	res, dir, pathO, inicioO:= crearDirectorioFILE(id, soloPath, atributoP) 
	if res==1{
		//crear el archivo
		i :=creacionArchivo(pathO, inicioO, dir, tamano, nameArchivo, contenido)
		if i == 3{
			fmt.Println("RESULTADO: El Archivo se encuentra repetido")
		}else if i == 0{
			fmt.Println("RESULTADO: No se ha podido crear el archivo")
		}else if i == 1{
			fmt.Println("RESULTADO: El archivo fue creado exitosamente")
		}
		
	} 
	//si no, ya mando el mensaje anteriormente
}

func creacionArchivo(pathDisco string, inicioParticion int64, inicioDirectorio int64, tamano int64, nombreArchivo string, contenidoArchivo string)int{
	file, err := os.OpenFile(strings.ReplaceAll(pathDisco, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	sbLeido := superbloque{}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		return 0
	}

	if sbLeido.MagicNum == 201314821 {
		avdLeido := avd{}
		file.Seek(inicioDirectorio, 0)
		data := readNextBytes(file, unsafe.Sizeof(avd{}))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &avdLeido)
		if err != nil {
			log.Fatal("binary.Read failed", err)
			return 0
		}
		if avdLeido.AVDApDetalleDirectorio != 0{
			//aqui buscar en los detalles de directorio en busca de nombres repetidos 
			i:= insertarArchivo(pathDisco, inicioParticion, avdLeido.AVDApDetalleDirectorio, tamano, nombreArchivo, contenidoArchivo)
			return i
		}
	}
	return 0
}


func insertarArchivo(pathDisco string, inicioParticion int64, inicioDetalle int64, tamano int64, nombreArchivo string, contenidoArchivo string) int {
	file, err := os.OpenFile(strings.ReplaceAll(pathDisco, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	sbLeido := superbloque{}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		return 0
	}

	if sbLeido.MagicNum == 201314821 {
		ddLeido := dd{}
		file.Seek(inicioDetalle, 0)
		data := readNextBytes(file, unsafe.Sizeof(dd{}))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &ddLeido)
		if err != nil {
			log.Fatal("binary.Read failed", err)
			return 0
		}
		contadorVacio:= 0
		if &ddLeido != nil{
			//aqui buscar en los detalles de directorio en busca de nombres repetidos 
			var convertido [16]byte
			copy(convertido[:], strings.ToLower(strings.TrimSpace(nombreArchivo)))
			for i :=0; i < len(ddLeido.DDArrayFiles) ; i++{

				if ddLeido.DDArrayFiles[i].FileApInodo != 0{
					if ddLeido.DDArrayFiles[i].FileNombre == convertido{
						return 3
					}
				}else {
					contadorVacio = contadorVacio + 1
				}
				
			}
			//si sale es porque no lo encontro 
			if ddLeido.DDApDetalleDirectorio == 0 {
				//crearlo aqui
				if contadorVacio == 0{
					//crear un nuevo dd y mandarlo a insertar en el nuevo DD \
					if sbLeido.DetalleDirectorioFree > 0{
						
						ddExtra := dd{}
						direccionDD := sbLeido.InicioBMDD + (sbLeido.PrimerLibreDD - 1)*int64(unsafe.Sizeof(dd{}))
						ddLeido.DDApDetalleDirectorio = direccionDD
						file.Seek(direccionDD,0)
						var binarioDDE bytes.Buffer
						binary.Write(&binarioDDE, binary.BigEndian, ddExtra)
						writeNextBytes(file, binarioDDE.Bytes())

						file.Seek(inicioDetalle,0)
						var binarioDDE1 bytes.Buffer
						binary.Write(&binarioDDE1, binary.BigEndian, ddLeido)
						writeNextBytes(file, binarioDDE1.Bytes())


						sbLeido.DetalleDirectorioFree = sbLeido.DetalleDirectorioFree - 1
						sbLeido.PrimerLibreDD = sbLeido.PrimerLibreDD + 1
						file.Seek(inicioParticion, 0)
						var binarioSBE bytes.Buffer
						binary.Write(&binarioSBE, binary.BigEndian, sbLeido)
						writeNextBytes(file, binarioSBE.Bytes())

						return insertarArchivo(pathDisco, inicioParticion, direccionDD, tamano, nombreArchivo, contenidoArchivo)
					}
				}else{
					//crearlo aqui
					//calcular tamano 
					numBloque:= 0
					if tamano <= 25 {
						numBloque = 1
					}else{
						residuo := math.Mod(float64(tamano), 25)
						numBloque = int(tamano / 25)
						if residuo != 0{
							numBloque = numBloque + 1
						}
						
					}

					numInodos:= 0
					if numBloque <= 4 {
						numInodos = 1
					}else{
						residuo := math.Mod(float64(numBloque), 4)
						numInodos = int(numBloque / 4)
						if residuo != 0{
							numInodos = numInodos + 1
						}
						
					}

					//preguntando si hay inodos y bloques suficientes 
					if sbLeido.InodosFree >= int64(numInodos) && sbLeido.BloquesFree >= int64(numBloque){
						//si hay espacio 

						archivoNuevo := archivo{}
						archivoNuevo.FileNombre = convertido
						archivoNuevo.FileDateCreacion = getFechaHora()
						archivoNuevo.FileDateModificacion = getFechaHora()
						archivoNuevo.FileApInodo =  sbLeido.InicioBMInodos + (sbLeido.PrimerLibreInodo - 1)*int64(unsafe.Sizeof(inodo{}))

						/***************************************************************************/
						//crear bloques 
						contadorBloques := 0
						tmpInodo := inodo{}
						tmpInodo.ICountInodo = sbLeido.InicioBMInodos
						tmpInodo.ISizeArchivo = tamano
						tmpInodo.ICountBloquesAsignados = 0
						escrito := 0
						for i := 0 ; i < numBloque ; i++{
							if contadorBloques == 3{
								//crear otro inodo
								direccionBloque :=  sbLeido.InicioBloques + (sbLeido.PrimerLibreBloque - 1)*int64(unsafe.Sizeof(bloque{}))
								bloqueNuevo := bloque{}
								/*Aqui le tendria que poner el contenido peor aun no lo hago*/
								file.Seek(direccionBloque,0)
								var binario bytes.Buffer
								binary.Write(&binario, binary.BigEndian, bloqueNuevo)
								writeNextBytes(file, binario.Bytes())
								tmpInodo.IArrayBloques[contadorBloques] = direccionBloque 
								tmpInodo.ICountBloquesAsignados = tmpInodo.ICountBloquesAsignados + 1
								sbLeido.PrimerLibreBloque = sbLeido.PrimerLibreBloque + 1
								/*tengo que escribir el inodo actual en el disco*/
								direccionInodo :=  sbLeido.InicioBMInodos + (sbLeido.PrimerLibreInodo - 1)*int64(unsafe.Sizeof(inodo{}))
								if i != numBloque -1 {
									//no lo estoy creando
									tmpInodo.IApIndirecto = direccionInodo + int64(unsafe.Sizeof(inodo{}))
								}
								file.Seek(direccionInodo,0)
								var binario2 bytes.Buffer
								binary.Write(&binario2, binary.BigEndian, tmpInodo)
								writeNextBytes(file, binario2.Bytes()) 
								sbLeido.PrimerLibreInodo = sbLeido.PrimerLibreInodo + 1
								escrito = 1
								nuevoInodo := inodo{}
								nuevoInodo.ICountInodo = sbLeido.InicioBMInodos
								nuevoInodo.ISizeArchivo = tamano
								nuevoInodo.ICountBloquesAsignados = 0
								tmpInodo = nuevoInodo
								contadorBloques = 0
							}else{
								escrito = 0
								//solo insertar bloque
								direccionBloque :=  sbLeido.InicioBloques + (sbLeido.PrimerLibreBloque - 1)*int64(unsafe.Sizeof(bloque{}))
								bloqueNuevo := bloque{}
								/*Aqui le tendria que poner el contenido peor aun no lo hago*/
								file.Seek(direccionBloque,0)
								var binario bytes.Buffer
								binary.Write(&binario, binary.BigEndian, bloqueNuevo)
								writeNextBytes(file, binario.Bytes())
								tmpInodo.IArrayBloques[contadorBloques] = direccionBloque 
								tmpInodo.ICountBloquesAsignados = tmpInodo.ICountBloquesAsignados + 1
								sbLeido.PrimerLibreBloque = sbLeido.PrimerLibreBloque + 1
								contadorBloques = contadorBloques + 1
							}
						}

						if escrito == 0{
							direccionInodo :=  sbLeido.InicioBMInodos + (sbLeido.PrimerLibreInodo - 1)*int64(unsafe.Sizeof(inodo{}))
							file.Seek(direccionInodo,0)
							var binario2 bytes.Buffer
							binary.Write(&binario2, binary.BigEndian, tmpInodo)
							writeNextBytes(file, binario2.Bytes()) 
							sbLeido.PrimerLibreInodo = sbLeido.PrimerLibreInodo + 1
							escrito = 1
						}
						/***************************************************************************/

						//escribiendo el archivo   
						for i := 0 ; i < len(ddLeido.DDArrayFiles); i++{
							if ddLeido.DDArrayFiles[i].FileApInodo == 0{
								ddLeido.DDArrayFiles[i] = archivoNuevo
							}
						}

						//ya que ingrese el archivo en el dd rescribo el dd
						file.Seek(inicioDetalle, 0)
						var binarioDD bytes.Buffer
						binary.Write(&binarioDD, binary.BigEndian, ddLeido)
						writeNextBytes(file, binarioDD.Bytes())

						//actualizar el superbloque 
						file.Seek(inicioParticion, 0)
						var binarioSB bytes.Buffer
						binary.Write(&binarioSB, binary.BigEndian, sbLeido)
						writeNextBytes(file, binarioSB.Bytes())
						return 1
					}else{
						//no hay suficientes inodos o bloques 
						return 0
					}
				}
			}else{
				return insertarArchivo(pathDisco, inicioParticion, ddLeido.DDApDetalleDirectorio, tamano, nombreArchivo, contenidoArchivo)
			}

		}
	}
	return 0
}

func crearDirectorioFILE(id string, pathCadena string, atributoP int) (int, int64, string, int64){
	var tipoParticion int
	var inicioParticion int64
	if strings.Compare(id, "") != 0 {
		s := strings.Split(strings.ToLower(strings.TrimSpace(id)), "")
		if len(s) > 3 {
			if s[0][0] == 'v' && s[1][0] == 'd' && s[2][0] > 96 && s[2][0] < 123 {
				var letra byte = s[2][0]
				inputFmt := id[3:len(id)] + ""
				idParticion := atributoSize(inputFmt)
				if idParticion > 0 {
					numResult, nombre, path := desmontarParticion2(letra, int64(idParticion))
					//fmt.Println("El nombre que retorna desmontar particion")
					//fmt.Println(nombre)
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
						//Aqui ya tengo el path del disco y el nombre de la particion en la que se va a crear
						//fmt.Println(pathEnviar)
						//fmt.Println(nombre)

						/************************************************************************/
						tipoParticion, inicioParticion, _ = obtenerDatosParticion(pathEnviar, *nombre)
						if tipoParticion == 0 {
							//fmt.Println("RESULTADO: No se encuentra la particion")
							return 0,0,"",0
						} else if tipoParticion == 3 {
							//fmt.Println("RESULTADO: No se puede formatear una particion extendida")
							return 0,0,"",0
						} else if tipoParticion == 1 {
							estadoActual := getEstadoFormato(letra, int64(idParticion))
							if estadoActual == 5 {
								//hay un problema
								//fmt.Println("RESULTADO: Existe un problema de formateo de la particion")
								return 0,0,"",0
							} else if estadoActual == 0 {
								//no esta formateada
								//fmt.Println("RESULTADO: La particion en la que intenta crear el directorio no se encuentra formateada")
								return 0,0,"",0
							} else if estadoActual == 1 {
								//ya esta formateada
								//mandar a crear la carpeta
								v1, v2 := crearCarpetaFILE(pathEnviar, inicioParticion, pathCadena, atributoP)
								return v1,v2, pathEnviar, inicioParticion
							}
						} else if tipoParticion == 2 {
							//es una extendida
							estadoActual := getEstadoFormato(letra, int64(idParticion))
							if estadoActual == 5 {
								//hay un problema
								//fmt.Println("RESULTADO: Existe un problema de formateo de la particion")
								return 0,0,"",0
							} else if estadoActual == 0 {
								//no esta formateada
								//fmt.Println("RESULTADO: La particion en la que intenta crear el directorio no se encuentra formateada")
								return 0,0,"",0
							} else if estadoActual == 1 {
								//ya esta formateada
								v1, v2 := crearCarpetaFILE(pathEnviar, inicioParticion+int64(unsafe.Sizeof(ebr{})), pathCadena, atributoP)
								return v1,v2, pathEnviar, inicioParticion
							}
						}
						/**************************************************/
					} else {
						//fmt.Println("RESULTADO: Error en el id de la particion en la que se desea crear el directorio")
						return 0,0, "",0
					}
				} else {
					//fmt.Println("RESULTADO: Error en el id de la particion en la que se desea crear el directorio")
					return 0,0,"",0
				}
			} else {
				//fmt.Println("RESULTADO: El formato del ID de la particion es incorrecto, no se creara el directorio ")
				return 0,0,"",0
			}
		}
	}
	return 0,0,"",0
}

func crearCarpetaFILE(pathDisco string, inicioSuperBloque int64, pathCrear string, atributoP int) (int, int64) {
	file, err := os.OpenFile(strings.ReplaceAll(pathDisco, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return 0 ,0
	}
	sbLeido := superbloque{}
	file.Seek(inicioSuperBloque, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		return 0,0
	}

	if sbLeido.MagicNum == 201314821 {
		apRaiz := sbLeido.InicioAV
		dirActual := apRaiz
		//aqui ya tengo mi superbloque leido
		s := strings.Split(pathCrear, "/")
		if len(s) > 0 {
 			for i := 0 ; i< len(s) ; i++{
				if strings.Compare(s[i] , "")!= 0{
					index1, index2 := crearcarpetaRecursivoFILE(pathDisco, s[i], inicioSuperBloque, dirActual, atributoP, 0)
					if index1 == 1{
						dirActual = index2
					}else{
						return 0,0
					}
				}
			}
		return 1, dirActual
		}

	}
	return 0,0
}

//el int retorna si el directorio existe o fue creado
//el int64 retorna la estructura nueva o ya existente
func crearcarpetaRecursivoFILE(pathDisco string, carpeta string, inicioSuperBloque int64, avdActual int64, atributoP int, fin int) (int, int64){
	file, err := os.OpenFile(strings.ReplaceAll(pathDisco, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return 0,0
	}
	sbLeido := superbloque{}
	file.Seek(inicioSuperBloque, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		return 0,0
	}

	if sbLeido.MagicNum == 201314821 {
		//moverse a avdActual y leerlo
		avdLeido := avd{}
		file.Seek(avdActual, 0)
		data := readNextBytes(file, unsafe.Sizeof(avd{}))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &avdLeido)
		if err != nil {
			log.Fatal("binary.Read failed", err)
			return 0,0
		}
		//revisar entre el arreglo actual  
		var convertido [16]byte
		copy(convertido[:], strings.ToLower(strings.TrimSpace(carpeta)))

		if avdLeido.AVDNombreDirectorio == convertido{
			return 1, avdActual
		}
		contadorVacio := 0
		for i := 0; i < len(avdLeido.AVDApArraySubdirectorios); i++ {
			
			if avdLeido.AVDApArraySubdirectorios[i] != 0{
				avdNuevo := avd{}
				file.Seek(avdLeido.AVDApArraySubdirectorios[i], 0)
				data := readNextBytes(file, unsafe.Sizeof(avd{}))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &avdNuevo)
				if err != nil {
					log.Fatal("binary.Read failed", err)
					return 0,0
				}
				if avdNuevo.AVDNombreDirectorio == convertido {
					if fin == 1{
						fmt.Println("RESULTADO: El directorio ya se encuentra creado")
					}
					return 1, avdLeido.AVDApArraySubdirectorios[i]
				}

			}else {
				contadorVacio = contadorVacio + 1
			}	
		}

		if atributoP == 0 && fin == 0 && avdLeido.AVDApArbolVirtualDirectorio == 0{
			fmt.Println("RESULTADO: No se ha encontrado el directorio " + carpeta)
			return 0,0
		}else if avdLeido.AVDApArbolVirtualDirectorio != 0{
			return crearcarpetaRecursivoFILE(pathDisco, carpeta, inicioSuperBloque, avdLeido.AVDApArbolVirtualDirectorio, atributoP, fin)
		}
		
		//si sale es porque no encontro el nombre del directorio
		if contadorVacio > 0 && atributoP == 1{
			//todavia hay directorios disponibles en el array actual 

			/*Verificar si hay disponibles*/

			if sbLeido.ArbolVirtualFree > 0 && sbLeido.DetalleDirectorioFree > 0 {
				//si lo hay
				ingresoAVD := avd{}
				ingresoAVD.AVDFechaCreacion = getFechaHora()
				ingresoAVD.AVDNombreDirectorio = convertido
				
				ingresoDD := dd{} //no lleva parametros 


				//2 escribir el dd  nuevo y modificar el bitmap
				posactual := sbLeido.InicioDD + (sbLeido.PrimerLibreDD - 1)*int64(unsafe.Sizeof(dd{})) //posicion para escribir el avd
				ingresoAVD.AVDApDetalleDirectorio = posactual
				file.Seek(posactual, 0 )
				var binario2 bytes.Buffer
				binary.Write(&binario2, binary.BigEndian, ingresoDD)
				writeNextBytes(file, binario2.Bytes())
				
				posactual = sbLeido.InicioBMDD + sbLeido.PrimerLibreDD -1
				file.Seek(posactual, 0)
				var cambioBitmap2 byte = 1
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, cambioBitmap2)
				writeNextBytes(file, binario3.Bytes())

				//1 escribir el avd nuevo y modificar el bitmap
				posactual = sbLeido.InicioAV + ((sbLeido.PrimerLibreAV-1)*int64(unsafe.Sizeof(avd{})))
				retorno := posactual //posicion para escribir el avd
				file.Seek(posactual, 0 )
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, ingresoAVD)
				writeNextBytes(file, binario.Bytes())
				
				posactual = sbLeido.InicioBMAV + sbLeido.PrimerLibreAV -1
				file.Seek(posactual, 0)
				var cambioBitmap byte = 1
				var binario1 bytes.Buffer
				binary.Write(&binario1, binary.BigEndian, cambioBitmap)
				writeNextBytes(file, binario1.Bytes())
				//meter el arbol actual a su directorio padre 
				for i:=0 ; i < len(avdLeido.AVDApArraySubdirectorios); i++{
					if avdLeido.AVDApArraySubdirectorios[i] == 0{
						avdLeido.AVDApArraySubdirectorios[i] = retorno
						break
					}
				}
				posactual = avdActual
				file.Seek(avdActual, 0)
				var binario10 bytes.Buffer

				binary.Write(&binario10, binary.BigEndian, avdLeido)
				writeNextBytes(file, binario10.Bytes())
				//3 modificar el superbloque
				sbLeido.ArbolVirtualCount = sbLeido.ArbolVirtualCount + 1			//se suma 1 directorio
				sbLeido.DetalleDirectorioCount = sbLeido.DetalleDirectorioCount + 1	//se suma un detalle 
				sbLeido.ArbolVirtualFree = sbLeido.ArbolVirtualFree -1 				//un avd libre menos
				sbLeido.DetalleDirectorioFree = sbLeido.DetalleDirectorioFree -1 	//un dd  libre menos 
				sbLeido.PrimerLibreAV = sbLeido.PrimerLibreAV + 1					//avanza uno en el bitmap de avd
				sbLeido.PrimerLibreDD = sbLeido.PrimerLibreDD + 1					//avanza uno en el bitmap de dd
				//4 escribir el superbloque
				posactual = inicioSuperBloque
				file.Seek(posactual, 0)
				var binario5 bytes.Buffer
				binary.Write(&binario5, binary.BigEndian, sbLeido)
				writeNextBytes(file, binario5.Bytes())
				fmt.Println("RESULTADO: Se ha creado el directorio " + carpeta)
				return 1, retorno

			}else{
				//no hay espacio 
				fmt.Println("RESULTADO: No hay espacio disponible para crear el directorio")
				return 0,0
			}
		}else{
			return 0,0
		}

		//si no se debe crear un nuevo avd con el nombre del actual para meter el subdirectorio 
		//lo repito aqui mismo 
		if sbLeido.ArbolVirtualFree > 0 && sbLeido.DetalleDirectorioFree > 0 && atributoP == 1{
				//si lo hay
				ingresoAVD := avd{}
				ingresoAVD.AVDFechaCreacion = avdLeido.AVDFechaCreacion
				ingresoAVD.AVDNombreDirectorio = avdLeido.AVDNombreDirectorio
				
				ingresoDD := dd{} //no lleva parametros 


				//2 escribir el dd  nuevo y modificar el bitmap
				posactual := sbLeido.InicioDD + (sbLeido.PrimerLibreDD - 1)*int64(unsafe.Sizeof(dd{})) //posicion para escribir el avd
				ingresoAVD.AVDApDetalleDirectorio = posactual
				file.Seek(posactual, 0 )
				var binario2 bytes.Buffer
				binary.Write(&binario2, binary.BigEndian, ingresoDD)
				writeNextBytes(file, binario2.Bytes())
				
				posactual = sbLeido.InicioBMDD + sbLeido.PrimerLibreDD -1
				file.Seek(posactual, 0)
				var cambioBitmap2 byte = 1
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, cambioBitmap2)
				writeNextBytes(file, binario3.Bytes())

				//1 escribir el avd nuevo y modificar el bitmap
				posactual = sbLeido.InicioAV + ((sbLeido.PrimerLibreAV-1)*int64(unsafe.Sizeof(avd{})))
				avdLeido.AVDApArbolVirtualDirectorio = posactual
				retorno2 := posactual //posicion para escribir el avd
				file.Seek(posactual, 0 )
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, ingresoAVD)
				writeNextBytes(file, binario.Bytes())
				
				posactual = sbLeido.InicioBMAV + sbLeido.PrimerLibreAV -1
				file.Seek(posactual, 0)
				var cambioBitmap byte = 1
				var binario1 bytes.Buffer
				binary.Write(&binario1, binary.BigEndian, cambioBitmap)
				writeNextBytes(file, binario1.Bytes())
				//meter el arbol actual a su directorio padre 
				for i:=0 ; i < len(avdLeido.AVDApArraySubdirectorios); i++{
					if avdLeido.AVDApArraySubdirectorios[i] == 0{
						avdLeido.AVDApArraySubdirectorios[i] = retorno2
						i = 100
					}
				}
				posactual = avdActual
				file.Seek(posactual, 0)
				var binario10 bytes.Buffer
				binary.Write(&binario10, binary.BigEndian, avdLeido)
				writeNextBytes(file, binario10.Bytes())
				//3 modificar el superbloque
				sbLeido.ArbolVirtualCount = sbLeido.ArbolVirtualCount + 1			//se suma 1 directorio
				sbLeido.DetalleDirectorioCount = sbLeido.DetalleDirectorioCount + 1	//se suma un detalle 
				sbLeido.ArbolVirtualFree = sbLeido.ArbolVirtualFree -1 				//un avd libre menos
				sbLeido.DetalleDirectorioFree = sbLeido.DetalleDirectorioFree -1 	//un dd  libre menos 
				sbLeido.PrimerLibreAV = sbLeido.PrimerLibreAV + 1					//avanza uno en el bitmap de avd
				sbLeido.PrimerLibreDD = sbLeido.PrimerLibreDD + 1					//avanza uno en el bitmap de dd
				//4 escribir el superbloque
				posactual = inicioSuperBloque
				file.Seek(posactual, 0)
				var binario5 bytes.Buffer
				binary.Write(&binario5, binary.BigEndian, sbLeido)
				writeNextBytes(file, binario5.Bytes())
				return crearcarpetaRecursivoFILE(pathDisco, carpeta, inicioSuperBloque, retorno2, atributoP, fin)
			}else{
				//no hay espacio 
				fmt.Println("RESULTADO: No hay espacio disponible para crear el directorio")
				return 0,0
			}

		//*****************************************************************************************

	}
	return 0,0
}
