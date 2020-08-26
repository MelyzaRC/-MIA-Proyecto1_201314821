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
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"
)

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

/**************************************************************
	Crear particion
	-El path ya fue validado, asi como que sea un disco
***************************************************************/
func crearParticion(path string, size int, unit string, name string, tipo string, fit string) {
	//leyendo el disco para saber sus atributos
	var disco *mbr = leerDisco(path)

	//determinando si el disco es nulo o no
	if disco != nil {
		//se leyo el disco sin problemas

		//obteniendo el tamano total actual
		var tablaParticiones [4]particion = disco.Tabla
		var tamanoOcupado = int64(unsafe.Sizeof(mbr{}))
		for i := 0; i < len(tablaParticiones); i++ {
			particionActual := tablaParticiones[i]
			if particionActual.Size != 0 {
				tamanoOcupado = tamanoOcupado + particionActual.Size
			}
		}
		//El espacio actual dispobible sera lo total - lo ocupado
		var tamDisponible = int64(disco.Tamano - tamanoOcupado)

		//determinando tamano total de la particion a crear
		var tamanoParticion int64 = 0
		switch strings.ToLower(unit) {
		case "b":
			tamanoParticion = int64(size)
		case "k":
			tamanoParticion = int64(size) * 1024
		case "m":
			tamanoParticion = int64(size) * 1024 * 1024
		default:
		}

		//determinando si hay espacio suficiente
		if tamDisponible-tamanoParticion >= 0 {
			//si hay espacio para crear la particion
			//formando la particion para pasarla al siguiente metodo
			nueva := particion{}
			nueva.Status = 0
			switch strings.ToLower(tipo) {
			case "p":
				nueva.Type = 'p'
			case "e":
				nueva.Type = 'e'
			case "l":
				nueva.Type = 'l'
			default:
			}
			switch strings.ToLower(fit) {
			case "ff":
				nueva.Fit = 'f'
			case "bf":
				nueva.Fit = 'b'
			case "wf":
				nueva.Fit = 'w'
			default:
			}
			nueva.Start = 0
			nueva.Size = tamanoParticion
			copy(nueva.Name[:], name)
			//la particion esta formada
			//mandar a colocarla en el disco
			res := insertarParticion(disco, &nueva, path)
			if res == 1 {
				reordenarParticiones(disco)
				reescribir(disco, path)
				/*Si es extendida**************************************************/
				if nueva.Type == 'e' {
					ebrVacio := ebr{}
					ebrVacio.Start = nueva.Start
					ebrVacio.Next = -1
					ebrVacio.Size = 0
					file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_RDWR, os.ModeAppend)
					defer file.Close()
					if err != nil {
						log.Fatal(err)
					} else {
						file.Seek(nueva.Start, 0)
						//Escribiendo el MBR
						var binario3 bytes.Buffer
						binary.Write(&binario3, binary.BigEndian, ebrVacio)
						writeNextBytes(file, binario3.Bytes())
					}
				}
				/******************************************************************/
				//fmt.Println(disco)
				fmt.Println("RESULTADO: Particion creada con exito")
				graficarMBR(path)
			} /*else {
				fmt.Println("RESULTADO: No se ha podido crear la particion")
			}*/
		} else {
			//no hay espacio para crear la particion
			fmt.Println("RESULTADO: No hay espacio suficiente para crear la particion")
		}
	} else {
		//si el disco regresa un nulo
		fmt.Println("RESULTADO: No se puede leer el disco")
	}

}

//primer ajuste para guardar particiones
//espacio total disponible si es suficiente
func insertarParticion(disco *mbr, nueva *particion, path string) int {
	var tabla = [4]particion{}
	tabla = disco.Tabla
	//verificar si se puede segun el tipo
	libre, primaria, extendida := 0, 0, 0
	for i := 0; i < len(tabla); i++ {
		if tabla[i].Size != 0 {
			if tabla[i].Name == nueva.Name {
				fmt.Println("RESULTADO: No se puede repetir el nombre de una particion")
				return 0
			}
		}
		switch tabla[i].Type {
		case 'p':
			primaria++
		case 'e':
			extendida++
		default:
			libre++
		}
	}
	//verificar la teoria de particiones
	//1 logica
	//para que hayan extendidas tiene que haber una logica
	//sumar 4

	if libre == 4 && nueva.Type != 'l' {
		//no hay particiones creadas aun
		nueva.Start = int64(unsafe.Sizeof(mbr{}))
		disco.Tabla[0] = *nueva
		return 1
	} else if libre == 4 && nueva.Type == 'l' {
		//ya hay particiones
		fmt.Print("RESULTADO: No se puede crear la particion logica, debe crear una particion extendida")
		return 0
	} else if libre == 0 && nueva.Type != 'l' {
		fmt.Println("RESULTADO: No se pueden crear mas particiones primarias ni extendidas")
		return 0
	} else if extendida > 0 && nueva.Type == 'l' {
		return creacionL(disco, nueva, path) //aqui que retorne lo que retorna el otro metodo que tengo que crear
	} else if nueva.Type == 'e' && extendida > 0 && libre > 0 {
		fmt.Println("RESULTADO: Solamente se puede crear una particion extendida")
		return 0
	} else if libre > 0 && nueva.Type != 'l' {
		return creacionPE(disco, nueva, path)
	}
	return 0
}

//tomando en cuenta que el arreglo de particiones esta en orden de part_start

//para particiones logicas y primarias
func creacionPE(disco *mbr, nueva *particion, path string) int {
	var inicioEspacio int64 = int64(unsafe.Sizeof(mbr{}))
	ingresoOk := 0
	var finalAnterior int64 = inicioEspacio

	//determinando el part_start
	for i := 0; i < 4; i++ {
		if disco.Tabla[i].Size != 0 {
			inicioActual := disco.Tabla[i].Start
			if inicioActual-finalAnterior >= nueva.Size {
				nueva.Start = finalAnterior
				ingresoOk = 1
			} else {
				finalAnterior = inicioActual + disco.Tabla[i].Size
			}
		}
	}

	if ingresoOk == 1 {
		//meter la particion en el primer nulo y ordenar
		for i := 0; i < 4; i++ {
			if disco.Tabla[i].Size == 0 {
				disco.Tabla[i] = *nueva
				return 1
			}
		}
		return 1
	}

	if disco.Tamano-finalAnterior >= 0 {
		nueva.Start = finalAnterior
		//meter la particion en el primer nulo y ordenar
		for i := 0; i < 4; i++ {
			if disco.Tabla[i].Size == 0 {
				disco.Tabla[i] = *nueva
				return 1
			}
		}
		return 1
	}
	fmt.Println("RESULTADO: No hay espacio adecuado para acomodar la particion")
	return 0

}

//para particiones extendidas
//ya se verifico que si hay una particion logica
func creacionL(disco *mbr, nueva *particion, path string) int {
	//encontrando la particion logica
	//este for recorre la tabla de particiones del disco en busca de una extendida
	for i := 0; i < len(disco.Tabla); i++ {
		if disco.Tabla[i].Type == 'e' {
			if disco.Tabla[i].Name == nueva.Name {
				fmt.Println("RESULTADO: Nombre de particion logica repetido")
				return 0
			}
			ebrEs := ebr{}
			ebrEs.Fit = nueva.Fit
			ebrEs.Name = nueva.Name
			ebrEs.Next = 0
			ebrEs.Size = nueva.Size
			ebrEs.Start = disco.Tabla[i].Start
			ebrEs.Status = 0

			return logicaRecursiva(path, disco.Tabla[i].Start, &ebrEs, disco.Tabla[i].Start+disco.Tabla[i].Size)
		}
	}
	//no encontro ninguna particion extendida
	fmt.Println("RESULTADO: No se encontro ninguna particion extendida creada")
	return 0
}

func logicaRecursiva(path string, pos int64, ebrNuevo *ebr, limite int64) int {
	//leyendo el archivo
	ebrTemp := ebr{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	file.Seek(int64(pos), 0)
	data := readNextBytes(file, unsafe.Sizeof(ebr{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		return 0
	}
	if &ebrTemp != nil {
		for i := ebrTemp.Start; i < limite; i++ {
			ebrLeido := ebr{}
			file.Seek(i, 0)
			data := readNextBytes(file, unsafe.Sizeof(ebr{}))
			buffer := bytes.NewBuffer(data)
			err = binary.Read(buffer, binary.BigEndian, &ebrLeido)
			if err != nil {
				log.Fatal("binary.Read failed", err)
				return 0
			}
			if &ebrLeido != nil {
				if ebrLeido.Name == ebrNuevo.Name {
					fmt.Println("RESULTADO: El nombre de la particion logica esta repetido")
					return 0
				}
				if ebrLeido.Next == -1 && ebrLeido.Size == 0 {
					//quiere decir que esta vacio
					disponible := limite - int64(unsafe.Sizeof(ebr{}))
					if disponible >= ebrNuevo.Size {
						ebrLeido.Fit = ebrNuevo.Fit
						ebrLeido.Name = ebrNuevo.Name
						ebrLeido.Size = ebrNuevo.Size
						ebrLeido.Status = 0
						ebrLeido.Next = -1
						escribirEbr(path, ebrLeido.Start, &ebrLeido)
						fmt.Print("EBR start: ")
						fmt.Println(ebrLeido.Start)
						return 1
					}
					fmt.Println("RESULTADO: no hay espacio disponible para crear la particion logica en esta particion")
					return 0
				}
				if ebrLeido.Next == -1 { //lego al utimo ebr
					disponible := limite - int64(unsafe.Sizeof(ebr{}))
					if disponible >= ebrNuevo.Size {
						ebrEs := ebr{}
						ebrEs.Fit = ebrNuevo.Fit
						ebrEs.Name = ebrNuevo.Name
						ebrEs.Next = -1
						ebrEs.Size = ebrNuevo.Size
						ebrEs.Start = ebrLeido.Start + int64(unsafe.Sizeof(ebr{})) + ebrLeido.Size
						ebrEs.Status = 0

						ebrLeido.Next = ebrEs.Start

						escribirEbr(path, ebrEs.Start, &ebrEs)
						escribirEbr(path, ebrLeido.Start, &ebrLeido)
						fmt.Print("EBR Start:")
						fmt.Println(ebrEs.Start)
						return 1
					}
					fmt.Println("RESULTADO: no hay espacio disponible para crear la particion logica en esta particion")
					return 0
				} else if ebrLeido.Next != -1 { //esta en los ebr antes del ultimo
					disponible := ebrLeido.Next - int64(unsafe.Sizeof(ebr{})) - ebrLeido.Size - ebrLeido.Start
					if disponible >= ebrNuevo.Size {
						ebrEs := ebr{}
						ebrEs.Fit = ebrNuevo.Fit
						ebrEs.Name = ebrNuevo.Name
						ebrEs.Next = ebrLeido.Next
						ebrEs.Size = ebrNuevo.Size
						ebrEs.Start = ebrLeido.Start + int64(unsafe.Sizeof(ebr{})) + ebrLeido.Size
						ebrEs.Status = 0

						ebrLeido.Next = ebrEs.Start

						escribirEbr(path, ebrEs.Start, &ebrEs)
						escribirEbr(path, ebrLeido.Start, &ebrLeido)
						fmt.Print("EBR Star: ")
						fmt.Println(ebrEs.Start)
						return 1
					}
					//porque al iterar el for le suma uno
					i = ebrLeido.Next - 1
				}

			}

		}
		fmt.Println("RESULTADO: no hay espacio disponible para crear la particion logica")
		return 0
	}
	return 0
}

func escribirEbr(path string, pos int64, ebr *ebr) {
	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		file.Seek(pos, 0)
		//Escribiendo el MBR
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, ebr)
		writeNextBytes(file, binario3.Bytes())
	}
}

func leerEbr(path string, pos int64, particion *particion) (int64, int64) {
	ebrTemp := ebr{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		file.Seek(int64(pos), 0)
		data := readNextBytes(file, unsafe.Sizeof(ebr{}))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
		if err != nil {
			log.Fatal("binary.Read failed", err)
			return 0, 0
		}
		return ebrTemp.Next, ebrTemp.Size
	}
	return 0, 0
}

//ordena la tabla de particiones del mbr de menor a mayor
func reordenarParticiones(disco *mbr) {
	tabla := disco.Tabla
	for i := 0; i < len(tabla); i++ {
		for j := 0; j < len(tabla)-1; j++ {
			if disco.Tabla[j].Start > disco.Tabla[j+1].Start {
				temporal := disco.Tabla[j]
				disco.Tabla[j] = disco.Tabla[j+1]
				disco.Tabla[j+1] = temporal
			}
		}
	}
}

//reescribe el disco con las actualizaciones
func reescribir(disco *mbr, path string) {
	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	} else {
		file.Seek(0, 0)
		//Escribiendo el MBR
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, disco)
		writeNextBytes(file, binario3.Bytes())
	}
}

/**************************************************************
	Leer disco
***************************************************************/
func leerDisco(path string) *mbr {
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
	var mDir *mbr = &m
	return mDir
}

func readNextBytes(file *os.File, number uintptr) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

/**************************************************************
	Crear disco
***************************************************************/
func crearDisco(tam int, unit string, ruta string) {
	archivo, err := os.Create(ruta)
	defer archivo.Close()
	if err != nil {
		log.Fatal(err)
	}
	var vacio int8 = 0
	s := &vacio
	var num int64 = 0
	//Definiendo tamano
	if strings.Compare(strings.ToLower(unit), "m") == 0 {
		num = int64(tam) * 1024 * 1024
	} else if strings.Compare(strings.ToLower(unit), "k") == 0 {
		num = int64(tam) * 1024
	}
	num = num - 1
	//Llenando el archivo

	//colocando el primer byte
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	writeNextBytes(archivo, binario.Bytes())

	//situando el cursor en la ultima posicion
	archivo.Seek(num, 0)

	//colocando el ultimo byte para rellenar
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	writeNextBytes(archivo, binario2.Bytes())

	//Regresando el cursor a 0 para escribir el mbr
	archivo.Seek(0, 0)

	//Formando el MBR
	disco := mbr{}
	disco.Tamano = num + 1

	fechahora := time.Now()
	fechahoraArreglo := strings.Split(fechahora.String(), "")
	fechahoraCadena := ""
	for i := 0; i < 16; i++ {
		fechahoraCadena = fechahoraCadena + fechahoraArreglo[i]
	}
	copy(disco.Fecha[:], fechahoraCadena)

	var signature int8
	binary.Read(rand.Reader, binary.LittleEndian, &signature)
	if signature < 0 {
		signature = signature * -1
	}
	disco.Firma = signature

	//Escribiendo el MBR
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, disco)
	writeNextBytes(archivo, binario3.Bytes())
	graficarMBR(ruta)
}

func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

/**************************************************************
	Eliminar disco
	-Ya se valido que existe y que es un disco
***************************************************************/
func removerDisco(path string) {
	// borrar el archivo
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}

/**************************************************************
	Eliminar particion
***************************************************************/
func eliminarParticion(path string, nombre string, tipo string) {
	s := leerDisco(path)
	nombre = strings.ReplaceAll(strings.ToLower(nombre), "\"", "")
	var nombreComparar [16]byte
	copy(nombreComparar[:], nombre)
	eliminado := 0
	if s != nil {
		for i := 0; i < len(s.Tabla); i++ {
			if eliminado == 0 {
				//recorriendo las particiones
				if nombreComparar == s.Tabla[i].Name {
					//encontro la particion entre las primarias y extendidas
					if strings.Compare(tipo, "fast") == 0 {
						particionVacia := particion{}
						s.Tabla[i] = particionVacia
						reordenarParticiones(s)
						reescribir(s, path)
						graficarMBR(path)
					} else if strings.Compare(tipo, "full") == 0 {
						h := s.Tabla[i].Start
						h2 := s.Tabla[i].Size
						particionVacia := particion{}
						s.Tabla[i] = particionVacia
						reordenarParticiones(s)
						reescribir(s, path)

						borrarFullParticion(path, h, h2)

						graficarMBR(path)
					}
					eliminado = 1

				} else if s.Tabla[i].Type == 'e' {
					/*Verificar las logicas dentro de la particion extendida*/
					/*Leyendo logicas***************************************/
					limite := s.Tabla[i].Start + int64(unsafe.Sizeof(ebr{})) + s.Tabla[i].Size
					ebrTemp := ebr{}
					ebrAnterior := ebr{}
					ebrAnterior.Start = 0
					ebrAnterior.Next = 0
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
					if &ebrTemp != nil {
						for i := ebrTemp.Start; i < limite; i++ {
							ebrLeido := ebr{}
							file.Seek(i, 0)
							data := readNextBytes(file, unsafe.Sizeof(ebr{}))
							buffer := bytes.NewBuffer(data)
							err = binary.Read(buffer, binary.BigEndian, &ebrLeido)
							if err != nil {
								log.Fatal("binary.Read failed", err)
							}
							if &ebrLeido != nil {

								if ebrLeido.Next == -1 { //lego al utimo ebr
									if ebrLeido.Name == nombreComparar {
										ebrVacio := ebr{}
										ebrVacio.Next = ebrLeido.Next
										ebrVacio.Size = 0
										ebrVacio.Start = ebrLeido.Start
										if ebrAnterior.Next == 0 {
											if strings.Compare(tipo, "fast") == 0 {
												escribirEbr(path, ebrLeido.Start, &ebrVacio)
											} else if strings.Compare(tipo, "full") == 0 {
												borrarFullParticion(path, ebrLeido.Start+int64(unsafe.Sizeof(ebr{})), s.Tabla[i].Size)
												escribirEbr(path, ebrLeido.Start, &ebrVacio)
											}
										} else {
											ebrAnterior.Next = ebrLeido.Next
											if strings.Compare(tipo, "fast") == 0 {
												escribirEbr(path, ebrAnterior.Start, &ebrAnterior)
											} else if strings.Compare(tipo, "full") == 0 {
												borrarFullParticion(path, ebrAnterior.Start+int64(unsafe.Sizeof(ebr{})), ebrAnterior.Next-ebrAnterior.Start-ebrAnterior.Size-int64(unsafe.Sizeof(ebr{})))
												escribirEbr(path, ebrLeido.Start, &ebrVacio)
											}
										}
										eliminado = 1
										i = limite + 1
									}
									i = limite + 1
								} else if ebrLeido.Next != -1 { //esta en los ebr antes del ultimo
									//verificar pero con el next
									if ebrLeido.Name == nombreComparar {
										ebrVacio := ebr{}
										ebrVacio.Next = ebrLeido.Next
										ebrVacio.Size = 0
										ebrVacio.Start = ebrLeido.Start
										if ebrAnterior.Next == 0 {
											if strings.Compare(tipo, "fast") == 0 {
												escribirEbr(path, ebrLeido.Start, &ebrVacio)
											} else if strings.Compare(tipo, "full") == 0 {
												borrarFullParticion(path, ebrLeido.Start+int64(unsafe.Sizeof(ebr{})), s.Tabla[i].Size)
												escribirEbr(path, ebrLeido.Start, &ebrVacio)
											}
										} else {
											ebrAnterior.Next = ebrLeido.Next
											if strings.Compare(tipo, "fast") == 0 {
												escribirEbr(path, ebrAnterior.Start, &ebrAnterior)
											} else if strings.Compare(tipo, "full") == 0 {
												borrarFullParticion(path, ebrAnterior.Start+int64(unsafe.Sizeof(ebr{})), ebrAnterior.Next-ebrAnterior.Start-ebrAnterior.Size-int64(unsafe.Sizeof(ebr{})))
												escribirEbr(path, ebrLeido.Start, &ebrVacio)
											}
										}
										eliminado = 1
										i = limite + 1
									}
									ebrAnterior = ebrLeido
									//porque al iterar el for le suma uno
									i = ebrLeido.Next - 1
								}

							}

						}
					}
					/*Final leyendo logicas*********************************/
				}
			}
		}
		if eliminado == 1 {
			graficarMBR(path)
			fmt.Println("RESULTADO: Se ha eliminado correctamente la particion")
		} else {
			fmt.Println("RESULTADO: No se ha podido eliminar la particion")
		}
	}
}

func borrarFullParticion(path string, inicio int64, tam int64) {
	archivo, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_RDWR, os.ModeAppend)
	defer archivo.Close()
	var vacio int8 = 0
	s := &vacio
	if err != nil {
		log.Fatal(err)
	} else {
		for i := inicio; i < inicio+tam; i++ {
			archivo.Seek(0, 0)
			archivo.Seek(i, 0)
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, s)
			writeNextBytes(archivo, binario.Bytes())
		}
	}
}

/**************************************************************
	Modificar tamano de particion
***************************************************************/
//el tamano viene en bytes
func modificarParticion(path string, nombre string, tam int64, tipo string) {
	s := leerDisco(path)
	nombre = strings.ReplaceAll(strings.ToLower(nombre), "\"", "")
	var nombreComparar [16]byte
	copy(nombreComparar[:], nombre)
	modificado := 0

	if s != nil {
		finAnterior := s.Tamano - 1
		for i := len(s.Tabla) - 1; i >= 0; i-- {
			if modificado == 0 {
				//recorriendo las particiones
				if s.Tabla[i].Size == 0 {
					//particion vacia
					//fin anterior se queda como esta

				} else {
					if nombreComparar == s.Tabla[i].Name {
						//encontro la particion en la primera tabla
						if strings.Compare(tipo, "quitar") == 0 {
							if s.Tabla[i].Type != 'e' {
								res := s.Tabla[i].Size - tam
								if res > 0 {
									s.Tabla[i].Size = res
									i = -10
									modificado = 1
									break
								} else {
									i = -10
									fmt.Println("RESULTADO: No se puede reducir el espacio en la particion")
								}
							} else {
								//verificar si hay espacio en la extendida con base a las logicas
								espacioD := obtenerUtilizacionEBR(path, s.Tabla[i].Start, s.Tabla[i].Size)
								if espacioD >= tam {
									s.Tabla[i].Size = s.Tabla[i].Size - tam
									i = -10
									modificado = 1
									break
								} else {
									i = -10
									fmt.Println("RESULTADO: No se puede reducir el espacio en la particion extendida")
								}
							}
						} else if strings.Compare(tipo, "agregar") == 0 {
							res := finAnterior - s.Tabla[i].Start + s.Tabla[i].Size
							if res >= tam {
								s.Tabla[i].Size = res
								modificado = 1
								i = -10
								break
							} else {
								fmt.Println("RESULTADO: No se puede ampliar el tama;o de la particion")
							}
						}
					} else if s.Tabla[i].Type == 'e' {
						//puede que la particion a modificar sea una logica****
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
									fmt.Print("Leido next ")
									fmt.Println(ebrLeido.Next)
									fmt.Println("Leido size")
									fmt.Println(ebrLeido.Size)
									if ebrLeido.Next != -1 && ebrLeido.Size == 0 {

										//Aqui no valuo porque es el primer EBR solo lo salto
										j = ebrLeido.Next - 1
										break
									} else if ebrLeido.Next == -1 && ebrLeido.Size == 0 {
										//la particion esta vacia
										j = limite + 1
										break
									} else if ebrLeido.Next == -1 && ebrLeido.Size > 0 { //lego al utimo ebr
										//valuo con el limite porque es el ultimo ebr
										if ebrLeido.Name == nombreComparar {
											if strings.Compare(tipo, "quitar") == 0 {
												espacio := ebrLeido.Size - tam
												if espacio > 0 {
													ebrLeido.Size = espacio
													modificado = 1
													j = limite + 1
													escribirEbr(path, ebrLeido.Start, &ebrLeido)
													break
												} else {
													j = limite + 1
													fmt.Print("RESULTADO: No se puede reducir la particion logica")
												}
											} else if strings.Compare(tipo, "agregar") == 0 {
												espacioDisponible := limite - ebrLeido.Start + ebrLeido.Size
												if espacioDisponible >= tam {
													ebrLeido.Size = ebrLeido.Size + tam
													escribirEbr(path, ebrLeido.Start, &ebrLeido)
													modificado = 1
													j = limite + 1
													break
												} else {
													j = limite + 1
													fmt.Print("RESULTADO: No se puede ampliar el tamano de la particion logica")
													break
												}
											}
										} else {
											j = ebrLeido.Next - 1
											break
										}

									} else if ebrLeido.Next != -1 && ebrLeido.Size > 0 { //esta en los ebr antes del ultimo
										//verificar pero con el next
										if ebrLeido.Name == nombreComparar {
											if strings.Compare(tipo, "quitar") == 0 {
												espacio := ebrLeido.Size - tam
												if espacio > 0 {
													ebrLeido.Size = espacio
													modificado = 1
													j = limite + 1
													escribirEbr(path, ebrLeido.Start, &ebrLeido)
													break
												} else {
													j = limite + 1
													fmt.Print("RESULTADO: No se puede reducir la particion logica")
												}
											} else if strings.Compare(tipo, "agregar") == 0 {
												espacioDisponible := ebrLeido.Next - ebrLeido.Start + ebrLeido.Size
												if espacioDisponible >= tam {
													ebrLeido.Size = ebrLeido.Size + tam
													escribirEbr(path, ebrLeido.Start, &ebrLeido)
													modificado = 1
													j = limite + 1
													break
												} else {
													j = limite + 1
													fmt.Print("RESULTADO: No se puede ampliar el tamano de la particion logica")
												}
											}
										} else {
											fmt.Print("Aqui seria a fuerzas ")
											j = ebrLeido.Next - 1
										}

									}
								}
								if modificado == 0 {
									j = ebrLeido.Next - 1
								}
							}
							if modificado == 0 {
								finAnterior = s.Tabla[i].Start - 1
							} else {
								i = -1
							}
						}
						if modificado == 0 {
							finAnterior = s.Tabla[i].Start - 1
						}
						//fin de las logicas***********************************
					} else {
						finAnterior = s.Tabla[i].Start - 1
					}
				}
			}
		}
		if modificado == 1 {
			fmt.Println("RESULTADO: Se ha modificado el tamano de la particion")
			reescribir(s, path)
			graficarMBR(path)
		}
	}
}

func obtenerUtilizacionEBR(path string, inicio int64, final int64) int64 {
	ebrTemp := ebr{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicio, 0)
	data := readNextBytes(file, unsafe.Sizeof(ebr{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	limite := inicio + final

	if &ebrTemp != nil {
		for i := ebrTemp.Start; i < limite; i++ {
			ebrLeido := ebr{}
			file.Seek(i, 0)
			data := readNextBytes(file, unsafe.Sizeof(ebr{}))
			buffer := bytes.NewBuffer(data)
			err = binary.Read(buffer, binary.BigEndian, &ebrLeido)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}
			if &ebrLeido != nil {
				if ebrLeido.Next == -1 && ebrLeido.Size == 0 {
					//particion vacia
					return limite - int64(unsafe.Sizeof(ebr{}))
				} else if ebrLeido.Next == -1 && ebrLeido.Size > 0 {
					return limite - ebrLeido.Start - int64(unsafe.Sizeof(ebr{})) - ebrLeido.Size
				} else {
					i = ebrLeido.Next - 1
				}
			}

		}
	}
	return 0
}
