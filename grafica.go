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
	"strconv"
	"strings"
	"unsafe"
)

/**************************************************************
	Colores
***************************************************************/

/*****************PARA LA ESTRUCTURA DEL DISCO*****************/
//			MBR 						#A3E4D7
//			Particion Primaria			#D7BDE2
//			Particion Extendida			#1E8449
//			EBR							#4B8DF1
//			Particion Logica			#D68910
//			Espacio Libre				#FFFFFF
//			Fondo de tabla				#FEFDBD

/**************************PARA EL MBR*************************/
//			Titulo 						#4A235A
//			Celdas						#E8DAEF

/**************************************************************
	Grafica la estructura del disco
***************************************************************/
func graficarDISCO(path string) {
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
		contenido = contenido + "<TD border=\"1\" bgcolor=\"#F7DC6F\"><b>MBR</b></TD>\n"
		//formar el contenido

		//determinando el part_start
		for i := 0; i < 4; i++ {
			if disco.Tabla[i].Size > 0 {
				if disco.Tabla[i].Start-finalAnterior > 1 {
					contenido = contenido + "<TD border=\"1\"  bgcolor=\"#BFC9CA\"><b>Libre</b></TD>\n"
				}
				switch disco.Tabla[i].Type {
				case 'p':
					contenido = contenido + "<TD border=\"1\"  bgcolor=\"#D7BDE2\"><b>Primaria</b></TD>\n"
				case 'e':
					contenido = contenido + "<TD border=\"1\"  bgcolor=\"#2ECC71\" cellpadding=\"5\">\n"
					//aqui graficar las logicas************************************
					contenido = contenido + "<TABLE border=\"1\" cellspacing=\"0\" cellpadding=\"5\" bgcolor=\"black\">\n" +
						"<TR>\n"

					//leyendo el archivo
					ebrTemp := ebr{}
					file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
					defer file.Close()
					if err != nil {
						log.Fatal(err)
					}
					file.Seek(disco.Tabla[i].Start, 0)
					data := readNextBytes(file, unsafe.Sizeof(ebr{}))
					buffer := bytes.NewBuffer(data)
					err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
					if err != nil {
						log.Fatal("binary.Read failed", err)
					}
					limite := disco.Tabla[i].Start + int64(unsafe.Sizeof(ebr{})) + disco.Tabla[i].Size
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
								if ebrLeido.Next != -1 && ebrLeido.Size == 0 {
									contenido = contenido + "<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
										"<TD border=\"1\"  bgcolor=\"#BFC9CA\"><b>Libre</b></TD>\n"
									i = ebrLeido.Next - 1
								} else if ebrLeido.Next == -1 && ebrLeido.Size == 0 {
									contenido = contenido + "<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
										"<TD border=\"1\"  bgcolor=\"#BFC9CA\"><b>Libre</b></TD>\n"
									i = limite + 1
								} else if ebrLeido.Next == -1 { //lego al utimo ebr
									disponible := limite - ebrLeido.Start - ebrLeido.Size - int64(unsafe.Sizeof(ebr{}))
									if disponible > 0 {
										contenido = contenido + "<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
											"<TD border=\"1\"  bgcolor=\"#D68910\"><b>Logica</b></TD>\n" +
											"<TD border=\"1\"  bgcolor=\"#BFC9CA\"><b>Libre</b></TD>\n"
										i = limite + 1
									} else {
										contenido = contenido + "<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
											"<TD border=\"1\"  bgcolor=\"#D68910\"><b>Logica</b></TD>\n"
										i = limite + 1
									}
								} else if ebrLeido.Next != -1 { //esta en los ebr antes del ultimo
									//verificar pero con el next
									disponible := ebrLeido.Next - ebrLeido.Start - ebrLeido.Size - int64(unsafe.Sizeof(ebr{}))
									if disponible > 0 {
										contenido = contenido + "<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
											"<TD border=\"1\"  bgcolor=\"#D68910\"><b>Logica</b></TD>\n" +
											"<TD border=\"1\"  bgcolor=\"#BFC9CA\"><b>Libre</b></TD>\n"
									} else {
										contenido = contenido + "<TD border=\"1\"  bgcolor=\"#4B8DF1\"><b>EBR</b></TD>\n" +
											"<TD border=\"1\"  bgcolor=\"#D68910\"><b>Logica</b></TD>\n"
									}
									//porque al iterar el for le suma uno
									i = ebrLeido.Next - 1
								}

							}

						}
					}
					//lo que voy copiando
					contenido = contenido + "</TR>\n" +
						"</TABLE>\n"
					//aqui graficar las logicas fin********************************

					contenido = contenido + "</TD>\n"
				}
				finalAnterior = disco.Tabla[i].Start + disco.Tabla[i].Size - 1
			}
		}

	}
	if disco.Tamano-finalAnterior > 1 {
		contenido = contenido + "<TD border=\"1\"  bgcolor=\"#BFC9CA\"><b>Libre</b></TD>\n"
	}
	contenido = contenido + "</TR>\n" +
		"</TABLE>\n" +
		"> shape = \"rectangle\" fontcolor = \"black\"];\n" +
		"}\n"
	//escribir el archivo formado
	escribirDot(1, contenido)
}

/**************************************************************
	Grafica del MBR
***************************************************************/
func graficarMBR(path string) {
	contenido := "digraph G {\n" +
		"label = \"Reporte de MBR\"\n" +
		"a0[label=<\n" +
		"<TABLE border=\"1\" cellspacing=\"0\" cellpadding=\"5\" bgcolor=\"#FEFDBD\">\n"
		//"<TR>\n"
	var disco *mbr = leerDisco(path)
	//ar inicioEspacio int64 = int64(unsafe.Sizeof(mbr{}))
	//var finalAnterior int64 = inicioEspacio
	if disco != nil {
		//colocar el MBR
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#4A235A\" width=\"250\" cellpadding=\"7\" align=\"left\"><font color=\"#FFFFFF\" face=\"Calibri\"><b>REPORTE DE MBR</b></font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"#4A235A\" width=\"200\" cellpadding=\"7\"></TD></TR>\n"
		//tamano claro
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">mbr_tamano</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(disco.Tamano)) + "</font></TD></TR>\n"
		//creacion oscuro
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">mbr_fecha_creacion</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + BytesToString(disco.Fecha[:]) + "</font></TD></TR>\n"
		//signature claro
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">mbr_disk_signature</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(disco.Firma)) + "</font></TD></TR>\n"

		//particiones
		for i := 0; i < len(disco.Tabla); i++ {
			if disco.Tabla[i].Type != 0 {
				particionActual := disco.Tabla[i]
				//particionActual := disco.Tabla[i]
				//titulo de tabla
				contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#4A235A\" width=\"250\" cellpadding=\"7\" align=\"left\"><font color=\"#FFFFFF\" face=\"Calibri\"><b>Particion</b></font></TD>\n"
				contenido = contenido + "<TD border=\"0\" bgcolor=\"#4A235A\" width=\"200\" cellpadding=\"7\"></TD></TR>\n"

				//status claro
				contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_status</font></TD>\n"
				contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(particionActual.Status)) + "</font></TD></TR>\n"
				//type oscuro
				contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_type</font></TD>\n"
				contenido = contenido + "<TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + string(particionActual.Type) + "</font></TD></TR>\n"
				//fit claro
				contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_fit</font></TD>\n"
				contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + string(particionActual.Fit) + "</font></TD></TR>\n"
				//start oscuro
				contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_start</font></TD>\n"
				contenido = contenido + "<TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(particionActual.Start)) + "</font></TD></TR>\n"
				//size claro
				contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_size</font></TD>\n"
				contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(particionActual.Size)) + "</font></TD></TR>\n"
				//name oscuro
				numDetener := 0
				for indice := 0; indice < len(particionActual.Name); indice++ {
					if particionActual.Name[indice] == 0 {
						numDetener = indice
						break
					}
				}
				contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_name</font></TD>\n"
				contenido = contenido + "<TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + BytesToString(particionActual.Name[:numDetener]) + "</font></TD></TR>\n"

				/*Logicas**************************************************************/
				if particionActual.Type == 'e' {
					ebrTemp := ebr{}
					file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
					defer file.Close()
					if err != nil {
						log.Fatal(err)
					}
					file.Seek(disco.Tabla[i].Start, 0)
					data := readNextBytes(file, unsafe.Sizeof(ebr{}))
					buffer := bytes.NewBuffer(data)
					err = binary.Read(buffer, binary.BigEndian, &ebrTemp)
					if err != nil {
						log.Fatal("binary.Read failed", err)
					}
					limite := particionActual.Start + int64(unsafe.Sizeof(ebr{})) + particionActual.Size
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
								//titulo de tabla
								contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#4A235A\" width=\"250\" cellpadding=\"7\" align=\"left\"><font color=\"#FFFFFF\" face=\"Calibri\"><b>Particion Logica</b></font></TD>\n"
								contenido = contenido + "<TD border=\"0\" bgcolor=\"#4A235A\" width=\"200\" cellpadding=\"7\"></TD></TR>\n"

								//status claro
								contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_status</font></TD>\n"
								contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Status)) + "</font></TD></TR>\n"
								//type oscuro
								contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_next</font></TD>\n"
								contenido = contenido + "<TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Next)) + "</font></TD></TR>\n"
								//fit claro
								contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_fit</font></TD>\n"
								contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + string(ebrLeido.Fit) + "</font></TD></TR>\n"
								//start oscuro
								contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_start</font></TD>\n"
								contenido = contenido + "<TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Start)) + "</font></TD></TR>\n"
								//size claro
								contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_size</font></TD>\n"
								contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Size)) + "</font></TD></TR>\n"
								//name oscuro
								numDetener := 0
								for indice := 0; indice < len(ebrLeido.Name); indice++ {
									if ebrLeido.Name[indice] == 0 {
										numDetener = indice
										break
									}
								}
								contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_name</font></TD>\n"
								contenido = contenido + "<TD border=\"0\" bgcolor=\"#E8DAEF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + BytesToString(ebrLeido.Name[:numDetener]) + "</font></TD></TR>\n"

								if ebrLeido.Next != -1 && ebrLeido.Size == 0 {
									i = ebrLeido.Next - 1
								} else if ebrLeido.Next == -1 && ebrLeido.Size == 0 {
									i = limite + 1
								} else if ebrLeido.Next == -1 { //lego al utimo ebr
									i = limite + 1
								} else if ebrLeido.Next != -1 { //esta en los ebr antes del ultimo
									//porque al iterar el for le suma uno
									i = ebrLeido.Next - 1
								}
							}

						}
					}
				}
				/*Logicas**************************************************************/
			}
		}
	}
	contenido = contenido +
		"</TABLE>\n" +
		"> shape = \"rectangle\" fontcolor = \"black\"];\n" +
		"}\n"
	//escribir el archivo formado
	escribirDot(2, contenido)
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
		crearArchivo("reportes/disco.dot", contenido)
		graficar("reportes/disco.dot", "reportes/disco.png")
	case 2:
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
