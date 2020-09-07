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
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/**************************************************************
	Lee una linea por medio de teclado
***************************************************************/
func obtenerLineaConsola(texto string) {
	fmt.Print("COMANDO: ")
	lector := bufio.NewReader(os.Stdin)
	comando, _ := lector.ReadString('\n')
	comando = texto + comando

	if strings.Compare(comando, "") == 1 {
		arreglo := strings.Split(comando, "")
		if len(arreglo) > 0 {
			if strings.Compare(arreglo[0], "#") == 0 { //Es un comentario
				fmt.Println(comando)
				obtenerLineaConsola(texto)
			} else if strings.Contains(comando, "#") { //Quitar comentario entre lineas
				pos := strings.LastIndex(comando, "#")
				comando = comando[0:pos]
				verificarLineaConsola(strings.TrimSpace(comando))
			} else { //Puede ser un comando
				verificarLineaConsola(strings.TrimSpace(comando))
			}
		}
	}
	obtenerLineaConsola("")
}

/**************************************************************
	Verifica y recompone la linea que le fue entregada
	0 , comando		El comando necesita una nueva linea
	1 , comando		El comando esta completo y puede analizar
	3 , comanod		Se ingreso vacio, terminar
***************************************************************/
func verificar(comando string) (int, string) {
	if strings.Compare(comando, "\n") == 0 { //viene solo enter
		return 2, comando
	} else if strings.Compare(comando, "") == 0 { //viene vacio
		return 2, comando
	} else {

		comando = strings.ReplaceAll(comando, "\n", "")
		s := strings.Split(comando, " -") //Descompone cada -
		i := len(s) - 1
		s2 := strings.Fields(s[i]) //Descompone el ultimo -
		i2 := len(s2) - 1

		if i2 > -1 {
			if strings.Compare(s2[i2], "\\*") == 0 { //Pregunta si es \*
				nuevoComando := strings.ReplaceAll(strings.TrimSpace(comando), "\\*", "")
				return 0, nuevoComando
			}
		}
		return 1, strings.ReplaceAll(strings.TrimSpace(comando), "\n", "")
	}
}

/**************************************************************
	Determina si manda a analizar o a obtener nueva linea
***************************************************************/
func verificarLineaConsola(comando string) {
	i, linea := verificar(comando)
	if i == 0 {
		obtenerLineaConsola(linea)
	} else if i == 1 {
		analizar(comando)
	}
}

/**************************************************************
	Analiza y clasifica el comando ya formado
***************************************************************/
func analizar(comando string) {
	comando = strings.ReplaceAll(comando, "\n", "")
	if strings.Compare(comando, "") == 1 {
		s := strings.Split(comando, " -")
		switch strings.ToLower(s[0]) {
		case "exec":
			comandoExec(comando)
		case "pause":
			comandoPause(comando)
		case "mkdisk":
			comandoMkdisk(comando)
		case "rmdisk":
			comandoRmdisk(comando)
		case "fdisk":
			comandoFkdisk(comando)
		case "mount":
			comandoMount(comando)
		case "unmount":
			comandoUnmount(comando)
		case "mkfs":
			comandoMKFS(comando)
		case "mkdir":
			comandoMKDIR(comando)
		case "rep":
			comandoRep(comando)
		default:
			fmt.Println("La instruccion " + s[0] + " no se reconoce")
		}
	}
}

/**************************************************************
	COMANDO EXEC
	Obligatorio:
		-	Path
***************************************************************/
func comandoExec(comando string) {
	fmt.Println("EJECUTANDO: " + comando)
	s := strings.Split(comando, " -")
	if len(s) == 2 {
		s2 := strings.Split(s[1], "->")
		if strings.Compare(strings.ToLower(s2[0]), "path") == 0 {
			_, err := os.Stat(strings.ReplaceAll(s2[1], "\"", ""))
			if err == nil {
				s3 := strings.Split(s2[1], ".")
				if strings.Compare(s3[1], "mia") == 0 {
					fmt.Println("RESULTADO: Leyendo archivo...")
					fmt.Println("")
					archivo := leerArchivo(s2[1])
					//mandar a analizar ese archivo
					analizarArchivo(archivo)
				} else {
					fmt.Println("RESULTADO: La extension del archivo debe ser .MIA")
				}
			}
			if os.IsNotExist(err) {
				fmt.Println("RESULTADO: No existe el archivo especificado")
			}
		} else {
			fmt.Println("RESULTADO: El parametro PATH es obligatorio")
		}
	} else {
		fmt.Println("RESULTADO: Demasiados parametros para el comando EXEC")
	}
}

/**************************************************************
	Leer archivo, ruta ya validada
	Devuelve el contenido del archivo completo
***************************************************************/
func leerArchivo(ruta string) string {
	file, err := os.Open(ruta)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	//concatenando el contenido
	archivo := ""
	for fileScanner.Scan() {
		archivo = archivo + fileScanner.Text() + "\n"
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %s", err)
	}
	file.Close()
	//ya tengo el contenido del archivo
	return archivo
}

func analizarArchivo(contenido string) {
	s := strings.Split(contenido, "\n")
	comandoActual := ""
	for i := 0; i < len(s); i++ {
		if strings.Compare(s[i], "") == 0 {
			if strings.Compare(comandoActual, "") == 1 {

				analizar(strings.TrimSpace(comandoActual))
				comandoActual = ""
			}
		} else {
			arreglo := strings.Split(s[i], "")
			if len(arreglo) > 0 {
				if strings.Compare(arreglo[0], "#") == 0 { //Es un comentario
					fmt.Println(s[i])
				} else if strings.Contains(s[i], "#") { //Quitar comentario entre lineas
					pos := strings.LastIndex(s[i], "#")
					s[i] = s[i][0:pos]
					comandoActual = comandoActual + s[i]
					num, lin := verificarLineaArchivo(strings.TrimSpace(comandoActual))
					if num == 0 {
						comandoActual = lin //aqui duplicaba
					} else {
						comandoActual = ""
						analizar(strings.TrimSpace(lin))
					}
				} else { //Puede ser un comando
					comandoActual = comandoActual + s[i]
					num, lin := verificarLineaArchivo(strings.TrimSpace(comandoActual))
					if num == 0 {
						comandoActual = strings.TrimSpace(lin) + " "
					} else {
						comandoActual = ""
						analizar(strings.TrimSpace(lin))
					}
				}
			}
		}
	}
}

func verificarLineaArchivo(comando string) (int, string) {
	i, linea := verificar(comando)
	if i == 0 {
		return 0, strings.TrimSpace(linea) + " "
	} else if i == 1 {
		return 1, strings.TrimSpace(linea)
	}
	return 0, linea
}

/**************************************************************
	COMANDO PAUSE
***************************************************************/
func comandoPause(comando string) {
	fmt.Print("Presione la tecla Enter para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println()
}

/**************************************************************
	COMANDO MKDISK
		Obligatorios
			-size
			-path
			-name
		Opcionales
			-unit
***************************************************************/
func comandoMkdisk(comando string) {
	fmt.Println("EJECUTANDO: " + comando)
	//Descomponiendo en atributos
	atributos := strings.Split(comando, " -")
	//verificando parametros
	if len(atributos) > 3 {
		size := 0
		path := ""
		name := ""
		unit := ""
		for i := 1; i < len(atributos); i++ {
			atributoActual := strings.Split(atributos[i], "->")
			switch strings.ToLower(atributoActual[0]) {
			case "size":
				size = atributoSize(atributoActual[1])
			case "path":
				path = strings.ReplaceAll(atributoActual[1], "\"", "")
			case "name":
				name = strings.ReplaceAll(atributoActual[1], "\"", "")
			case "unit":
				unit = atributoUnit(atributoActual[1])
			default:
				fmt.Println("RESULTADO: El atributo " + atributoActual[0] + " no se reconoce")
			}
		}
		//verificando tamano
		if size < 1 {
			fmt.Println("RESULTADO: Error en el atributo SIZE")
		} else {
			//verificando unidad
			if strings.Compare(unit, "error") == 0 {
				fmt.Println("RESULTADO: Error en el atributo UNIT")
			} else if strings.Compare(unit, "") == 0 {
				unit = "m"
			} else {
				//verificando path
				if strings.Compare(path, "") == 1 {
					_, err := os.Stat(strings.ReplaceAll(path, "\"", ""))
					if err == nil {
					} else {
						crearDirectorioSiNoExiste(path)
					}
					//En este punto ya tiene que estar creado el directorio si es que no existia
				}
				//verificando nombre
				if strings.Compare(name, "") == 1 {
					_, err := os.Stat(strings.ReplaceAll(path+"/"+name, "\"", ""))
					if err == nil {
						fmt.Println("RESULTADO: El disco ya se encuentra creado, cambie de nombre")
						/*Esto lo tengo que quitar*/
						sName := strings.Split(name, ".")
						if strings.Compare(strings.ToLower(strings.TrimSpace(sName[1])), "dsk") == 0 {
							//Aqui mando a crear el archivo
							crearDisco(size, unit, path+"/"+name)
							fmt.Println("RESULTADO: Disco creado")
						} else {
							fmt.Println("RESULTADO: Solo se pueden crear discos con extension .DSK")
						}
						/*Esto lo tengo que quitar*/
					} else {
						//verificar extension
						sName := strings.Split(name, ".")
						if len(sName) > 1 {
							if strings.Compare(strings.ToLower(strings.TrimSpace(sName[1])), "dsk") == 0 {
								//Aqui mando a crear el archivo
								crearDisco(size, unit, path+"/"+name)
								fmt.Println("RESULTADO: Disco creado")
							} else {
								fmt.Println("RESULTADO: Solo se pueden crear discos con extension .DSK")
							}
						} else {
							fmt.Println("RESULTADO: Error en el nombre del disco a crear")
						}
					}
				}
			}
		}
	} else {
		fmt.Println("RESULTADO: Faltan atributos obligatorios para el comando MKDISK")
	}
}

func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.MkdirAll(directorio, 0777)
		if err != nil {
			//manejar el error aqui
		}
	}
}

/**************************************************************
	COMANDO RMDISK
	Obligatorio:
		-path
***************************************************************/
func comandoRmdisk(comando string) {
	fmt.Println("EJECUTANDO: " + comando)
	atributos := strings.Split(comando, " -")
	if len(atributos) > 1 {
		atributoActual := strings.Split(atributos[1], "->")
		if len(atributoActual) > 1 {
			if strings.Compare(strings.ToLower(atributoActual[0]), "path") == 0 {
				pathActual := strings.ReplaceAll(atributoActual[1], "\"", "")
				/********************/
				//verificando nombre
				if strings.Compare(pathActual, "") == 1 {
					_, err := os.Stat(pathActual)
					if err == nil {
						sName := strings.Split(pathActual, ".")
						if strings.Compare(strings.ToLower(strings.TrimSpace(sName[1])), "dsk") == 0 {
							//Aqui mando a borrar el disco
							removerDisco(pathActual)
							fmt.Println("RESULTADO: Disco eliminado")
						} else {
							fmt.Println("RESULTADO: Solo se pueden eliminar discos con extension .DSK")
						}
					} else {
						//no existe el archivo solicitado
						fmt.Println("RESULTADO: No existe el disco a eliminar")
					}
				}
				/*******************/
			}
		}
	} else {
		fmt.Println("RESULTADO: Faltan atributos obligatorios para el comando RMDISK")
	}
}

/**************************************************************
	COMANDO FDISK
***************************************************************/
func comandoFkdisk(comando string) {
	if strings.Compare(comando, "") == 1 {
		deleteFlag, addFlag := 0, 0
		fmt.Println("EJECUTANDO: " + comando)
		//Separa el primer comando general para determinar que accion realizar
		s := strings.Split(comando, " -")
		//verificar que hayan atributos
		if len(s) > 1 {
			for i := 0; i < len(s); i++ {
				s1 := strings.Split(s[i], "->")
				if len(s1) > 0 {
					if strings.Compare(strings.ToLower(s1[0]), "add") == 0 {
						addFlag = 1
					} else if strings.Compare(strings.ToLower(s1[0]), "delete") == 0 {
						deleteFlag = 1
					}
				}
			}
			if deleteFlag > 0 && addFlag > 0 {
				fmt.Println("RESULTADO: Existe una combinacion de instrucciones delete y add en el mismo comando FDISK")
			} else if deleteFlag > 0 {
				fDiskEliminar(comando)
			} else if addFlag > 0 {
				fDiskAdd(comando)
			} else {
				//si no es add ni delete puede que sea crear
				fDiskCrear(comando)
			}
		}
	}
}

/**************************************************************
	ELIMINAR PARTICION
	Obligatorios
	-name
	-path
	-tipo de eliminacion, va con el comando
***************************************************************/
func fDiskEliminar(comando string) {
	//verifico que el comando exista
	if strings.Compare(comando, "") != 0 {
		tipoEliminacion := ""
		nombreEliminacion := ""
		pathEliminacion := ""
		pathOk := 0
		s := strings.Split(comando, " -")
		if len(s) > 3 {
			for i := 1; i < len(s); i++ {
				s1 := strings.Split(s[i], "->")
				if len(s1) > 1 {
					switch strings.ToLower(strings.TrimSpace(s1[0])) {
					case "delete":
						tipoEliminacion = atributoDelete(s1[1])
					case "name":
						nombreEliminacion = strings.TrimSpace(strings.ReplaceAll(s1[1], "\"", ""))
					case "path":
						pathOk, pathEliminacion = verificarPath(s1[1])
					default:
						fmt.Println("RESULTADO: El atributo " + s1[0] + " no se reconoce para el comando DELETE")
					}
				}
			}
			if pathOk == 1 {
				//si existe el path
				if strings.Compare(tipoEliminacion, "error") != 0 {
					//si esta bien el tipo de eliminacion
					if strings.Compare(nombreEliminacion, "") != 0 {
						//si esta bien el nombre de la particion
						eliminarParticion(pathEliminacion, nombreEliminacion, tipoEliminacion)
					} else {
						fmt.Println("RESULTADO: Debe ingresar el nombre de la particion a eliminar")
					}
				} else {
					fmt.Println("RESULTADO: Error en el tipo de eliminacion de la particion")
				}
			} else if pathOk == 2 {
				fmt.Println("RESULTADO: El archivo indicado no representa un disco")
			} else {
				fmt.Println("RESULTADO: No existe el archivo especificado")
			}
		} else {
			fmt.Println("RESULTADO: Faltan atributos obligatorios para el comando FDISK DELETE")
		}
	}
}

/**************************************************************
	CREAR PARTICION
***************************************************************/
func fDiskCrear(comando string) {

	/********************************************
		Para crear una particion
		Obligatorios:
			-size
			-path
			-name
		Opcionales:
			-unit
			-type
			-fit
	*********************************************/

	//comando para ver si es una creacion
	if strings.Compare(comando, "") == 1 {
		//separando el comando " -"
		s := strings.Split(comando, " -")
		//verificando que s exista
		if len(s) > 0 {
			//verificando que venga el numero de atributos obligatorios
			atributosObligatorios := 3
			if len(s) > atributosObligatorios {
				//si cumple con los atributos obligatorios
				//verificar cada uno de los atributos
				size := 0
				pathOk := 0
				path := ""
				name := ""
				unit := "k"
				tipo := "p"
				fit := "wf"
				for i := 1; i < len(s); i++ {
					//dividiendo cada uno de los atributos en nombre y valor
					atributo := strings.Split(s[i], "->")
					//asegurando que atributo tenga unidades de informacion
					if len(atributo) > 1 {
						//verificando que atributo es
						switch strings.ToLower(atributo[0]) {
						case "size":
							size = atributoSize(atributo[1])
						case "path":
							pathOk, path = verificarPath(atributo[1])
						case "name":
							name = strings.ToLower(strings.ReplaceAll(atributo[1], "\"", ""))
						case "unit":
							unit = atributoUnitParticion(atributo[1])
						case "type":
							tipo = atributoType(atributo[1])
						case "fit":
							fit = atributoFit(atributo[1])
						default:
							fmt.Println("ADVERTENCIA: No se reconoce el atributo " + atributo[0])
						}
					} else {
						//si no tiene al menos dos unidades de informacion es un error
						fmt.Println("RESULTADO: Error en atributos del comando FDISK, falta informacion")
					}
				}
				//terminando el for para formar los atributos

				/*fmt.Println(size)
				fmt.Println(pathOk)
				fmt.Println(path)
				fmt.Println(name)
				fmt.Println(unit)
				fmt.Println(tipo)
				fmt.Println(fit)*/

				//verificando los parametros de creacion de la particion
				//verificando el tamano
				if size > 0 {
					//verificando path
					if pathOk == 0 {
						fmt.Println("RESULTADO: No existe el archivo indicado")
					} else if pathOk == 1 {
						//verificar el nombre
						if strings.Compare(name, "") == 1 {
							//veriicando unidad de tamano
							if strings.Compare(strings.ToLower(unit), "error") == 1 || strings.Compare(unit, "b") == 0 {
								//verificando tipo
								if strings.Compare(strings.ToLower(tipo), "error") == 1 || strings.Compare(tipo, "e") == 0 {
									//verificando ajuste
									if strings.Compare(strings.ToLower(fit), "error") == 1 || strings.Compare(fit, "bf") == 0 {
										//mandar a crear la particion
										crearParticion(path, size, unit, name, tipo, fit)
									} else {
										fmt.Println("RESULTADO: Error en ajuste de particion")
									}
								} else {
									fmt.Println("RESULTADO: Error en el tipo de particion a crear")
								}
							} else {
								fmt.Println("RESULTADO: Error en unidad de tamano de particion")
							}
						} else {
							fmt.Println("RESULTADO: El nombre de la particion no puede estar vacio")
						}
					} else if pathOk == 2 {
						fmt.Println("RESULTADO: El archivo indicado no representa un disco")
					}
				} else {
					fmt.Println("RESULTADO: Error en tamano de particion")
				}
			} else {
				//no cumple con los atributos obligatorios
				fmt.Println("RESULTADO: Faltan atributos obligatorios para el comando FDISK")
			}
		}
	}
}

/**************************************************************
	MODIFICAR TAMANO PARTICION
***************************************************************/
func fDiskAdd(comando string) {
	//verifico que el comando exista
	if strings.Compare(comando, "") != 0 {
		var tamAdd int = 0
		tipoAdd := ""
		nombreAdd := ""
		pathAdd := ""
		pathOk := 0
		unitAdd := "k"
		s := strings.Split(comando, " -")
		if len(s) > 3 {
			for i := 1; i < len(s); i++ {
				s1 := strings.Split(s[i], "->")
				if len(s1) > 1 {
					switch strings.ToLower(strings.TrimSpace(s1[0])) {
					case "add":
						tipoAdd, tamAdd = atributoAdd(s1[1])
					case "name":
						nombreAdd = strings.TrimSpace(strings.ReplaceAll(s1[1], "\"", ""))
					case "path":
						pathOk, pathAdd = verificarPath(s1[1])
					case "unit":
						unitAdd = atributoUnitParticion(s1[1])
					default:
						fmt.Println("RESULTADO: El atributo " + s1[0] + " no se reconoce para el comando ADD")
					}
				}
			}
			if pathOk == 1 {
				//si existe el path
				if strings.Compare(tipoAdd, "error") != 0 {
					//si esta bien el tipo de eliminacion
					if strings.Compare(nombreAdd, "") != 0 {
						//verificar la unidad
						if strings.Compare(unitAdd, "error") != 0 {
							switch unitAdd {
							case "k":
								tamAdd = tamAdd * 1024
							case "b":
							case "m":
								tamAdd = tamAdd * 1024 * 1024
							}
							//si esta bien el nombre de la particion
							modificarParticion(pathAdd, nombreAdd, int64(tamAdd), tipoAdd)
						} else {

						}
					} else {
						fmt.Println("RESULTADO: Debe ingresar el nombre de la particion a modificar")
					}
				} else {
					fmt.Println("RESULTADO: Error en el tipo de modificacion de la particion")
				}
			} else if pathOk == 2 {
				fmt.Println("RESULTADO: El archivo indicado no representa un disco")
			} else {
				fmt.Println("RESULTADO: No existe el archivo especificado")
			}
		} else {
			fmt.Println("RESULTADO: Faltan atributos obligatorios para el comando FDISK ADD")
		}
	}
}

/**************************************************************
	COMANDO MOUNT
	Obligatorios
	-path
	-name
***************************************************************/
func comandoMount(comando string) {
	if strings.Compare(comando, "") != 0 {
		fmt.Println("EJECUTANDO: " + comando)
		s := strings.Split(comando, " -")
		if len(s) == 1 {
			imprimirMOUNT()
		} else if len(s) == 3 {
			nombre := ""
			path := ""
			pathOk := 0
			for i := 1; i < len(s); i++ {
				s2 := strings.Split(s[i], "->")
				if len(s2) > 1 {
					switch strings.ToLower(strings.TrimSpace(s2[0])) {
					case "path":
						pathOk, path = verificarPath(s2[1])
					case "name":
						nombre = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(s2[1], "\"", "")))
					default:
						fmt.Println("RESULTADO: No se reconoce el comando " + s2[0] + " para el comando MOUNT")
					}
				}
			}
			if pathOk == 1 {
				if strings.Compare(nombre, "") != 0 {
					num := montarParticion(path, nombre)
					if num == 1 {
						fmt.Println("RESULTADO: Particion montada con exito")
					}
				} else {
					fmt.Println("RESULTADO: debe ingresar el nombre de la particion")
				}
			} else if pathOk == 0 {
				fmt.Println("RESULTADO: No existe la ruta especificada")
			} else if pathOk == 2 {
				fmt.Println("RESULTADO: El archivo especificado no representa un disco")
			}
		} else if len(s) < 3 {
			fmt.Println("RESULTADO: Faltan parametros obligatorios para el comando MOUNT")
		} else if len(s) > 3 {
			fmt.Println("RESULTADO: Se ingresaron mas parametros de los requeridos por el comando MOUNT")
		}
	}
}

/**************************************************************
	COMANDO UNMOUNT
***************************************************************/
func comandoUnmount(comando string) {
	if strings.Compare(comando, "") != 0 {
		fmt.Println("EJECUTANDO: " + comando)
		s := strings.Split(comando, " -")
		if len(s) > 1 {
			var listaID []string
			for i := 1; i < len(s); i++ {
				s2 := strings.Split(strings.ToLower(strings.TrimSpace(s[i])), "->")
				if len(s2) > 1 {
					/*
						Colocar un switch
						Aqui aun no he valudado que venga IDn u otra cosa pero para avanzar por el momento lo dejo asi
					*/
					listaID = append(listaID, s2[1])
				}
			}
			if len(listaID) > 0 {
				desmontar(listaID)
			} else {
				fmt.Println("RESULTADO: Debe ingresar el id de la(s) particion(es) a desmontar")
			}
		} else {
			fmt.Println("RESULTADO: Debe ingresar el id de la(s) particion(es) a desmontar")
		}
	}
}

/**************************************************************
	COMANDO MKFS
***************************************************************/
func comandoMKFS(comando string) {
	if strings.Compare(comando, "") != 0 {
		fmt.Println("EJECUTANDO: " + comando)
		s := strings.Split(comando, " -")
		addCounter := 0
		idFormatear := ""
		tipoFormato := ""
		tipoAdd := ""
		tamAdd := 0
		unitMkfs := ""
		if len(s) > 1 {
			for i := 1; i < len(s); i++ {
				s2 := strings.Split(strings.ToLower(strings.TrimSpace(s[i])), "->")
				if len(s2) > 1 {
					switch strings.ToLower(s2[0]) {
					case "id":
						idFormatear = strings.ToLower(strings.TrimSpace(s2[1]))
					case "type":
						tipoFormato = atributoTypeFormat(strings.ToLower(strings.TrimSpace(s2[1])))
					case "add":
						addCounter = 1
						tipoAdd, tamAdd = atributoAdd(strings.ToLower(strings.TrimSpace(s2[1])))
					case "unit":
						unitMkfs = atributoUnitParticion(strings.ToLower(strings.TrimSpace(s2[1])))
					}
				}
			}
			if addCounter == 1 {
				//es un add
				if strings.Compare(tipoAdd, "") == 0 || strings.Compare(tipoAdd, "error") == 0 {
					fmt.Println("RESULTADO: Error en el comando MKFS ADD")
				} else {
					if tamAdd <= 0 {
						fmt.Println("RESULTADO: Error en la cantidad a modificar en el comando MKFS ADD")
					} else {
						if strings.Compare(unitMkfs, "error") == 0 {
							fmt.Println("RESULTADO: Error en la unidad del comando MKFS ADD")
						} else {
							switch unitMkfs {
							case "m":
								tamAdd = tamAdd * 1024 * 1024
							case "k":
							case "b":
								tamAdd = tamAdd * 1024
							}
							if strings.Compare(idFormatear, "") == 0 {
								fmt.Println("RESULTADO: En el id a modificar en el comando MKFS ADD")
							} else {
								//todo esta bien
								mkfsAdd(idFormatear, int64(tamAdd), tipoAdd)
							}
						}
					}
				}
			} else {
				//no es format
				if strings.Compare(tipoAdd, "") == 0 {
					if strings.Compare(unitMkfs, "") == 0 {
						if strings.Compare(idFormatear, "") != 0 {
							if strings.Compare(tipoFormato, "error") != 0 {
								formatear(idFormatear, tipoFormato)
							} else {
								fmt.Println("RESULTADO: Error en el tipo de formato a realizar")
							}
						} else {
							fmt.Println("RESULTADO: Error en el ID de la particion a formatear 2")
						}
					} else {
						fmt.Println("RESULTADO: Parametro no permitido para MKFS Formatear 1 ")
					}
				} else {
					fmt.Println("RESULTADO: Parametro no permitido para MKFS Formatear")
				}
			}
		} else {
			fmt.Println("RESULTADO: Faltan parametros obligatorios para el comando MKFS")
		}
	}
}

func mkfsAdd(idParticion string, tamAdd int64, tipoAdd string) {
}

/**************************************************************
	COMANDO MKDIR
	Obligatorios
		-id (de la particion ya formateada)
		-path (DE LA CARPETA A CREAR)
	Opcional
		-p crear padres
***************************************************************/
func comandoMKDIR(comando string) {
	fmt.Println("EJECUTANDO: " + comando)
	if strings.Compare(comando, "") != 0 {
		s := strings.Split(comando, " -")
		atribP := 0
		atribPath := ""
		atribID := ""
		if len(s) > 2 {
			for i := 1; i < len(s); i++ {
				s2 := strings.Split(s[i], "->")
				if len(s2) > 0 {
					if len(s2) == 1 {
						//parametro p
						atribP = 1
					} else if len(s2) > 1 {
						switch strings.ToLower(strings.TrimSpace(s2[0])) {
						case "id":
							atribID = strings.ToLower(strings.TrimSpace(s2[1]))
						case "path":
							atribPath = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(s2[1], "\"", "")))
						default:
							fmt.Println("RESULTADO: Parametro no permitido para el comando MKDIR")
							return
						}
					}
				}
			}
			/*fmt.Println(atribP)
			fmt.Println(atribPath)
			fmt.Println(atribId)*/
			if strings.Compare(atribPath, "") != 0 {
				if strings.Compare(atribID, "") != 0 {
					crearDirectorio(atribID, atribPath, atribP)
				} else {
					fmt.Println("RESULTADO: Debe ingresar el id de la particion en la que desea crear la carpeta")
					return
				}
			} else {
				fmt.Println("RESULTADO: Debe ingresar la ruta de la carpeta a crear")
				return
			}
		} else {
			fmt.Println("RESULTADO: Faltan parametros obligatorios para el comando MKDIR")
		}
	}
}

///exec -path->/usr/local/go/src/archivos_proyecto1/archivo5.mia

/**************************************************************
	COMANDO REP
***************************************************************/
func comandoRep(comando string) {
	fmt.Println("EJECUTANDO: " + comando)
	id := ""
	path := ""
	ruta := ""
	nombre := ""
	if strings.Compare(comando, "") != 0 {
		s := strings.Split(comando, " -")
		if len(s) > 3 {
			for i := 1; i < len(s); i++ { //[0] = rep
				s2 := strings.Split(s[i], "->")
				if len(s2) > 1 {
					switch strings.ToLower(strings.TrimSpace(s2[0])) {
					case "path":
						path = strings.TrimSpace(s2[1])
					case "id":
						id = strings.TrimSpace(s2[1])
					case "ruta":
						ruta = strings.TrimSpace(s2[1])
					case "nombre":
						nombre = strings.TrimSpace(s2[1])
					default:
						fmt.Println("RESULTADO: No se reconoce el parametro " + s2[0] + " para el comando REP")
						return
					}
				}
			}
			/*revisar los parametros */
		} else {
			fmt.Println("RESULTADO: Faltan parametros obligatorios para el comando REP")
		}
	}
}

//exec -path->/usr/local/go/src/archivos_proyecto1/archivo5.mia

/**************************************************************
	Atributos
***************************************************************/
func atributoSize(cadena string) int {
	if strings.Compare(cadena, "") == 1 {
		i, err := strconv.Atoi(cadena)
		if err != nil {
			return 0
		}
		if i < 1 {
			return 0
		}
		return i
	}
	return 0
}

func atributoUnit(cadena string) string {
	if strings.Compare(cadena, "") == 1 {
		switch strings.ToLower(cadena) {
		case "k":
			return "k"
		case "m":
			return "m"
		default:
			return "error"
		}
	}
	return "error"
}

func atributoUnitParticion(cadena string) string {
	if strings.Compare(cadena, "") == 1 {
		switch strings.ToLower(cadena) {
		case "k":
			return "k"
		case "m":
			return "m"
		case "b":
			return "b"
		default:
			return "error"
		}
	}
	return "error"
}

func atributoFit(cadena string) string {
	if strings.Compare(cadena, "") == 1 {
		switch strings.ToLower(strings.TrimSpace(cadena)) {
		case "ff":
			return "ff"
		case "bf":
			return "bf"
		case "wf":
			return "wf"
		default:
			return "Error"
		}
	}
	return "Error"
}

func atributoType(cadena string) string {
	if strings.Compare(cadena, "") == 1 {
		switch strings.ToLower(strings.TrimSpace(cadena)) {
		case "p":
			return "p"
		case "e":
			return "e"
		case "l":
			return "l"
		default:
			return "error"
		}
	}
	return "error"
}

func verificarPath(pathActual string) (int, string) {
	pathActual = strings.ReplaceAll(pathActual, "\"", "")
	if strings.Compare(pathActual, "") == 1 {
		_, err := os.Stat(pathActual)
		if err == nil {
			sName := strings.Split(pathActual, ".")
			if strings.Compare(strings.ToLower(strings.TrimSpace(sName[1])), "dsk") == 0 {
				//Si existe el disco
				return 1, pathActual
			}
			//Existe el archivo pero no es un disco
			return 2, ""
		}
		//no existe el archivo
		return 0, ""
	}
	//se envio un path vacio
	return 0, ""
}

func atributoDelete(cadena string) string {
	if strings.Compare(cadena, "") == 0 {
		return "error"
	}
	switch strings.ToLower(strings.TrimSpace(strings.ReplaceAll(cadena, "\"", ""))) {
	case "fast":
		return "fast"
	case "full":
		return "full"
	default:
		return "error"
	}
}

func atributoAdd(cadena string) (string, int) {
	if strings.Compare(cadena, "") == 1 {
		i, err := strconv.Atoi(cadena)
		if err != nil {
			return "error", 0
		}
		if i < 0 {
			return "quitar", i * -1
		}
		return "agregar", i
	}
	return "error", 0
}

func atributoTypeFormat(cadena string) string {
	if strings.Compare(cadena, "") == 1 {
		switch strings.ToLower(strings.TrimSpace(cadena)) {
		case "fast":
			return "fast"
		case "full":
			return "full"
		default:
			return "error"
		}
	}
	return "error"
}
