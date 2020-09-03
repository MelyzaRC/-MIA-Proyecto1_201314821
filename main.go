/**************************************************************
	Melyza Alejandra Rodriguez Contreras
	201314821
	Laboratorio de Manejo e implementacion de Archivos
	Segundo Semestre 2020
	Proyecto No. 1
***************************************************************/
package main

func main() {
	obtenerLineaConsola("")
}

/*
digraph structs {
    node [shape=record];



    carpetaRaiz [style=bold color="#6F080C" label="{ / |{<f0>1|<f1>1|<f2>0|0|0|0|0|0}}"];


    carpetaHome[style=bold color="#6F080C" label="{<f0>/home |{<f1>1|<f2>1|<f3>1|<f4>1|<f5>1|<f6>1|<f7>0|<f8>1}}"];
    carpetaHome2[style=bold color="#6F080C" label="{<f0>/home |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];

    carpeta1 [style=bold color="#6F080C" label="{<f0>/carpeta1 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpeta2 [ style=bold color="#6F080C" label="{<f0>/carpeta2 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpeta3 [style=bold color="#6F080C" label="{<f0>/carpeta3 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpeta4 [style=bold color="#6F080C" label="{<f0>/carpeta4 |{<f1>1|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpeta4_1 [style=bold color="#6F080C" label="{<f0>/carpeta4.1 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>1|<f8>0}}"];
    carpeta5 [style=bold color="#6F080C" label="{<f0>/carpeta5 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpeta6 [style=bold color="#6F080C" label="{<f0>/carpeta6 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpeta7 [style=bold color="#6F080C" label="{<f0>/carpeta7 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpeta8 [style=bold color="#6F080C" label="{<f0>/carpeta8 |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];

    carpetaRaiz:f0 -> carpetaHome:f0 [color="#14106C"];

    carpetaHome:f1 -> carpeta1:f0 [color="#14106C"];
    carpetaHome:f2 -> carpeta2:f0 [color="#14106C"];
    carpetaHome:f3 -> carpeta3:f0 [color="#14106C"];
    carpetaHome:f4 -> carpeta4:f0 [color="#14106C"];
    carpetaHome:f5 -> carpeta5:f0 [color="#14106C"];
    carpetaHome:f6 -> carpeta6:f0 [color="#14106C"];
    carpetaHome:f8 -> carpetaHome2:f0 [color="#14106C"];
    carpetaHome2:f1 -> carpeta7:f0 [color="#14106C"];
    carpetaHome2:f2 -> carpeta8:f0 [color="#14106C"];
    carpeta4:f1 -> carpeta4_1:f0 [color="#14106C"];

    detalle4_1 [style=bold color="#370A19" label="{<f0>DETALLE|{hola.txt|<f1>}|{|<f2>}|{|<f3>}|{|<f4>}|{|<f5>}|{<f6>0}}"];
    carpeta4_1:f7 -> detalle4_1:f0 [color="#14106C"];

    inodo1 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>1}|{Bloque 3|<f3>1}|{Bloque 4|<f4>1}|{<f5>1}}"];
    inodo2 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>0}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalle4_1:f1 -> inodo1:f0
    inodo1:f5 -> inodo2:f0

    data1 [style=bold color="#085E6F" label="{<f0>Data}"];
    data2 [style=bold color="#085E6F" label="{<f0>Data}"];
    data3 [style=bold color="#085E6F" label="{<f0>Data}"];
    data4 [style=bold color="#085E6F" label="{<f0>Data}"];
    data5 [style=bold color="#085E6F" label="{<f0>Data}"];
    data6 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodo1:f1 -> data1:f0 [color="#14106C"];
    inodo1:f2 -> data2:f0 [color="#14106C"];
    inodo1:f3 -> data3:f0 [color="#14106C"];
    inodo1:f4 -> data4:f0 [color="#14106C"];

    inodo2:f1 -> data5:f0 [color="#14106C"];
    inodo2:f2 -> data6:f0 [color="#14106C"];


    carpetaUser [style=bold color="#6F080C" label="{<f0>/user |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpetaRaiz:f1 -> carpetaUser:f0 [color="#14106C"];

    detalleUser [style=bold color="#370A19" label="{<f0>DETALLE|{file1.txt|<f1>}|{file2.txt|<f2>}|{file3.txt|<f3>}|{file4.txt|<f4>}|{file5.txt|<f5>}|{<f6>1}}"];
    carpetaUser:f7 -> detalleUser:f0 [color="#14106C"];
    detalleUser2 [style=bold color="#370A19" label="{<f0>DETALLE|{file6.txt|<f1>}|{file7.txt|<f2>}|{|<f3>}|{|<f4>}|{|<f5>}|{<f6>0}}"];
    detalleUser:f6 -> detalleUser2:f0 [color="#14106C"];

    inodoFile1 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>1}|{Bloque 3|<f3>1}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleUser:f1 -> inodoFile1:f0 [color="#14106C"];
    data1File1 [style=bold color="#085E6F" label="{<f0>Data}"];
    data2File1 [style=bold color="#085E6F" label="{<f0>Data}"];
    data3File1 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoFile1:f1 -> data1File1:f0 [color="#14106C"];
    inodoFile1:f2 -> data2File1:f0 [color="#14106C"];
    inodoFile1:f3 -> data3File1:f0 [color="#14106C"];

    inodoFile2 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleUser:f2 -> inodoFile2:f0 [color="#14106C"];
    data1File2 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoFile2:f1 -> data1File2:f0 [color="#14106C"];

    inodoFile3 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleUser:f3 -> inodoFile3:f0 [color="#14106C"];
    data1File3 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoFile3:f1 -> data1File3:f0 [color="#14106C"];

    inodoFile4 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleUser:f4 -> inodoFile4:f0 [color="#14106C"];
    data1File4 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoFile4:f1 -> data1File4:f0 [color="#14106C"];

    inodoFile5 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleUser:f5 -> inodoFile5:f0 [color="#14106C"];
    data1File5 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoFile5:f1 -> data1File5:f0 [color="#14106C"];

    inodoFile6 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleUser2:f1 -> inodoFile6:f0 [color="#14106C"];
    data1File6 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoFile6:f1 -> data1File6:f0 [color="#14106C"];

    inodoFile7 [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleUser2:f2 -> inodoFile7:f0 [color="#14106C"];
    data1File7 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoFile7:f1 -> data1File7:f0 [color="#14106C"];


    carpetaBin [style=bold color="#6F080C" label="{<f0>/bin |{<f1>0|<f2>0|<f3>0|<f4>0|<f5>0|<f6>0|<f7>0|<f8>0}}"];
    carpetaRaiz:f2 -> carpetaBin:f0 [color="#14106C"];

    detalleBin [style=bold color="#370A19" label="{<f0>DETALLE|{archivo.txt|<f1>}|{|<f2>}|{|<f3>}|{|<f4>}|{|<f5>}|{<f6>0}}"];
    carpetaBin:f7 -> detalleBin:f0 [color="#14106C"];

    inodoArchivo [style=bold color="#0A370B" label="{<f0>INODOS|{Bloque 1|<f1>1}|{Bloque 2|<f2>0}|{Bloque 3|<f3>0}|{Bloque 4|<f4>0}|{<f5>0}}"];
    detalleBin:f1 -> inodoArchivo:f0 [color="#14106C"];
    dataArchivo [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoArchivo:f1 -> dataArchivo:f0 [color="#14106C"];
     dataArchivo2 [style=bold color="#085E6F" label="{<f0>Data}"];
    inodoArchivo:f2 -> dataArchivo2:f0 [color="#14106C"];
}
*/
