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
	"os/exec"
	"strings"
	"unsafe"
)

/**************************************************************
	Colores
***************************************************************/
//MBR 						#C0ECEB
//Particion Primaria		#F5AFDC
//Particion Logica			#A9FA52
//EBR						#4B8DF1
//Particion Extendida		#FCCA43
//Espacio Libre				#FFFFFF
//Fondo de tabla			#FEFDBD
/**************************************************************
	Grafica del mbr
***************************************************************/
func graficarMBR(path string) {
	contenido := "digraph G {\n" +
		"label = \"Estructura del disco\"\n" +
		"a0[label=<\n" +
		"<TABLE border=\"1\" cellspacing=\"0\" cellpadding=\"30\" bgcolor=\"#FEFDBD\">\n" +
		"<TR>\n"
	var disco *mbr = leerDisco(path)
	var inicioEspacio int64 = int64(unsafe.Sizeof(mbr{}))
	var finalAnterior int64 = inicioEspacio
	if disco != nil {
		//colocar el MBR
		contenido = contenido + "<TD border=\"1\" bgcolor=\"#C0ECEB\"><b>MBR</b></TD>\n"
		//formar el contenido

		//determinando el part_start
		for i := 0; i < 4; i++ {
			if disco.Tabla[i].Size > 0 {
				if disco.Tabla[i].Start-finalAnterior > 1 {
					contenido = contenido + "<TD border=\"1\"  bgcolor=\"#FFFFFF \"><b>Libre</b></TD>\n"
					switch disco.Tabla[i].Type {
					case 'p':
						contenido = contenido + "<TD border=\"1\"  bgcolor=\"#F5AFDC\"><b>Primaria</b></TD>\n"
					case 'e':
						//aqui graficar las logicas
						contenido = contenido + "<TD border=\"1\"  bgcolor=\"#A9FA52\" cellpadding=\"5\">"
						/*Grafica de la parte logica****************************/
						ebrTemp := ebr{}
						file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
						defer file.Close()
						if err != nil {
							log.Fatal(err)
						} else {
							file.Seek(0, int(disco.Tabla[i].Start))
							data := readNextBytes(file, unsafe.Sizeof(ebr{}))
							buffer := bytes.NewBuffer(data)
							err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
							if err != nil {
								log.Fatal("binary.Read failed", err)
							}
						}
						if ebrTemp.Next == -1 && ebrTemp.Size == 0 {

							contenido = contenido + "<TABLE border=\"1\" cellspacing=\"0\" cellpadding=\"30\" bgcolor=\"yellow\">\n" +
								"<TR>\n" +
								"<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
								"<TD border=\"1\"  bgcolor=\"#FFFFFF\"><b>Libre</b></TD>\n" +
								"</TR>\n" +
								"</TABLE>\n"
						} else {
							//Aqui graficar los demas ebr
						}
						/*******************************************************/
						contenido = contenido + "</TD>\n"
					}
					finalAnterior = disco.Tabla[i].Start + disco.Tabla[i].Size - 1
				} else {
					switch disco.Tabla[i].Type {
					case 'p':
						contenido = contenido + "<TD border=\"1\"  bgcolor=\"#F5AFDC\"><b>Primaria</b></TD>\n"
					case 'e':
						//aqui graficar las logicas
						contenido = contenido + "<TD border=\"1\"  bgcolor=\"#A9FA52\" cellpadding=\"5\">"
						/*Grafica de la parte logica****************************/
						ebrTemp := ebr{}
						file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
						defer file.Close()
						if err != nil {
							log.Fatal(err)
						} else {
							file.Seek(0, 0)
							file.Seek(int64(disco.Tabla[i].Start), 0)
							data := readNextBytes(file, unsafe.Sizeof(ebr{}))
							buffer := bytes.NewBuffer(data)
							err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
							if err != nil {
								log.Fatal("binary.Read failed", err)
							}
						}
						if ebrTemp.Next == -1 && ebrTemp.Size == 0 {
							contenido = contenido + "<TABLE border=\"1\" cellspacing=\"0\" cellpadding=\"30\" bgcolor=\"yellow\">\n" +
								"<TR>\n" +
								"<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
								"<TD border=\"1\"  bgcolor=\"#FFFFFF\"><b>Libre</b></TD>\n" +
								"</TR>\n" +
								"</TABLE>\n"
						} else {
							//Aqui graficar los demas ebr
						}
						/*******************************************************/
						contenido = contenido + "</TD>\n"
					}
					finalAnterior = disco.Tabla[i].Start + disco.Tabla[i].Size - 1
				}
			}
		}

	}
	if disco.Tamano-finalAnterior > 1 {
		contenido = contenido + "<TD border=\"1\"  bgcolor=\"#FFFFFF \"><b>Libre</b></TD>\n"
	}
	contenido = contenido + "</TR>\n" +
		"</TABLE>\n" +
		"> shape = \"rectangle\" fontcolor = \"black\"];\n" +
		"}\n"
	//escribir el archivo formado
	escribirDot(1, contenido)
}

/**************************************************************
	Metodo graficar general
***************************************************************/
func graficar(arg3 string, arg5 string) {
	arg0 := "/usr/bin/dot"
	arg1 := "-Tpng"
	arg4 := "-o"
	out := exec.Command(arg0, arg1, arg3, arg4, arg5)
	out.Run()
}

func escribirDot(tipo int, contenido string) {
	switch tipo {
	case 1:
		crearArchivo("reportes/mbr.dot", contenido)
		graficar("reportes/mbr.dot", "reportes/mbr.png")
	default:
	}
}

func crearArchivo(path string, contenido string) {
	//Verifica que el archivo existe
	var _, err = os.Stat(path)
	//Crea el archivo si no existe
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if existeError(err) {
			return
		}
		defer file.Close()
	}
	escribeArchivo(path, contenido)
}

func escribeArchivo(path string, contenido string) {
	// Abre archivo usando permisos READ & WRITE
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if existeError(err) {
		return
	}
	defer file.Close()
	// Escribe algo de texto linea por linea
	_, err = file.WriteString(contenido)
	if existeError(err) {
		return
	}
	// Salva los cambios
	err = file.Sync()
	if existeError(err) {
		return
	}
}

func existeError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}
