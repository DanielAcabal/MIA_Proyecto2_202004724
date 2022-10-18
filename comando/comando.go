package comando

import (
	"MIA-Proyecto2_202004724/Estructuras"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*======================MKDISK=======================*/
func Mkdisk(commandArray []string) string{
	//mkdisk -Size=3000 -unit=K -path=/home/user/Disco4.dk
	// mkdisk -size=5 -unit=M -path="/home/mis discos/Disco3.dk"
	consola :=""
	tamano := 0
	dimensional := "MB"
	ajuste := "FF"
	cantidad := 0
	tamano_archivo := 0
	limite := 0
	bloque := make([]byte, 1024)
	seguir := true
	ruta := ""

	// Lectura de parametros del comando
	for i := 1; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data, "size=") {
			strtam := strings.Replace(data, "size=", "", 1)
			strtam = strings.Replace(strtam, "\"", "", 2)
			strtam = strings.Replace(strtam, "\r", "", 1)
			tamano2, err := strconv.Atoi(strtam)
			tamano = tamano2
			if err != nil {
				msg_error(err)
			}
			if tamano<=0{
				consola +="El tamaño debe ser mayor o igual a cero\n"
				seguir = false
			}
			cantidad +=1
		} else if strings.Contains(data, "unit=") {
			dimensional = strings.Replace(data, "unit=", "", 1)
			dimensional = strings.Replace(dimensional, "\"", "", 2)
		}else if strings.Contains(data,"fit="){
			ajuste = strings.Replace(data,"fit=","",1)
			ajuste = strings.Replace(ajuste, "\"", "", 2)
			if ajuste!="bf" && ajuste!="ff" && ajuste!="wf"{
				seguir = false
				consola +="Ajuste no válido\n"
			}
		}else if strings.Contains(data,"path="){
			ruta = strings.Replace(data,"path=","",1)
			ruta = strings.Replace(ruta, "\"", "", 2)
			reg := regexp.MustCompile("/[a-zA-Z0-9]+.dk")
			junto := reg.ReplaceAllString(ruta,"")
			//seguir = false
			junto = "/home"+junto+""
			err := os.MkdirAll(junto,os.ModePerm)
			if err != nil {
				msg_error(err)
			}
			cantidad+=1
		}else{
			consola +="Parámetro no válido:"+data+"\n"
		}
	}
	if !seguir{
		return consola
	}
	if cantidad!=2{
		return consola
	}
	// Calculo de tamaño del archivo
	if strings.Contains(dimensional, "k") {
		tamano_archivo = tamano
	} else if strings.Contains(dimensional, "m") {
		tamano_archivo = tamano * 1024
	} else if strings.Contains(dimensional, "g") {
		tamano_archivo = tamano * 1024 * 1024
	}

	// Preparacion del bloque a escribir en archivo
	for j := 0; j < 1024; j++ {
		bloque[j] = 0
	}
	ruta = "/home"+ruta+""
	// Creacion, escritura y cierre de archivo
    disco, err := os.Create(ruta)
	if err != nil {
		msg_error(err)
	}
	for limite < tamano_archivo {
		_, err := disco.Write(bloque)
		if err != nil {
			msg_error(err)
		}
		limite++
	}
	mbr := estructuras.MBR{}
	copy(mbr.Mbr_tamano[:],strconv.Itoa(tamano))
	tiempo := time.Now()
	s1 := rand.NewSource(tiempo.UnixNano())
    r1 := rand.New(s1)
	copy(mbr.Mbr_fecha_creacion[:],tiempo.String())
	copy(mbr.Mbr_dsk_signature[:],strconv.Itoa(r1.Intn(100)))
	copy(mbr.Dsk_fit[:],ajuste)
	res := Struct_to_bytes(mbr)
	puntero,er := disco.Seek(0,io.SeekStart)
	if er != nil {
		msg_error(err)
	}
	_, err = disco.WriteAt(res, puntero)
		if err != nil {
			msg_error(err)
		}
	disco.Close()
	if err != nil {
		msg_error(err)
	}
	// Resumen de accion realizada
	consola +="Creacion de Disco:\n"
	consola +=" Tamaño: "
	consola += strconv.Itoa(tamano)
	consola +="\n Dimensional: "
	consola += dimensional
	return consola
}
func Rmdisk(commandArray []string) string{
	// rmdisk -path=/home/user2/Disco1.dk
	consola := ""
	ruta := ""
	seguir := true
	cant := 0
	for i := 1; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data,"path="){
			ruta = strings.Replace(data,"path=","",1)
			ruta = strings.Replace(ruta,"\"","",2)
			cant += 1
		}else{
			consola += "Parámetro no válido: "+data+"\n"
			seguir = false
		}
	}
	if !seguir{return consola}
	if cant != 1 {return consola}
	ruta = "/home"+ruta
	if _, err := os.Stat(ruta); err == nil {
		e := os.Remove(ruta)
		if e != nil {
		msg_error(e)
		return consola
		}
	 } else {
		consola += "File does not exist\n"
		return consola
	 }
	consola += "Disco eliminado\n"
	return consola
}
func msg_error(err error){
	fmt.Println("Error: ",err)
}
func Struct_to_bytes(p interface{}) []byte {
	// Codificacion de Struct a []Bytes
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return buf.Bytes()
}
func Fdisk(commandArray []string) string {
	consola :=""
	nuevo := estructuras.Particion{}
	copy(nuevo.Part_fit[:],[]byte("F"))
	copy(nuevo.Part_status[:],[]byte("U"))
	copy(nuevo.Part_type[:],[]byte("P"))
	nombreParticion :="" 
	tamanio := -1
    ruta :=""
	unidad := 'K'
	eliminar := false
	//nuevoEspacio := 0
    primero := ""
	cant := 0
	for i := 1; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		if strings.Contains(data,"path="){
			ruta = strings.Replace(data,"path=","",1)
			ruta = strings.Replace(ruta,"\"","",2)
			cant += 1
		}else if strings.Contains(data,"size="){
			tam := strings.Replace(data,"size=","",1)
			tamanio1, e := strconv.Atoi(tam)	
			if e != nil{
				msg_error(e)
			}
			if tamanio1 <= 0{
				consola += "Tamaño es menor que cero"
			}else{
				if(primero==""){primero="Size";tamanio = tamanio1}
			}
		}else if strings.Contains(data,"unit="){
			aux := strings.Replace(data,"unit=","",1)
			if (aux == "B" || aux == "b"){unidad = 'B';
			}else if (aux == "K" || aux == "k"){unidad = 'K';
			}else if (aux == "M" || aux == "m"){unidad = 'M';
			}else {consola += "Dimensión de partición incorrecta"
			cant--}
		}else if strings.Contains(data,"type="){
			tipo := strings.Replace(data,"type=","",1)
			if (tipo == "P" || tipo == "p"){copy(nuevo.Part_type[:],[]byte("P"))
			}else if (tipo == "E" || tipo == "e"){copy(nuevo.Part_type[:],[]byte("E"))
			}else if (tipo == "L" || tipo == "l"){copy(nuevo.Part_type[:],[]byte("L"))
			}else {consola += "Tipo de partición incorrecta";cant--;}
		}else if strings.Contains(data,"fit="){
			ajuste := strings.Replace(data,"fit=","",1)
			if (ajuste == "BestFit"||ajuste=="bf"||ajuste=="BF"){copy(nuevo.Part_fit[:],[]byte("B"))
			}else if (ajuste == "FirstFit"||ajuste=="ff"||ajuste=="FF"){copy(nuevo.Part_fit[:],[]byte("F"))
			}else if (ajuste == "WorstFit"||ajuste=="wf"||ajuste=="WF"){copy(nuevo.Part_fit[:],[]byte("W"))
			}else {consola += "Ajuste de partición incorrecta"; cant--;}
            
		}else if strings.Contains(data,"name="){
			name := strings.Replace(data,"name=","",1)
			nombreParticion = name
			cant++
		}else{
			consola += "Parámetro no válido"
		}
	}	
	if (cant==2){
        //Ejecutar fdisk
        //Leer disco
        ruta = strings.Replace(ruta,"\"","",2);
		ruta = "/home"+ruta
        nombreParticion = strings.Replace(nombreParticion,"\"","",2);
        copy(nuevo.Part_name[:],nombreParticion);
        if (eliminar){
            fmt.Print("Eliminando particion: ",nombreParticion)
           // eliminarParticion(path,nombreParticion);
        }else{
            if(primero=="Size"){
            consola += "Creando particion: "+nombreParticion
            crearParticion(nuevo,ruta,tamanio,unidad,0);
            }else{
            fmt.Print("Cambio de tamaño a particion: ",nombreParticion);
            //addParticion(path.c_str(),nuevoEspacio,unidad,nombreParticion);
            }
        }
    }else{
        consola += "Error de creacion/eliminacion de partición, faltan parámetros"
    }
	return consola
}
func crearParticion(particion estructuras.Particion,path string, size int, unit rune,add int){
    //fdisk -size=1000  -unit=B -path=/home/user/Disco4.dk -name=Particion4
    //fdisk -add=8 -s=10 -unit=K -path=/home/user/Disco3.dk -name=Particion3
    //fdisk -delete=full -name=Particion4 -path=/home/user/Disco3.dk
	disco, err := os.OpenFile(path,os.O_RDWR,0660);
	if err != nil{
		msg_error(err)
		return
	}
    mbrEmpty := estructuras.MBR{};
	si := Struct_to_bytes(mbrEmpty)
	data := make([]byte,len(si))
	_, err = disco.ReadAt(data,int64(0))
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	mbr := BytesToStructMBR(data)
    //Verificar tipo de particiones en mbr
    aux := [4]estructuras.Particion{};
    aux[0] = mbr.Mbr_partition_1;
    aux[1] = mbr.Mbr_partition_2;
    aux[2] = mbr.Mbr_partition_3;
    aux[3] = mbr.Mbr_partition_4;
    primarias := 0; extendidas := 0; logicas := 0;
    crear := true;
    for  i := 0; i < 4; i++ {
        if !bytes.Equal(aux[i].Part_start[:],[]byte{0}){ //Cero por si no está creada
        if bytes.Equal(aux[i].Part_type[:],[]byte{'P'}){primarias++;}
        if bytes.Equal(aux[i].Part_type[:],[]byte{'E'}){extendidas++;}
        if bytes.Equal(aux[i].Part_type[:], []byte{'L'}){logicas++;} //Creo que deberia ir aquí aux
        if bytes.Equal(aux[i].Part_name[:],particion.Part_name[:]){
			crear=false; 
			break;}
    }   
    }
    //Verificar el nombre de la particion, si existe
    if (!crear) {fmt.Print("La partición ",particion.Part_name," ya existe");
	}else{
    if ((primarias+extendidas)<4){ // Para otra primaria o una extendida
        if bytes.Equal(mbr.Dsk_fit[:],[]byte{'F'}){
			bs:=calcularTamanio(size,unit)
            copy(particion.Part_size[:],strconv.Itoa(bs));
            c := primerAjuste(&mbr,&particion); // 0 == inicio disco
            if (c!=0){
				b := make([]byte,2)
				binary.LittleEndian.PutUint16(b,c)
				copy(particion.Part_start[:],b);
			}else{crear=false;}
			b := make([]byte,2)
			binary.LittleEndian.PutUint16(b,0)
			if bytes.Equal(aux[0].Part_start[:],b){
                mbr.Mbr_partition_1 = particion;
            }else if bytes.Equal(aux[1].Part_start[:],b){
                mbr.Mbr_partition_2 = particion;
            }else if bytes.Equal(aux[2].Part_start[:],b){
                mbr.Mbr_partition_3 = particion;
            }else if bytes.Equal(aux[3].Part_start[:],b){
                mbr.Mbr_partition_4 = particion;
            }  
    }
	}else{
        fmt.Print("Máximo de particiones creadas");
        crear = false;
    }
    
    if (crear){
		data = Struct_to_bytes(mbr)
		/*puntero,er := disco.Seek(int64(0),io.SeekStart)
		if er != nil {
			msg_error(err)
		}*/
		_,err =disco.WriteAt(data,0)
		if err != nil{
			msg_error(err)
		}
    }
    disco.Close()
}
}
func calcularTamanio(s int,u rune) int{
    tamanio_real := 0;
    if (u == 'B' || u == 'b'){
        tamanio_real = s;
    }else if (u == 'K' || u == 'k'){
        tamanio_real = s * 1024;
    }else if (u == 'M' || u == 'm'){
        tamanio_real = s * 1024 *1024;
    }else{
        return 0;
    }
    return tamanio_real;
}
func primerAjuste(mbr *estructuras.MBR, particion *estructuras.Particion) uint16{
    comienza := uint16(unsafe.Sizeof(mbr));
    aux := [4]estructuras.Particion{};
    aux[0] = mbr.Mbr_partition_1;
    aux[1] = mbr.Mbr_partition_2;
    aux[2] = mbr.Mbr_partition_3;
    aux[3] = mbr.Mbr_partition_4;
    n := 4;
    i, j:=0,0;
    for i = 0; i < n - 1; i++{
        for j = 0; j < n - i - 1; j++{
            if (binary.BigEndian.Uint16(aux[j].Part_start[:]) > binary.BigEndian.Uint16(aux[j + 1].Part_start[:])){
                aux[j], aux[j + 1] = aux[j + 1] , aux[j] 
			}
		}
	}
    for i := 0; i < n; i++ {
        espacio := binary.BigEndian.Uint16(aux[i].Part_start[:]) - comienza; //int
        if (espacio>=binary.BigEndian.Uint16(particion.Part_size[:])){
			b := make([]byte, 2)
			binary.LittleEndian.PutUint16(b, uint16(i))
            copy(particion.Part_start[:],b);
            return comienza;
        }
        if(comienza<=binary.BigEndian.Uint16(aux[i].Part_start[:])){
            comienza = binary.BigEndian.Uint16(aux[i].Part_size[:]) + binary.BigEndian.Uint16(aux[i].Part_start[:]);
        }
    }
    if (comienza==0){comienza = uint16(unsafe.Sizeof(mbr));}
    x:= binary.BigEndian.Uint16(mbr.Mbr_tamano[:]); //int
    if ((x-comienza)>=binary.BigEndian.Uint16(particion.Part_size[:])){
        return comienza;
    }else{
        fmt.Print("No hay espacio para esta particion: ",particion.Part_name);
    }
    return 0;
}

func BytesToStructMBR(s []byte) estructuras.MBR {
	// Decodificacion de [] Bytes a Struct ejemplo
	p := estructuras.MBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}