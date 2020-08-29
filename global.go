package main

import (
	"fmt"
)

func tarea(tamParticion int64, inicioParticion int64) { //convertir a megas

	var tamAVD int64 = 20
	var tamDD int64 = 15
	var tamInodo int64 = 24
	var tamBloque int64 = 25
	var tamBitacora int64 = 18
	var tamSuperBloque int64 = 50

	var nEstructuras = (tamParticion - (2 * tamSuperBloque)) / (27 + tamAVD + tamDD + (5*tamInodo + (20 * tamBloque) + tamBitacora))

	cantidadAVD := nEstructuras
	cantidadDD := nEstructuras
	cantidadInodos := 5 * nEstructuras
	cantidadBloques := 20 * nEstructuras //5*cantidadInodos
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
	finalparticion := iniciocopiaSB + tamSuperBloque

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
}
