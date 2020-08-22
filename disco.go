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
	enc := 0
	for i := 0; i < len(disco.Tabla); i++ {
		if enc == 0 {
			if disco.Tabla[i].Type == 'e' {
				enc = 1
				//se encontro la particion extendida
				//verificar si el tamano alcanza
				if nueva.Size <= disco.Tabla[i].Size-int64(unsafe.Sizeof(ebr{})) {
					//alcanza el tamano
					//leer el ebr inicial
					nextEBR, sizeEbr := leerEbr(path, disco.Tabla[i].Start, nueva)
					if nextEBR == -1 && sizeEbr == 0 {
						ebrEs := ebr{}
						ebrEs.Fit = nueva.Fit
						ebrEs.Name = nueva.Name
						ebrEs.Next = -1
						ebrEs.Size = nueva.Size
						ebrEs.Start = disco.Tabla[i].Start
						ebrEs.Status = 0
						escribirEbr(path, ebrEs.Start, &ebrEs)
						return 1
						//significa que la particion esta vacia
					} else {
						//hay particiones logicas ya registradas
						ebrEs := ebr{}
						ebrEs.Fit = nueva.Fit
						ebrEs.Name = nueva.Name
						ebrEs.Next = -1
						ebrEs.Size = nueva.Size
						ebrEs.Start = disco.Tabla[i].Start
						ebrEs.Status = 0
						return logicaRecursiva(path, disco.Tabla[i].Start, &ebrEs, disco.Tabla[i].Size)
					}
				} else {
					//no alcanza el tamano
					fmt.Println("RESULTADO: La particion logica excede el tamano de la particion extendida")
					return 0
				}
			}
		}
	}
	return 0
}

func logicaRecursiva(path string, pos int64, ebrNuevo *ebr, limite int64) int {
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
	if ebrTemp.Next == -1 {
		//es el ultimo ebr
		//no puede ser el primero porque ya se valido antes
		disponible := limite - ebrTemp.Size - int64(unsafe.Sizeof(ebr{}))
		if disponible >= ebrNuevo.Size {
			ebrEs := ebr{}
			ebrEs.Fit = ebrNuevo.Fit
			ebrEs.Name = ebrNuevo.Name
			ebrEs.Next = -1
			ebrEs.Size = ebrNuevo.Size
			ebrEs.Start = ebrNuevo.Start
			ebrEs.Status = 0

			ebrTemp.Next = ebrEs.Start

			escribirEbr(path, ebrEs.Start, &ebrEs)
			escribirEbr(path, ebrTemp.Start, &ebrTemp)
			return 1
		}
		//si no alcanza el tamao sale de aca
		fmt.Println("RESULTADO: Espacio insuficiente para crear la particion logica")
		return 0
	} else {
		//no es el ultimo ebr
		disponible := ebrTemp.Size - int64(unsafe.Sizeof(ebr{}))
		if disponible >= ebrNuevo.Size {
			ebrEs := ebr{}
			ebrEs.Fit = ebrNuevo.Fit
			ebrEs.Name = ebrNuevo.Name
			ebrEs.Next = -1
			ebrEs.Size = ebrNuevo.Size
			ebrEs.Start = ebrNuevo.Start
			ebrEs.Status = 0

			ebrTemp.Next = ebrEs.Start

			escribirEbr(path, ebrEs.Start, &ebrEs)
			escribirEbr(path, ebrTemp.Start, &ebrTemp)
			return 1
		}
		//si no alcanza el tamao sale de aca

		return logicaRecursiva(path, ebrTemp.Next, ebrNuevo, limite)
	}
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
}

func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

/**************************************************************
	Eliminar disco
***************************************************************/
func removerDisco(path string) {
	// borrar el archivo archivoBorrable.txt
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}
