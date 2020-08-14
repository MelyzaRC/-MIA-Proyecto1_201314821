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
	fmt.Println("Ejecutando: " + comando)
	s := strings.Split(comando, " -")
	if len(s) == 2 {
		s2 := strings.Split(s[1], "->")
		if strings.Compare(strings.ToLower(s2[0]), "path") == 0 {
			_, err := os.Stat(strings.ReplaceAll(s2[1], "\"", ""))
			if err == nil {
				s3 := strings.Split(s2[1], ".")
				if strings.Compare(s3[1], "mia") == 0 {
					fmt.Println("Resultado: Leyendo archivo...")
					fmt.Println("")
					archivo := leerArchivo(s2[1])
					//mandar a analizar ese archivo
					analizarArchivo(archivo)
				} else {
					fmt.Println("Resultado: La extension del archivo debe ser .MIA")
				}
			}
			if os.IsNotExist(err) {
				fmt.Println("Resultado: No existe el archivo especificado")
			}
		} else {
			fmt.Println("Resultado: El parametro PATH es obligatorio")
		}
	} else {
		fmt.Println("Resultado: Demasiados parametros para el comando EXEC")
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
	fmt.Println("Comando PAUSE")
	fmt.Print("Presione la tecla Enter para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
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
	fmt.Println("Ejecutando: " + comando)
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
				fmt.Println("Resultado: El atributo " + atributoActual[0] + " no se reconoce")
			}
		}
		//verificando tamano
		if size < 1 {
			fmt.Println("Resultado: Error en el atributo SIZE")
		} else {
			//verificando unidad
			if strings.Compare(unit, "error") == 0 {
				fmt.Println("Resultado: Error en el atributo UNIT")
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
						fmt.Println("Resultado: El disco ya se encuentra creado, cambie de nombre")
						/*Esto lo tengo que quitar*/
						sName := strings.Split(name, ".")
						if strings.Compare(strings.ToLower(strings.TrimSpace(sName[1])), "disk") == 0 {
							//Aqui mando a crear el archivo
							crearDisco(size, unit, path+"/"+name)
							fmt.Println("Resultado: Disco creado")
						} else {
							fmt.Println("Resultado: Solo se pueden crear discos con extension .DISK")
						}
						/*Esto lo tengo que quitar*/
					} else {
						//verificar extension
						sName := strings.Split(name, ".")
						if len(sName) > 1 {
							if strings.Compare(strings.ToLower(strings.TrimSpace(sName[1])), "disk") == 0 {
								//Aqui mando a crear el archivo
								crearDisco(size, unit, path+"/"+name)
								fmt.Println("Resultado: Disco creado")
							} else {
								fmt.Println("Resultado: Solo se pueden crear discos con extension .DISK")
							}
						} else {
							fmt.Println("Resultado: Error en el nombre del disco a crear")
						}
					}
				}
			}
		}
	} else {
		fmt.Println("Resultado: Faltan atributos obligatorios para el comando MKDISK")
	}
}

func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.MkdirAll(directorio, 0777)
		if err != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			panic(err)
		}
	}
}

/**************************************************************
	COMANDO RMDISK
***************************************************************/
func comandoRmdisk(comando string) {
	fmt.Println("Ejecutando: " + comando)
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
						if strings.Compare(strings.ToLower(strings.TrimSpace(sName[1])), "disk") == 0 {
							//Aqui mando a borrar el disco
							removerDisco(pathActual)
							fmt.Println("Resultado: Disco eliminado")
						} else {
							fmt.Println("Resultado: Solo se pueden eliminar discos con extension .DISK")
						}
						/*Esto lo tengo que quitar*/
					} else {
						//no existe el archivo solicitado
						fmt.Println("Resultado: No existe el disco a eliminar")
					}
				}
				/*******************/
			}
		}
	} else {
		fmt.Println("Resultado: Faltan atributos obligatorios para el comando RMDISK")
	}
}

///exec -path->/usr/local/go/src/archivos_proyecto1/archivo.mia
///exec -path->"/home/melyza/Escritorio/hola.jpg"
//exec -path->"/home/mis discos/Disco Prueba/Disco1.mia"
//mkdisk -SiZe->8 -pAth->”/home/mis discos/DISCO Prueba/” \*   #comentario
//-name->Disco1.dsk -uniT->k
/**************************************************************
	COMANDO FKDISK
***************************************************************/
func comandoFkdisk(comando string) {
	fmt.Println("Comando FDISK")
}

/**************************************************************
	COMANDO MOUNT
***************************************************************/
func comandoMount(comando string) {
	fmt.Println("Comando MOUNT")
}

/**************************************************************
	COMANDO UNMOUNT
***************************************************************/
func comandoUnmount(comando string) {
	fmt.Println("Comando UNMOUNT")
}

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

func validarRuta(cadena string) int {
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
