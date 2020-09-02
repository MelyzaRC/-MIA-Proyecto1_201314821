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
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
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
					graficarMBR(path)
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
func formatear(idFormatear string, tipoFormato string) { //convertir a megas
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
						fmt.Println(pathEnviar)
						fmt.Println(nombre)
						tipoParticion, inicioParticion, tamParticion = obtenerDatosParticion(pathEnviar, *nombre)
						if tipoParticion == 0 {
							fmt.Println("RESULTADO: No se encuentra la particion")
						} else if tipoParticion == 3 {
							fmt.Println("RESULTADO: No se puede formatear una particion extendida")
						} else if tipoParticion == 1 || tipoParticion == 2 {
							realizarFormato(pathEnviar, inicioParticion, tamParticion)
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
						for j := ebrTemp.Start; j < limite; j++ {
							ebrLeido := ebr{}

							file.Seek(j, 0)
							data := readNextBytes(file, unsafe.Sizeof(ebr{}))
							buffer := bytes.NewBuffer(data)
							err = binary.Read(buffer, binary.BigEndian, &ebrLeido)
							if err != nil {
								log.Fatal("binary.Read failed", err)
							}
							if &ebrLeido != nil {
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

func realizarFormato(path string, inicioParticion int64, tamParticion int64) {

	//datos para formula
	var tamAVD int64 = int64(unsafe.Sizeof(avd{}))
	var tamDD int64 = int64(unsafe.Sizeof(dd{}))
	var tamInodo int64 = int64(unsafe.Sizeof(inodo{}))
	var tamBloque int64 = int64(unsafe.Sizeof(bloque{}))
	var tamBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
	var tamSuperBloque int64 = int64(unsafe.Sizeof(superbloque{}))

	//formula
	var nEstructuras = (tamParticion - (2 * tamSuperBloque)) / (27 + tamAVD + tamDD + (5*tamInodo + (20 * tamBloque) + tamBitacora))

	//cantidad de tipos de estructura
	cantidadAVD := nEstructuras
	cantidadDD := nEstructuras
	cantidadInodos := 5 * nEstructuras
	cantidadBloques := 20 * nEstructuras
	//cantidadBitacoras := nEstructuras

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
	//iniciocopiaSB := iniciobitacora + (tamBitacora * cantidadBitacoras)
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
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	//formando el superbloque
	nuevoSB := superbloque{}
	//nuevoSB.NombreHD
	nuevoSB.ArbolVirtualCount = cantidadAVD
	nuevoSB.DetalleDirectorioCount = cantidadDD
	nuevoSB.InodosCount = cantidadInodos
	nuevoSB.BloquesCount = cantidadBloques
	nuevoSB.ArbolVirtualFree = cantidadAVD
	nuevoSB.DetalleDirectorioFree = cantidadDD
	nuevoSB.InodosFree = cantidadInodos
	nuevoSB.BloquesFree = cantidadBloques
	nuevoSB.DateCreacion = getFechaHora()
	nuevoSB.DateUltimoMontaje = getFechaHora()
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
	nuevoSB.PrimerLibreAV = iniciobitmapAVD
	nuevoSB.PrimerLibreDD = iniciobitmapDD
	nuevoSB.PrimerLibreInodo = iniciobitMapInodo
	nuevoSB.PrimerLibreBloque = iniciobitmapBloque
	nuevoSB.MagicNum = 201314821

	/*Escribiendo el superbloque*/
	file.Seek(inicioSuperBloque, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, nuevoSB)
	writeNextBytes(file, binario.Bytes())

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
