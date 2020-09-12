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
func graficarDISCO(path string, pathDestino string, nombreDestino string, formatoDestino string) {
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
	escribirDot(1, contenido, pathDestino, nombreDestino, formatoDestino)
}

/**************************************************************
	Grafica del MBR
***************************************************************/
func graficarMBR(path string, pathDestino string, nombreDestino string, formatoDestino string) {
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
					if particionActual.Name[indice] != 0 {
						numDetener = indice
					}
				}
				numDetener = numDetener + 1
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
						if ebrTemp.Size > 0 {
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
									contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#F08080\" width=\"250\" cellpadding=\"7\" align=\"left\"><font color=\"#FFFFFF\" face=\"Calibri\"><b>Particion Logica</b></font></TD>\n"
									contenido = contenido + "<TD border=\"0\" bgcolor=\"#F08080\" width=\"200\" cellpadding=\"7\"></TD></TR>\n"

									//status claro
									contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_status</font></TD>\n"
									contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Status)) + "</font></TD></TR>\n"
									//type oscuro
									contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#F5B7B1\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_next</font></TD>\n"
									contenido = contenido + "<TD border=\"0\" bgcolor=\"#F5B7B1\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Next)) + "</font></TD></TR>\n"
									//fit claro
									contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_fit</font></TD>\n"
									switch ebrLeido.Fit {
									case 'f':
										contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + "f" + "</font></TD></TR>\n"
									case 'b':
										contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + "b" + "</font></TD></TR>\n"
									case 'w':
										contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + "w" + "</font></TD></TR>\n"
									default:
										contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + "0" + "</font></TD></TR>\n"
									}
									//start oscuro
									contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#F5B7B1\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_start</font></TD>\n"
									contenido = contenido + "<TD border=\"0\" bgcolor=\"#F5B7B1\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Start)) + "</font></TD></TR>\n"
									//size claro
									contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_size</font></TD>\n"
									contenido = contenido + "<TD border=\"0\" bgcolor=\"#FFFFFF\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + strconv.Itoa(int(ebrLeido.Size)) + "</font></TD></TR>\n"
									//name oscuro
									numDetener := 0
									for indice := 0; indice < len(ebrLeido.Name); indice++ {
										if ebrLeido.Name[indice] != 0 {
											numDetener = indice
										}
									}
									numDetener = numDetener + 1
									contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"#F5B7B1\" width=\"250\" cellpadding=\"5\"><font color=\"#000000\" face=\"Calibri\">part_name</font></TD>\n"
									contenido = contenido + "<TD border=\"0\" bgcolor=\"#F5B7B1\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"#000000\" face=\"Calibri\">" + BytesToString(ebrLeido.Name[:numDetener]) + "</font></TD></TR>\n"

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
	escribirDot(2, contenido, pathDestino, nombreDestino, formatoDestino)
}

/**************************************************************
	Graficar superbloque
***************************************************************/
func graficarSB(path string, inicioParticion int64, pathDestino string, nombreDestino string, formatoDestino string) {
	sbTemp := superbloque{}
	numDetener := 0
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if &sbTemp != nil {

		colorClaro := "#FFFFFF"
		colorTitulo := "#145A32"
		colorOscuro := "#27AE60"
		contenido := "digraph G {\n" +
			"label = \"Reporte de SUPERBLOQUE\"\n" +
			"a0[label=<\n" +
			"<TABLE border=\"1\" cellspacing=\"0\" cellpadding=\"5\" bgcolor=\"#145A32\">\n"
		/*Titulo de la tabla*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorTitulo + "\" width=\"250\" cellpadding=\"7\" align=\"left\"><font color=\"#FFFFFF\" face=\"Calibri\"><b>Reporte de SUPERBLOQUE</b></font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorTitulo + "\" width=\"200\" cellpadding=\"7\"></TD></TR>\n"

		/*nombre claro*/
		numDetener = 0
		for indice := 0; indice < len(sbTemp.NombreHD); indice++ {
			if sbTemp.NombreHD[indice] != 0 {
				numDetener = indice
			}
		}
		numDetener = numDetener + 1
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_nombre_hd</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + BytesToString(sbTemp.NombreHD[:numDetener]) + "</font></TD></TR>\n"
		/*contador de avd*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_arbol_virtual_count</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.ArbolVirtualCount)) + "</font></TD></TR>\n"

		/*contador dd claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_detalle_directorio_count</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.DetalleDirectorioCount)) + "</font></TD></TR>\n"
		/*contador inodos oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_inodos_count</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InodosCount)) + "</font></TD></TR>\n"

		/*contador bloques claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_bloques_count</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.BloquesCount)) + "</font></TD></TR>\n"
		/*contador free avd oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_arbol_virtual_free</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.ArbolVirtualFree)) + "</font></TD></TR>\n"

		/*contador dd free claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_detalle_directorio_free</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.DetalleDirectorioFree)) + "</font></TD></TR>\n"
		/*contador inodos oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_inodos_free</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InodosFree)) + "</font></TD></TR>\n"

		/*contador bloques free claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_bloques_free</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.BloquesFree)) + "</font></TD></TR>\n"
		/*contador datecreation oscuro*/
		numDetener = 0
		for indice := 0; indice < len(sbTemp.DateCreacion); indice++ {
			if sbTemp.DateCreacion[indice] != 0 {
				numDetener = indice
			}
		}
		numDetener = numDetener + 1
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_date_creacion</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + BytesToString(sbTemp.DateCreacion[:numDetener]) + "</font></TD></TR>\n"

		/*ultimo montaje claro*/
		numDetener = 0
		for indice := 0; indice < len(sbTemp.DateUltimoMontaje); indice++ {
			if sbTemp.DateUltimoMontaje[indice] != 0 {
				numDetener = indice
			}
		}
		numDetener = numDetener + 1
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_date_ultimo_montaje</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + BytesToString(sbTemp.DateUltimoMontaje[:numDetener]) + "</font></TD></TR>\n"
		/*contador montajes count oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_montajes_count</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.MontajesCount)) + "</font></TD></TR>\n"

		/*contador ap bitmap avd claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_bitmap_arbol_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioBMAV)) + "</font></TD></TR>\n"
		/*contador ap_arboldirectorio oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_arbol_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioAV)) + "</font></TD></TR>\n"

		/*contador ap bitmap dd claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_bitmap_detalle_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioBMDD)) + "</font></TD></TR>\n"
		/*contador ap_dd oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_detalle_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioDD)) + "</font></TD></TR>\n"

		/*contador ap bitmap inodos claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_bitmap_inodos</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioBMInodos)) + "</font></TD></TR>\n"
		/*contador ap_inodos oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_inodos</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioInodos)) + "</font></TD></TR>\n"

		/*contador ap bitmap bloques claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_bitmap_bloques</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioBMBloques)) + "</font></TD></TR>\n"
		/*contador ap_bloques oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_bloques</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioBloques)) + "</font></TD></TR>\n"

		/*ap log claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_ap_log</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.InicioLog)) + "</font></TD></TR>\n"

		/*ap_size avd oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_size_struct_arbol_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(unsafe.Sizeof(avd{}))) + "</font></TD></TR>\n"
		/*ap size dd claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_size_struct_detalle_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(unsafe.Sizeof(dd{}))) + "</font></TD></TR>\n"
		/*ap_size inodo oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_size_struct_inodo</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(unsafe.Sizeof(inodo{}))) + "</font></TD></TR>\n"
		/*ap size bloque claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_size_struct_bloque</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(unsafe.Sizeof(bloque{}))) + "</font></TD></TR>\n"

		/*contador pimer bit arbol directorio oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_first_free_bit_arbol_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.PrimerLibreAV)) + "</font></TD></TR>\n"
		/*dd claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_first_free_bit_detalle_directorio</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.PrimerLibreDD)) + "</font></TD></TR>\n"
		/*inodos oscuro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_first_free_bit_tabla_inodos</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.PrimerLibreInodo)) + "</font></TD></TR>\n"
		/*bloques claro*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\"><font color=\"" + "#000000" + "\" face=\"Calibri\">sb_first_free_bit_bloques</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorClaro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.PrimerLibreBloque)) + "</font></TD></TR>\n"

		/*contador magic num*/
		contenido = contenido + "<TR><TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">sb_magic_num</font></TD>\n"
		contenido = contenido + "<TD border=\"0\" bgcolor=\"" + colorOscuro + "\" width=\"250\" cellpadding=\"5\" align=\"left\"><font  color=\"" + "#000000" + "\" face=\"Calibri\">" + strconv.Itoa(int(sbTemp.MagicNum)) + "</font></TD></TR>\n"

		//final
		contenido = contenido +
			"</TABLE>\n" +
			"> shape = \"rectangle\" fontcolor = \"black\"];\n" +
			"}\n"
		//escribir el archivo formado
		escribirDot(3, contenido, pathDestino, nombreDestino, formatoDestino)
	}
}

/**************************************************************
	Graficar BITMAP AVD
***************************************************************/
func graficarBitMapDirectorio(path string, inicioParticion int64, pathDestino string, nombreDestino string, formatoDestino string) {
	sbTemp := superbloque{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if &sbTemp != nil {
		if sbTemp.MagicNum == 201314821 {
			lineaContenido := ""
			counter := 0
			for i := sbTemp.InicioBMAV; i < sbTemp.InicioAV; i++ {
				file.Seek(i, 0)
				var n byte
				data := readNextBytes(file, unsafe.Sizeof(n))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &n)
				if err != nil {
					log.Fatal("binary.Read failed", err)
				}
				if n == 0 {
					lineaContenido = lineaContenido + "0 |"
				} else if n == 1 {
					lineaContenido = lineaContenido + "1 |"
				}
				if counter == 14 {
					counter = 0
					lineaContenido = lineaContenido + "\n"
				} else {
					counter = counter + 1
				}
			}
			crearArchivo(pathDestino+"/"+nombreDestino+"."+formatoDestino, lineaContenido)
		}
	}
}

/**************************************************************
	Graficar BITMAP DD
***************************************************************/
func graficarBitMapDetalle(path string, inicioParticion int64, pathDestino string, nombreDestino string, formatoDestino string) {
	sbTemp := superbloque{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if &sbTemp != nil {
		if sbTemp.MagicNum == 201314821 {
			lineaContenido := ""
			counter := 0
			for i := sbTemp.InicioBMDD; i < sbTemp.InicioDD; i++ {
				file.Seek(i, 0)
				var n byte
				data := readNextBytes(file, unsafe.Sizeof(n))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &n)
				if err != nil {
					log.Fatal("binary.Read failed", err)
				}
				if n == 0 {
					lineaContenido = lineaContenido + "0 |"
				} else if n == 1 {
					lineaContenido = lineaContenido + "1 |"
				}
				if counter == 14 {
					counter = 0
					lineaContenido = lineaContenido + "\n"
				} else {
					counter = counter + 1
				}
			}
			crearArchivo(pathDestino+"/"+nombreDestino+"."+formatoDestino, lineaContenido)
		}
	}
}

/**************************************************************
	Graficar BITMAP INODOS
***************************************************************/
func graficarBitMapInodo(path string, inicioParticion int64, pathDestino string, nombreDestino string, formatoDestino string) {
	sbTemp := superbloque{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if &sbTemp != nil {
		if sbTemp.MagicNum == 201314821 {
			lineaContenido := ""
			counter := 0
			for i := sbTemp.InicioBMInodos; i < sbTemp.InicioInodos; i++ {
				file.Seek(i, 0)
				var n byte
				data := readNextBytes(file, unsafe.Sizeof(n))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &n)
				if err != nil {
					log.Fatal("binary.Read failed", err)
				}
				if n == 0 {
					lineaContenido = lineaContenido + "0 |"
				} else if n == 1 {
					lineaContenido = lineaContenido + "1 |"
				}
				if counter == 14 {
					counter = 0
					lineaContenido = lineaContenido + "\n"
				} else {
					counter = counter + 1
				}
			}
			crearArchivo(pathDestino+"/"+nombreDestino+"."+formatoDestino, lineaContenido)
		}
	}
}

/**************************************************************
	Graficar BITMAP BLOQUES
***************************************************************/
func graficarBitMapBloque(path string, inicioParticion int64, pathDestino string, nombreDestino string, formatoDestino string) {
	sbTemp := superbloque{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if &sbTemp != nil {
		if sbTemp.MagicNum == 201314821 {
			lineaContenido := ""
			counter := 0
			for i := sbTemp.InicioBMBloques; i < sbTemp.InicioBloques; i++ {
				file.Seek(i, 0)
				var n byte
				data := readNextBytes(file, unsafe.Sizeof(n))
				buffer := bytes.NewBuffer(data)
				err = binary.Read(buffer, binary.BigEndian, &n)
				if err != nil {
					log.Fatal("binary.Read failed", err)
				}
				if n == 0 {
					lineaContenido = lineaContenido + "0 |"
				} else if n == 1 {
					lineaContenido = lineaContenido + "1 |"
				}
				if counter == 14 {
					counter = 0
					lineaContenido = lineaContenido + "\n"
				} else {
					counter = counter + 1
				}
			}
			crearArchivo(pathDestino+"/"+nombreDestino+"."+formatoDestino, lineaContenido)
		}
	}
}

/**************************************************************
	Graficar BITMAP BLOQUES
***************************************************************/
func graficarDirectorioGeneral(path string, inicioParticion int64, pathDestino string, nombreDestino string, formatoDestino string) {
	sbTemp := superbloque{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if &sbTemp != nil {
		if sbTemp.MagicNum == 201314821 {
			contenido := ""
			contenido = "digraph G {\n" +
			"label = \"Reporte de DIRECTORIO\"\nnode [shape=record penwidth = 3];"
			/*Formar los nodos*/
			contenido = contenido + "\n\n\n" + graficaDirectorioRecursiva(path, sbTemp.InicioAV) + "\n\n\n"
			/*Crear los apuntadores*/
			contenido = contenido + "}\n"
			escribirDot(4, contenido, pathDestino, nombreDestino, formatoDestino)
		}
	}
}

func graficaDirectorioRecursiva(path string, inicioactual int64) string {
	avdLeido := avd{}
	contenido := ""
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		//return ""
	}
	file.Seek(inicioactual, 0)
	data := readNextBytes(file, unsafe.Sizeof(avd{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &avdLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		//return ""
	}

	if &avdLeido != nil {
		if avdLeido.AVDApDetalleDirectorio != 0{
			///<f0>1|<f1>1|<f2>0|0|0|0|0|0
			numDetener := 0
			for indice := 0; indice < len(avdLeido.AVDNombreDirectorio); indice++ {
				if avdLeido.AVDNombreDirectorio[indice] != 0 {
					numDetener = indice
				}
			}
			numDetener = numDetener + 1
			copiaIn := inicioactual
			contenido = contenido + "node"+ strconv.Itoa(int(copiaIn)) +"[style=bold color=\"#9C640C\" fillcolor=\"#F4D03F\" style=\"filled\" label=\"{"+ BytesToString(avdLeido.AVDNombreDirectorio[:numDetener]) +" |{"
			/*recorrer los apuntadores*/
			for i:=0 ; i<len(avdLeido.AVDApArraySubdirectorios) ; i++{
				contenido = contenido + " <f"
				contenido = contenido + strconv.Itoa(int(i))
				contenido = contenido + ">"
				/*Aqui el contenido de sus subdirectorios*/
				if avdLeido.AVDApArraySubdirectorios[i] == 0{
					contenido = contenido + "0"
				}else {
					contenido = contenido + "1"
				}
				if i == len(avdLeido.AVDApArraySubdirectorios) -1{
					//contenido = contenido + "||" //esto seria para el completo 
					contenido = contenido + "|<f6>"
				}else{
					contenido = contenido + "|"
				}
			}
			contenido = contenido + "}}\"];\n"

			for i:=0 ; i<len(avdLeido.AVDApArraySubdirectorios) ; i++{
				if avdLeido.AVDApArraySubdirectorios[i] != 0 {
					contenido = contenido +"\n"+ graficaDirectorioRecursiva(path, avdLeido.AVDApArraySubdirectorios[i])
				}
			}

			if avdLeido.AVDApArbolVirtualDirectorio != 0{
				contenido = contenido + "\n" + graficaDirectorioRecursiva(path, avdLeido.AVDApArbolVirtualDirectorio)
			}
			/*Aqui hacer los enlaces*/

			//de la tabla
			for i := 0 ; i < len(avdLeido.AVDApArraySubdirectorios) ; i++{
				if avdLeido.AVDApArraySubdirectorios[i] != 0 {
					contenido = contenido + "\nnode"
					
					nuNodo :=  strconv.Itoa(int(inicioactual))
					contenido = contenido + nuNodo+ ":f" + strconv.Itoa(int(i)) + "->node"
					nuNodo2 :=  strconv.Itoa(int(avdLeido.AVDApArraySubdirectorios[i]))
					contenido = contenido + nuNodo2 + "[color=\"#641E16\"];"
				}
			}
			//con el extra
			if avdLeido.AVDApArbolVirtualDirectorio != 0{
				contenido = contenido + "\nnode"
					
				nuNodo :=  strconv.Itoa(int(inicioactual))
				contenido = contenido + nuNodo+ ":f6->node"
				nuNodo2 :=  strconv.Itoa(int(avdLeido.AVDApArbolVirtualDirectorio))
				contenido = contenido + nuNodo2 + "[color=\"#6F080C\"];"
			}
			return contenido 

		}
		return contenido
	}
	return ""
}

/**************************************************************
	Graficar Completo
***************************************************************/
func graficarCompleto(path string, inicioParticion int64, pathDestino string, nombreDestino string, formatoDestino string) {
	sbTemp := superbloque{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioParticion, 0)
	data := readNextBytes(file, unsafe.Sizeof(superbloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sbTemp)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if &sbTemp != nil {
		if sbTemp.MagicNum == 201314821 {
			contenido := ""
			contenido = "digraph G {\n" +
			"label = \"Reporte TREE COMPLETE\"\nnode [shape=record penwidth = 3];\nratio = auto;\n"
			/*Formar los nodos*/
			contenido = contenido + "\n\n\n" + graficaCompletoRecursiva(path, sbTemp.InicioAV) + "\n\n\n"
			/*Crear los apuntadores*/
			contenido = contenido + "}\n"
			escribirDot(4, contenido, pathDestino, nombreDestino, formatoDestino)
		}
	}
}

func graficaCompletoRecursiva(path string, inicioactual int64) string {
	avdLeido := avd{}
	contenido := ""
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		//return ""
	}
	file.Seek(inicioactual, 0)
	data := readNextBytes(file, unsafe.Sizeof(avd{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &avdLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		//return ""
	}

	if &avdLeido != nil {
		if avdLeido.AVDApDetalleDirectorio != 0{
			///<f0>1|<f1>1|<f2>0|0|0|0|0|0
			numDetener := 0
			for indice := 0; indice < len(avdLeido.AVDNombreDirectorio); indice++ {
				if avdLeido.AVDNombreDirectorio[indice] != 0 {
					numDetener = indice
				}
			}

			//fillcolor=\"#BB8FCE\" style=\"filled\"
			numDetener = numDetener + 1
			copiaIn := inicioactual
			contenido = contenido + "node"+ strconv.Itoa(int(copiaIn)) +"[style=bold color=\"#6F080C\" fillcolor=\"#D98880\" style=\"filled\"  label=\"{DIRECTORIO: "+ BytesToString(avdLeido.AVDNombreDirectorio[:numDetener]) +" |{"
			/*recorrer los apuntadores*/
			for i:=0 ; i<len(avdLeido.AVDApArraySubdirectorios) ; i++{
				contenido = contenido + " <f"
				contenido = contenido + strconv.Itoa(int(i))
				contenido = contenido + ">"
				/*Aqui el contenido de sus subdirectorios*/
				if avdLeido.AVDApArraySubdirectorios[i] == 0{
					contenido = contenido + "0"
				}else {
					contenido = contenido + "1"
				}
				if i == len(avdLeido.AVDApArraySubdirectorios) -1{
					//contenido = contenido + "||" //esto seria para el completo 
					contenido = contenido + "|<f6>▤|"
					if avdLeido.AVDApArbolVirtualDirectorio == 0{
						contenido = contenido + "<f7>□ "
					}else{
						contenido = contenido + "<f7>■"
					}
				}else{
					contenido = contenido + "|"
				}
			}
			contenido = contenido + "}}\"];\n"

			for i:=0 ; i<len(avdLeido.AVDApArraySubdirectorios) ; i++{
				if avdLeido.AVDApArraySubdirectorios[i] != 0 {
					contenido = contenido +"\n"+ graficaCompletoRecursiva(path, avdLeido.AVDApArraySubdirectorios[i])
				}
			}

			if avdLeido.AVDApArbolVirtualDirectorio != 0{
				contenido = contenido + "\n" + graficaCompletoRecursiva(path, avdLeido.AVDApArbolVirtualDirectorio)
			}
			/*Aqui hacer los enlaces*/

			//de la tabla
			for i := 0 ; i < len(avdLeido.AVDApArraySubdirectorios) ; i++{
				if avdLeido.AVDApArraySubdirectorios[i] != 0 {
					contenido = contenido + "\nnode"
					
					nuNodo :=  strconv.Itoa(int(inicioactual))
					contenido = contenido + nuNodo+ ":f" + strconv.Itoa(int(i)) + "->node"
					nuNodo2 :=  strconv.Itoa(int(avdLeido.AVDApArraySubdirectorios[i]))
					contenido = contenido + nuNodo2 + "[color=\"#14106C\"];"
				}
			}
			//con el extra
			if avdLeido.AVDApArbolVirtualDirectorio != 0{
				contenido = contenido + "\nnode"
					
				nuNodo :=  strconv.Itoa(int(inicioactual))
				contenido = contenido + nuNodo+ ":f7->node"
				nuNodo2 :=  strconv.Itoa(int(avdLeido.AVDApArbolVirtualDirectorio))
				contenido = contenido + nuNodo2 + "[color=\"#6F080C\"];"
			}

			/*Graficar los dd*/
			if avdLeido.AVDApDetalleDirectorio != 0 {
				/*Si tiene dd*/
				con1:= graficarDDs(path, avdLeido.AVDApDetalleDirectorio)
				if strings.Compare(con1, "") != 0{
					contenido = contenido + con1
					contenido = contenido + "\nnode"
					nuNodo :=  strconv.Itoa(int(inicioactual))
					contenido = contenido + nuNodo+ ":f6->node"
					nuNodo2 :=  strconv.Itoa(int(avdLeido.AVDApDetalleDirectorio))
					contenido = contenido + nuNodo2 + "[color=\"#196F3D\"];"
				}

			}


			return contenido 

		}
		return contenido
	}
	return ""
}

func graficarDDs(path string, inicioDD int64) string {
	ddLeido := dd{}
	contenido := ""
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		//return ""
	}
	file.Seek(inicioDD, 0)
	data := readNextBytes(file, unsafe.Sizeof(dd{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &ddLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		//return ""
	}

	if &ddLeido != nil {
		contadorArchivos := 0
		
	
		for i:= 0 ; i < len(ddLeido.DDArrayFiles); i++{
			if ddLeido.DDArrayFiles[i].FileApInodo != 0{
				contadorArchivos = contadorArchivos + 1
	
		}
		if contadorArchivos == 0{
			return ""
		}else{
			nuNodo :=  strconv.Itoa(int(inicioDD))
			contenido = contenido + "\nnode" + nuNodo + "[style=bold color=\"#370A19\" fillcolor=\"#BB8FCE\" style=\"filled\" label=\"{DETALLE|" 
			//{<f1>nombre de archivo}|{<f2>}|{<f3>}|{<f4>}|{<f5>}|{<f6>Aqui va el apuntador a otro }  ir armando
			for i:= 0 ; i < len(ddLeido.DDArrayFiles); i++{
				contenido = contenido + "{<f" + strconv.Itoa(int(i)) + ">"
				numDetener := 0
				for indice := 0; indice < len(ddLeido.DDArrayFiles[i].FileNombre); indice++ {
					if ddLeido.DDArrayFiles[i].FileNombre[indice] != 0 {
						numDetener = indice
					}
				}
				numDetener = numDetener
				comp := BytesToString(ddLeido.DDArrayFiles[i].FileNombre[:numDetener])
				contenido = contenido + comp 
				contenido = contenido + "}|"
			}

			

			
			if ddLeido.DDApDetalleDirectorio != 0{
				contenido = contenido + "{<f6>1}"
				contenido = contenido + "}\"];"
				con1:= graficarDDs(path, ddLeido.DDApDetalleDirectorio)
				if strings.Compare(con1, "")!=0{
					contenido = contenido + con1

					contenido = contenido + "\nnode"
					nuNodo :=  strconv.Itoa(int(inicioDD))
					contenido = contenido + nuNodo+ ":f6->node"
					nuNodo2 :=  strconv.Itoa(int(ddLeido.DDApDetalleDirectorio))
					contenido = contenido + nuNodo2 + "[color=\"#196F3D\"];"
				}
			}else{
				contenido = contenido + "{<f6>0}"
				contenido = contenido + "}\"];"
			}

			


			for i:= 0 ; i < len(ddLeido.DDArrayFiles); i++{
				if ddLeido.DDArrayFiles[i].FileApInodo != 0 {
					con2:= graficarInodos(path, ddLeido.DDArrayFiles[i].FileApInodo)
					contenido = contenido + con2
				}
			}
			for i:= 0 ; i < len(ddLeido.DDArrayFiles); i++{
				if ddLeido.DDArrayFiles[i].FileApInodo != 0 {
					contenido = contenido + "\nnode"
					nuNodo :=  strconv.Itoa(int(inicioDD))
					contenido = contenido + nuNodo+ ":f"+strconv.Itoa(int(i))+"->node"
					nuNodo2 :=  strconv.Itoa(int(ddLeido.DDArrayFiles[i].FileApInodo))
					contenido = contenido + nuNodo2 + "[color=\"#196F3D\"];"
				}
			}
			


			return contenido
		}
	}
	return ""
	}
	return ""
}


func graficarInodos(path string, inicioDD int64) string {
	inodoLeido := inodo{}
	contenido := ""
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		//return ""
	}
	file.Seek(inicioDD, 0)
	data := readNextBytes(file, unsafe.Sizeof(inodo{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &inodoLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		//return ""
	}

	if &inodoLeido != nil {
		
		nuNodo :=  strconv.Itoa(int(inicioDD))
		contenido = contenido + "\nnode" + nuNodo + "[style=bold color=\"#0B5345\" fillcolor=\"#2ECC71\" style=\"filled\" label=\"{INODO|" 
		//{<f1>nombre de archivo}|{<f2>}|{<f3>}|{<f4>}|{<f5>}|{<f6>Aqui va el apuntador a otro }  ir armando
		for i:= 0 ; i < len(inodoLeido.IArrayBloques); i++{
			contenido = contenido + "{<f" + strconv.Itoa(int(i)) + ">"
			
			comp := "Bloque " +  strconv.Itoa(int(i+1))
			contenido = contenido + comp 
			contenido = contenido + "}|"
		}
		
		if inodoLeido.IApIndirecto != 0{
			con1:= graficarInodos(path, inodoLeido.IApIndirecto)

			contenido = contenido + "{<f5>1}"
			contenido = contenido + "}\"];"
			
			if strings.Compare(con1, "")!=0{
				contenido = contenido + con1

				contenido = contenido + "\nnode"
				nuNodo :=  strconv.Itoa(int(inicioDD))
				contenido = contenido + nuNodo+ ":f5->node"
				nuNodo2 :=  strconv.Itoa(int(inodoLeido.IApIndirecto))
				contenido = contenido + nuNodo2 + "[color=\"#196F3D\"];\n"
			}
		}else{
			contenido = contenido + "{<f5>0}"
			contenido = contenido + "}\"];"
		}
		for i:= 0 ; i < len(inodoLeido.IArrayBloques); i++{
			if inodoLeido.IArrayBloques[i] != 0{
				contenido = contenido + graficarBloque(path, inodoLeido.IArrayBloques[i])
				contenido = contenido + "\nnode"
				nuNodo :=  strconv.Itoa(int(inicioDD))
				contenido = contenido + nuNodo+ ":f"+strconv.Itoa(int(i))+"->node"
				nuNodo2 :=  strconv.Itoa(int(inodoLeido.IArrayBloques[i]))
				contenido = contenido + nuNodo2 + "[color=\"#B03A2E\"];\n"
			}
		}


		return contenido
	}
	return ""
}

func graficarBloque(path string, inicio int64) string {
	bloqueLeido := bloque{}
	contenido := ""
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		//return ""
	}
	file.Seek(inicio, 0)
	data := readNextBytes(file, unsafe.Sizeof(bloque{}))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &bloqueLeido)
	if err != nil {
		log.Fatal("binary.Read failed", err)
		//return ""
	}

	if &bloqueLeido != nil {
		contenido = contenido + "\nnode"+strconv.Itoa(int(inicio))+" [style=bold color=\"#085E6F\" fillcolor=\"#5DADE2\" style=\"filled\"  label=\"{BLOQUE|{<f0>"
		numDetener := 0
		for indice := 0; indice < len(bloqueLeido.DBData); indice++ {
			if bloqueLeido.DBData[indice] != 0 {
				numDetener = indice
			}
		}
		numDetener = numDetener + 1
		comp := BytesToString(bloqueLeido.DBData[:numDetener])
		contenido = contenido + comp
		contenido = contenido +"}}\"];"
	}
	return contenido
}



/**************************************************************
	Metodo graficar general
***************************************************************/
func graficar(arg3 string, arg5 string, formato string) {
	arg0 := "/usr/bin/dot"
	arg1 := "-T" + formato
	arg4 := "-o"
	out := exec.Command(arg0, arg1, arg3, arg4, arg5)
	out.Run()
}

func escribirDot(tipo int, contenido string, pathDestino string, nombreDestino string, formatoDestino string) {
	switch tipo {
	case 1:
		crearArchivo("reportes/disco.dot", contenido)
		graficar("reportes/disco.dot", pathDestino+"/"+nombreDestino+"."+formatoDestino, formatoDestino)
	case 2:
		crearArchivo("reportes/mbr.dot", contenido)
		graficar("reportes/mbr.dot", pathDestino+"/"+nombreDestino+"."+formatoDestino, formatoDestino)
	case 3:
		crearArchivo("reportes/superbloque.dot", contenido)
		graficar("reportes/superbloque.dot", pathDestino+"/"+nombreDestino+"."+formatoDestino, formatoDestino)
	case 4:
		crearArchivo("reportes/directorio.dot", contenido)
		graficar("reportes/directorio.dot", pathDestino+"/"+nombreDestino+"."+formatoDestino, formatoDestino)
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
