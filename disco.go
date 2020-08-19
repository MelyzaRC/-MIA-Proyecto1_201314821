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
				fmt.Println("RESULTADO: Particion creada")
				fmt.Print(disco)
			} else {
				fmt.Println("RESULTADO: No se ha podido crear la particion")
			}
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
	libre, primaria, logica := 0, 0, 0
	for i := 0; i < len(tabla); i++ {
		switch tabla[i].Type {
		case 'p':
			primaria++
		case 'l':
			logica++
		default:
			libre++
		}
	}

	//verificar la teoria de particiones
	//1 logica
	//para que hayan extendidas tiene que haber una logica
	//sumar 4

	if libre == 4 && nueva.Type != 'e' {
		//no hay particiones creadas aun
		nueva.Start = int64(unsafe.Sizeof(mbr{}))
		disco.Tabla[0] = *nueva
		return 1
	} else if libre == 4 && nueva.Type == 'e' {
		//ya hay particiones
		fmt.Print("RESULTADO: No se puede crear la particion extendida, debe crear una particion logica")
	}
	return 0

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
