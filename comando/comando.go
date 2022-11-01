package comando

import (
	"MIA-Proyecto2_202004724/Estructuras"
	"MIA-Proyecto2_202004724/helpers"
	"bytes"
	"container/list"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)
var particionesMontadas = list.New()
var particionActual = estructuras.Pmontada{}
var usuarioActual = estructuras.Usuario{}
var sesionInicida = false
var RepFile = "digraph{\"RepFile\"}"
var RepSB = "digraph{\"RepSB\"}"
/*======================MKDISK=======================*/
func Mkdisk(commandArray []string) string{
	//mkdisk -Size=3000 -unit=K -path=/home/user/Disco4.dk
	// mkdisk -size=5 -unit=M -path="/home/mis discos/Disco3.dk"
	consola :="==============MKDISK==============\n"
	tamano := int64(0)
	dimensional := "MB"
	ajuste := "FF"
	cantidad := 0
	tamano_archivo := int64(0)
	limite := int64(0)
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
			tamano = int64(tamano2)
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
			if ajuste!="bf" && ajuste!="ff" && ajuste!="wf" && ajuste!="bestfit" && ajuste!="firstfit" && ajuste!="worstfit"{
				seguir = false
				consola +="Ajuste no válido\n"
			}
		}else if strings.Contains(data,"path="){
			ruta = strings.Replace(data,"path=","",1)
			ruta = strings.Replace(ruta, "\"", "", 2)
			reg := regexp.MustCompile("/[a-zA-Z0-9 ]+.dk")
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
	copy(mbr.Mbr_tamano[:],helpers.IntToByteArray(tamano_archivo*1024))
	tiempo := time.Now()
	s1 := rand.NewSource(tiempo.UnixNano())
    r1 := rand.New(s1)
	copy(mbr.Mbr_fecha_creacion[:],tiempo.String())
	copy(mbr.Mbr_dsk_signature[:],helpers.IntToByteArray(int64(r1.Intn(100))))
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
	consola += strconv.Itoa(int(tamano))
	consola +="\n Dimensional: "
	consola += dimensional
	return consola
}
func Rmdisk(commandArray []string) string{
	// rmdisk -path=/home/user2/Disco1.dk
	consola := "==========RMDISK==========\n"
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
	consola :="==========FDISK==========\n"
	nuevo := estructuras.Particion{}
	copy(nuevo.Part_fit[:],[]byte("F"))
	copy(nuevo.Part_status[:],[]byte("U"))
	copy(nuevo.Part_type[:],[]byte("P"))
	nombreParticion :="" 
	tamanio := -1
    ruta :=""
	unidad := 'K'
	eliminar := false
	nuevoEspacio := 1
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
		}else if strings.Contains(data,"add="){
			tam := strings.Replace(data,"add=","",1)
			if strings.Contains(tam,"~"){
				nuevoEspacio = -1
				tam = strings.Replace(tam,"~","",1)
			}
			x,err := strconv.Atoi(tam)
			if err!=nil{
				msg_error(err)
			}
			nuevoEspacio *= x
            if (nuevoEspacio==0){consola+="No se puede añadir 0 espacio\n";cant--;
			}else {if (primero==""){primero="Add";}}
		}else if strings.Contains(data,"delete="){
			dato := strings.Replace(data,"delete=","",1)
			if (dato == "Full" || dato == "full"){eliminar = true;
			}else {consola += "Opción de eliminación de partición incorrecta\n";cant--;}
		} else{
			consola += "Parámetro no válido\n"
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
            consola += "Eliminando particion: "+nombreParticion+"\n"
            eliminarParticion(ruta,nombreParticion,&consola);
        }else{
            if(primero=="Size"){
            consola += "Creando particion: "+nombreParticion
            crearParticion(nuevo,ruta,tamanio,unidad,0,&consola);
            }else{
            consola += "Cambio de tamaño a particion: "+nombreParticion+"\n"
            addParticion(ruta,nuevoEspacio,unidad,nombreParticion,&consola);
            }
        }
    }else{
        consola += "Error de creacion/eliminacion de partición, faltan parámetros\n"
    }
	return consola
}
func crearParticion(particion estructuras.Particion,path string, size int, unit rune,add int,consola *string){
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
    if (!crear) {*consola+="La partición "+string(particion.Part_name[:])+" ya existe";
	}else{
    if ((primarias+extendidas)<4){ // Para otra primaria o una extendida
        c :=int64(0)
		if bytes.Equal(mbr.Dsk_fit[:],[]byte{'F'}){
			bs:=int64(calcularTamanio(size,unit))
            copy(particion.Part_size[:],helpers.IntToByteArray(bs));
            c = primerAjuste(&mbr,&particion,consola); // 0 == inicio disco
		}else if bytes.Equal(mbr.Dsk_fit[:],[]byte{'B'}){
			bs:=int64(calcularTamanio(size,unit))
            copy(particion.Part_size[:],helpers.IntToByteArray(bs));
            c = mejorAjuste(&mbr,&particion,consola); // 0 == inicio disco
		}else if bytes.Equal(mbr.Dsk_fit[:],[]byte{'W'}){
			bs:=int64(calcularTamanio(size,unit))
            copy(particion.Part_size[:],helpers.IntToByteArray(bs));
            c = PeorAjuste(&mbr,&particion,consola); // 0 == inicio disco
		}
            if (c!=0){
				copy(particion.Part_start[:],helpers.IntToByteArray(c));
			}else{crear=false;}

			if helpers.ByteArrayToInt64(aux[0].Part_start[:])==0{
                mbr.Mbr_partition_1 = particion;
            }else if helpers.ByteArrayToInt64(aux[1].Part_start[:])==0{
                mbr.Mbr_partition_2 = particion;
            }else if helpers.ByteArrayToInt64(aux[2].Part_start[:])==0{
                mbr.Mbr_partition_3 = particion;
            }else if helpers.ByteArrayToInt64(aux[3].Part_start[:])==0{
                mbr.Mbr_partition_4 = particion;
            }  
    
	}else{
        *consola += "Máximo de particiones creadas";
        crear = false;
    }
    
    if (crear){
		data = Struct_to_bytes(mbr)
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
func primerAjuste(mbr *estructuras.MBR, particion *estructuras.Particion,consola *string) int64{
    var info estructuras.MBR
	comienza := helpers.HandleSizeof(info);
    aux := [4]estructuras.Particion{};
    aux[0] = mbr.Mbr_partition_1;
    aux[1] = mbr.Mbr_partition_2;
    aux[2] = mbr.Mbr_partition_3;
    aux[3] = mbr.Mbr_partition_4;
    n := 4;
    i, j:=0,0;
    for i = 0; i < n - 1; i++{
        for j = 0; j < n - i - 1; j++{
			actual := helpers.ByteArrayToInt64(aux[j].Part_start[:])
			siguiente := helpers.ByteArrayToInt64(aux[j+1].Part_start[:])
            if (actual > siguiente){
                aux[j], aux[j + 1] = aux[j + 1] , aux[j] 
			}
		}
	}
    for i = 0; i < n; i++ {
		actual := helpers.ByteArrayToInt64(aux[i].Part_start[:])
        espacio := actual - comienza; //int
		necesito := helpers.ByteArrayToInt64(particion.Part_size[:])
        if (espacio>=necesito){
			copy(particion.Part_start[:],helpers.IntToByteArray(comienza))
            return comienza;
        }

        if(espacio>=0){
			nuevo := helpers.ByteArrayToInt64(aux[i].Part_size[:])
			actual := helpers.ByteArrayToInt64(aux[i].Part_start[:])
            comienza = actual+ nuevo;
        }
    }
    if (comienza==0){comienza = helpers.HandleSizeof(info);}
    x := helpers.ByteArrayToInt64(mbr.Mbr_tamano[:]); //int
	nuevo := helpers.ByteArrayToInt64(particion.Part_size[:])
    if ((x-comienza)>=nuevo){
        return comienza;
    }else{
        *consola += "No hay espacio para esta particion: "+string(particion.Part_name[:]);
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

func mejorAjuste(mbr *estructuras.MBR,particion *estructuras.Particion,consola *string) int64{
    
    aux := [4]estructuras.Particion{};
    aux[0] = mbr.Mbr_partition_1;
    aux[1] = mbr.Mbr_partition_2;
    aux[2] = mbr.Mbr_partition_3;
    aux[3] = mbr.Mbr_partition_4;
    pares := [5]estructuras.Pares{};
    n:= 4;
    i, j :=0,0
    for i = 0; i < n - 1; i++{
        for j = 0; j < n - i - 1; j++{
            if (helpers.ByteArrayToInt64(aux[j].Part_start[:]) > helpers.ByteArrayToInt64(aux[j + 1].Part_start[:])){
                aux[j], aux[j + 1] = aux[j + 1] ,aux[j] 
			}
		}
	}
	var info estructuras.MBR
	comienza := helpers.HandleSizeof(info);
    for i = 0; i < n; i++{
        espacio := helpers.ByteArrayToInt64(aux[i].Part_start[:]) - comienza;
        pares[i].Inicio = comienza;
        pares[i].Tamanio = espacio;
        if (espacio>=0){
        comienza = helpers.ByteArrayToInt64(aux[i].Part_start[:])+helpers.ByteArrayToInt64(aux[i].Part_size[:]);
        }
    }
    if (comienza==0){comienza = helpers.HandleSizeof(info);}
    espacio := helpers.ByteArrayToInt64(mbr.Mbr_tamano[:]) - comienza;
    pares[4].Inicio = comienza;
    pares[4].Tamanio = espacio;
    for i = 0; i < 5 - 1; i++{
        for j = 0; j < 5 - i - 1; j++{
            if (pares[j].Tamanio > pares[j + 1].Tamanio){
                pares[j], pares[j + 1] = pares[j+1], pares[j] 
			}
		}
	}
    for i = 0; i < 5; i++{
        x :=pares[i].Tamanio;
    	a := helpers.ByteArrayToInt64(particion.Part_size[:]);
        if(x>a){
            return pares[i].Inicio;
        }
    }
        *consola += "No hay espacio para esta particion: "+string(particion.Part_name[:]);
    return 0;
}

func PeorAjuste(mbr *estructuras.MBR, particion *estructuras.Particion,consola *string)int64{
    
	aux := [4]estructuras.Particion{};
    aux[0] = mbr.Mbr_partition_1;
    aux[1] = mbr.Mbr_partition_2;
    aux[2] = mbr.Mbr_partition_3;
    aux[3] = mbr.Mbr_partition_4;

    n := 4;
    i, j := 0,0;
    for i = 0; i < n - 1; i++{
        for j = 0; j < n - i - 1; j++{
            if helpers.ByteArrayToInt64(aux[j].Part_start[:]) > helpers.ByteArrayToInt64(aux[j + 1].Part_start[:]){
                aux[j], aux[j + 1] = aux[j+1], aux[j]
			}
		}
	}
	var info estructuras.MBR
    tamanio := int64(0); comienza := helpers.HandleSizeof(info); start := comienza;
    for i = 0; i < n; i++{
        espacio := helpers.ByteArrayToInt64(aux[i].Part_start[:]) - comienza;
        if (espacio>tamanio){
            tamanio = espacio;
            start = comienza;

        }
        if(espacio>=0){
        comienza = helpers.ByteArrayToInt64(aux[i].Part_start[:])+helpers.ByteArrayToInt64(aux[i].Part_size[:]);
        }
    }
    if (comienza==0){comienza = helpers.HandleSizeof(info);}

    x := helpers.ByteArrayToInt64(mbr.Mbr_tamano[:]); 
    espacio := x - comienza;
    if (espacio>tamanio){
        tamanio = espacio;
        start = comienza;
    }
    if (tamanio>=helpers.ByteArrayToInt64(particion.Part_size[:])){
        return start;
    }else{
        *consola += "No hay espacio para esta particion: "+string(particion.Part_name[:]);
    }
    return 0;    
}
func addParticion(path string, size int,unit rune,nombre string,consola *string){
    tamanioBytes := int64(calcularTamanio(size,unit));//Ha sumar o restar
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
    aux := [4]*estructuras.Particion{};
    aux[0] = &mbr.Mbr_partition_1;
    aux[1] = &mbr.Mbr_partition_2;
    aux[2] = &mbr.Mbr_partition_3;
    aux[3] = &mbr.Mbr_partition_4;
    
    n := 4;
    i, j := 0,0;
    for i = 0; i < n - 1; i++{ //Se ordena
        for j = 0; j < n - i - 1; j++{
            if (helpers.ByteArrayToInt64(aux[j].Part_start[:]) > helpers.ByteArrayToInt64(aux[j + 1].Part_start[:])){
                aux[j], aux[j + 1] = aux[j+1], aux[j]
			}
		}

	}
    encontrado := false;
    guardar := false;
    for i = 0; i < n; i++{
        p:= aux[i].Part_name[:];
		aux := make([]byte,15)
		copy(aux,nombre)
        if bytes.Equal(p[:],aux[:]){
            encontrado =true;
            break;
        }
    }
    if (encontrado){
        res := helpers.ByteArrayToInt64(aux[i].Part_size[:]) + tamanioBytes; //Positivo o negativo
        if(res<=0){
            *consola += "No quedará espacio en la partición\n";
        }else{
            fin := helpers.ByteArrayToInt64(mbr.Mbr_tamano[:])
            if(i+1<4){
                fin =helpers.ByteArrayToInt64(aux[i+1].Part_start[:])
            }
            if ((fin-helpers.ByteArrayToInt64(aux[i].Part_start[:]))>=res){//Si lo puede guardar
                vacio := make([]byte,10)
				copy(aux[i].Part_size[:],vacio);
				copy(aux[i].Part_size[:],helpers.IntToByteArray(int64(res)));
                guardar = true;
            }else{
                *consola += "Espacio insuficiente\n"
            }
        }
    }else{
        *consola +="La partición no existe\n"
    }
    if (guardar){
	data := Struct_to_bytes(mbr)
	pos,_:=disco.Seek(0,io.SeekStart)
	_,err :=disco.WriteAt(data,pos)
	if err!=nil{
		msg_error(err)
	}
    }
    disco.Close()
}

func eliminarParticion(path string, name string,consola *string){
    //Verificar que la partición exista
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
    encontrado := false;
	i:=0
    for  i = 0; i < 4; i++{
        if helpers.ByteArrayToInt64(aux[i].Part_start[:])!=0{
            p:= aux[i].Part_name[:];
			aux := make([]byte,50)
			copy(aux,name)
        	if bytes.Equal(p[:],aux[:]){
            encontrado =true;
            break;
        }
        }
    }
    if (encontrado){ //FULL
        nuevo :=estructuras.Particion{};
        ceros := make([]byte,1)
        ceros[0] = 0;
        j:=int64(0);
		puntero, err :=disco.Seek(helpers.ByteArrayToInt64(aux[i].Part_start[:]),io.SeekStart)
        if err!=nil{
			msg_error(err)
		}
		for (j!=helpers.ByteArrayToInt64(aux[i].Part_size[:])){
			disco.WriteAt(ceros,puntero)
            j++;
        }
        if (i==0){mbr.Mbr_partition_1 = nuevo
		}else if (i==1){mbr.Mbr_partition_2 = nuevo
		}else if (i==2){mbr.Mbr_partition_3 = nuevo
		}else if (i==3){mbr.Mbr_partition_4 = nuevo}
		puntero,err = disco.Seek(0,io.SeekStart)
		if err !=nil{
			msg_error(err)
		}
		data := Struct_to_bytes(mbr)
		disco.WriteAt(data,puntero)
        //mostrarMBR(disco);       
    }
	disco.Close()
}
func Mount(parametros []string) string{
	consola := "==========MOUNT==========\n"
	ruta := "/home"
	name := ""
	cant := 0
	for i := 1; i < len(parametros); i++ {
		parametro :=  strings.ToLower(parametros[i])
		if strings.Contains(parametro,"path="){
			ruta += strings.Replace(parametro,"path=","",1)
			cant++
		}else if strings.Contains(parametro,"name="){	
			name += strings.Replace(parametro,"name=","",1)
			cant++
		}else{
			consola += "Parámetro "+parametro+" no válido\n"
			return consola
		}
	}
	if cant==2{
		consola += montarPart(ruta,name)
	}else{
		consola += "Parámetros insuficientes\n"
		return consola
	}

	return consola
}
func montarPart(ruta string,name string)string{
	consola :=""
	disco, err := os.OpenFile(ruta,os.O_RDWR,0660);
	if err != nil{
		msg_error(err)
		return ""
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
    aux := [4]*estructuras.Particion{};
    aux[0] = &mbr.Mbr_partition_1;
    aux[1] = &mbr.Mbr_partition_2;
    aux[2] = &mbr.Mbr_partition_3;
    aux[3] = &mbr.Mbr_partition_4;
    encontrado := false;
	i:=0
    for  i = 0; i < 4; i++{
        if helpers.ByteArrayToInt64(aux[i].Part_start[:])!=0{
			p:= aux[i].Part_name[:];
			aux := make([]byte,15)
			copy(aux,name)
			if bytes.Equal(p[:],aux[:]){
            encontrado =true;
            break;
        }
        }
    }	
	if encontrado{
		letra := rune(i+97)
		id:= "24"+strconv.Itoa(i)
		id += string(letra)
		partMontada := estructuras.Pmontada{}
		aux[i].Part_status = [1]byte{'M'}
		partMontada.Id = id
		partMontada.Particion = *aux[i]
		partMontada.Path = ruta
		tiempo := time.Now()
		partMontada.TiempoM = tiempo.String()
		particionesMontadas.PushFront(partMontada)
		consola += "Partición "+name+" montada, id: "+id
	}else{
		consola += "Partición "+name+" no encontrada\n"		
	}
	disco.Close()
	return consola
}
func Unmount(parametros []string)string{
	idBuscar :=""
	seguir := false
	for i := 1; i < len(parametros); i++ {
		parametro := parametros[i]
		if strings.Contains(parametro,"id="){
			idBuscar = strings.Replace(parametro,"id=","",1)
			seguir = true
		}else{
			return "Parámetro "+parametro+" no válido\n"
		}
	}
	if seguir{
		for element := particionesMontadas.Front(); element != nil; element = element.Next() {
			// do something with element.Value
			part := estructuras.Pmontada(element.Value.(estructuras.Pmontada))
			if part.Id == idBuscar{
				particionesMontadas.Remove(element)
				return "Partición "+idBuscar+" desmontada\n"				
			}
			}
	}else{
		return "Faltan parámetros\n"
	}
	return ""
}
func Mkfs(parametros []string) string{
	id := ""
	tipo := "full"
	cant := 0
	for i := 1; i < len(parametros); i++ {
		parametro := strings.ToLower(parametros[i])
		if strings.Contains(parametro,"id="){
			id = strings.Replace(parametro,"id=","",1)
			cant++
		}else if strings.Contains(parametro,"type="){
			tipo = strings.Replace(parametro,"type=","",1)
			if tipo != "full"{
				cant = -1
				return "Tipo de formateo no válido\n"
			}
		}else{
			return "Parámetro "+parametro+" no válido\n"
		}
	}
	if cant==1{
		return formatear(id,tipo)
	}else{
		return "Faltan parámetros obligatorios (id)\n"
	}
}
func obtenerMBR(ruta string, mbr *estructuras.MBR,aux *[4]estructuras.Particion) *os.File{
	disco, err := os.OpenFile(ruta,os.O_RDWR,0660);
	if err != nil{
		msg_error(err)
	}
    mbrEmpty := estructuras.MBR{};
	si := Struct_to_bytes(mbrEmpty)
	data := make([]byte,len(si))
	_, err = disco.ReadAt(data,int64(0))
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	*mbr = BytesToStructMBR(data)
	aux[0] = mbr.Mbr_partition_1
	aux[1] = mbr.Mbr_partition_2
	aux[2] = mbr.Mbr_partition_3
	aux[3] = mbr.Mbr_partition_4
    return disco 
}
func getParticionMontada(id string,parti *estructuras.Pmontada) string{
	for element := particionesMontadas.Front(); element != nil; element = element.Next() {
		// do something with element.Value
		part := estructuras.Pmontada(element.Value.(estructuras.Pmontada))
		if part.Id == id{
			//fmt.Print(parti)
			*parti = part
			return ""				
		}
	}
	return "Partición "+id+" no encontrada, asegurese de que esté montada"
}
func formatear(id string,tipo string) string{
	particion := estructuras.Pmontada{}
	consola := "=====Formateando partición "+id+"=====\n"
	consola += getParticionMontada(id,&particion)
	if strings.Contains(consola,"no encontrada"){return consola}
	mbr := estructuras.MBR{}
	aux := [4]estructuras.Particion{}
	disco := obtenerMBR(particion.Path,&mbr,&aux)	
	//Limpiar partición en el MBR
	encontrado := false
	i :=0
	for i = 0; i < len(aux); i++ {
        	if bytes.Equal(aux[i].Part_name[:],particion.Particion.Part_name[:]){
            encontrado =true;
            break;
        }
	}
	if encontrado{
		//Limpiado
		ceros := make([]byte,1)
		puntero ,err := disco.Seek(int64(helpers.ByteArrayToInt64(aux[i].Part_start[:])),io.SeekStart)
		if err!=nil{
			msg_error(err)
		}
		j := int64(0)
		fin := helpers.ByteArrayToInt64(aux[i].Part_size[:])
		for j!= fin{
			disco.WriteAt(ceros,puntero)
			j++
		}
		//Cantidad de archivos
		inicio := helpers.ByteArrayToInt64(aux[i].Part_start[:])
		tamanioParticion := helpers.ByteArrayToInt64(aux[i].Part_size[:])
		tamanioSuper := helpers.HandleSizeof(estructuras.SuperBloque{})
		tamanioInodo := helpers.HandleSizeof(estructuras.Inodo{})
		tamanioBloque := helpers.HandleSizeof(estructuras.BloqueCarpeta{})
		n := (tamanioParticion-tamanioSuper)/(4+3*tamanioBloque+tamanioInodo);
    	numero := n; //Cantidad de inodos
		fmt.Println(numero)
		//llenando SuperBloque
		superbloque := estructuras.SuperBloque{}
		copy(superbloque.S_filesystem_type[:],helpers.IntToByteArray(2))
		copy(superbloque.S_inodes_count[:],helpers.IntToByteArray(numero))
		copy(superbloque.S_blocks_count[:],helpers.IntToByteArray(3*numero))
		copy(superbloque.S_free_inodes_count[:],helpers.IntToByteArray(numero))
		copy(superbloque.S_free_blocks_count[:],helpers.IntToByteArray(3*numero))
		tim := time.Now()
		copy(superbloque.S_mtime[:],[]byte(tim.String()))
		copy(superbloque.S_mnt_count[:],helpers.IntToByteArray(0))
		copy(superbloque.S_magic[:],"EF53")
		copy(superbloque.S_inode_size[:],helpers.IntToByteArray(tamanioInodo))
		copy(superbloque.S_block_size[:],helpers.IntToByteArray(tamanioBloque))
		copy(superbloque.S_firts_ino[:],helpers.IntToByteArray(0))
		copy(superbloque.S_first_blo[:],helpers.IntToByteArray(0))
		copy(superbloque.S_bm_inode_start[:],helpers.IntToByteArray(inicio+tamanioSuper))
		copy(superbloque.S_bm_block_start[:],helpers.IntToByteArray(inicio+tamanioSuper+numero))
		copy(superbloque.S_inode_start[:],helpers.IntToByteArray(inicio+tamanioSuper+4*numero))
		copy(superbloque.S_block_start[:],helpers.IntToByteArray(inicio+tamanioSuper+4*numero+numero*tamanioInodo))
		puntero, e := disco.Seek(inicio,io.SeekStart)
		if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(superbloque),puntero)
		disco.Close()
		iniciarBitmaps(particion.Path,helpers.ByteArrayToInt64(superbloque.S_bm_inode_start[:]),
						helpers.ByteArrayToInt64(superbloque.S_inodes_count[:]),
						helpers.ByteArrayToInt64(superbloque.S_bm_block_start[:]),
						helpers.ByteArrayToInt64(superbloque.S_blocks_count[:]))
	
		carpetaRaiz(&superbloque,aux[i],particion.Path)
		consola += archivoUsers("1,G,root\n1,U,root,root,123\n",particion.Path,&superbloque,aux[i])
	}else{
		consola += "Partición no encontrada en disco"
	}
	return consola
}
func carpetaRaiz(super *estructuras.SuperBloque, particion estructuras.Particion,ruta string){
	disco, err := os.OpenFile(ruta,os.O_RDWR,0660);
	if err != nil{
		msg_error(err)
	}
	//Posiciones para saber dónde escribir
	IniciaBitmapInodo := super.S_bm_inode_start[:]
	InodoLibreBM := helpers.ByteArrayToInt64(super.S_firts_ino[:])
    Inodo_libreI := helpers.ByteArrayToInt64(super.S_inode_start[:])+(InodoLibreBM)*helpers.ByteArrayToInt64(super.S_inode_size[:]);//posicion archivo del inodo libre
	

	inodoRaiz := estructuras.Inodo{}
	estructuras.NuevoInodo(&inodoRaiz,1,1,0,"0","777")
	//Actualizamos Bitmap, inodos libres, primer inodo libre
	actualizarBitmapInodo(disco,helpers.ByteArrayToInt64(super.S_inodes_count[:]),helpers.ByteArrayToInt64(IniciaBitmapInodo),super)
	nuevoLibre := helpers.ByteArrayToInt64(super.S_free_inodes_count[:])-1
	copy(super.S_free_inodes_count[:],helpers.IntToByteArray(nuevoLibre))

	// Bloque carpetaRaiz
	IniciaBitmapBloque := super.S_bm_block_start[:]
	PrimerBloqueLibre := helpers.ByteArrayToInt64(super.S_first_blo[:])
	BloqueLibre := helpers.ByteArrayToInt64(super.S_block_start[:]) + helpers.ByteArrayToInt64(super.S_block_size[:])*(PrimerBloqueLibre)
	
	carpetaR := estructuras.BloqueCarpeta{}
	estructuras.NuevoBloqueCarpeta(&carpetaR)//Constructor xd
	// Los primeros 2 registros del primer apuntador directo del Inodo son la carpeta y carpeta padre
	copy(carpetaR.B_content[0].B_name[:],"/")
	copy(carpetaR.B_content[0].B_inodo[:],helpers.IntToByteArray(InodoLibreBM)) // Apuntan al inodo creado antes
	copy(carpetaR.B_content[1].B_name[:],"/")
	copy(carpetaR.B_content[1].B_inodo[:],helpers.IntToByteArray(InodoLibreBM))
	//Primer apuntador directo 
	inodoRaiz.I_block[0] = helpers.IntToByteArray(PrimerBloqueLibre)[0]

	actualizarBitmapBloque(disco,helpers.ByteArrayToInt64(super.S_blocks_count[:]),
								helpers.ByteArrayToInt64(IniciaBitmapBloque[:]),super)

	nuevoBloqueLibre := helpers.ByteArrayToInt64(super.S_free_blocks_count[:])-1
	copy(super.S_free_blocks_count[:],helpers.IntToByteArray(nuevoBloqueLibre))	
	
	//Guardamos estructuras
	puntero,e:=disco.Seek(helpers.ByteArrayToInt64(particion.Part_start[:]),io.SeekStart)
	if e!=nil{msg_error(e)}
	disco.WriteAt(Struct_to_bytes(super),puntero)


	puntero,e =disco.Seek(Inodo_libreI,io.SeekStart)
	if e!=nil{msg_error(e)}
	disco.WriteAt(Struct_to_bytes(inodoRaiz),puntero)
	
	puntero,e =disco.Seek(BloqueLibre,io.SeekStart)
	if e!=nil{msg_error(e)}
	disco.WriteAt(Struct_to_bytes(carpetaR),puntero)
	disco.Close()
}
func actualizarBitmapInodo(disco *os.File,fin int64,inicio int64,super *estructuras.SuperBloque){
	//Agregamos en el último libre
	prt,e:=disco.Seek(helpers.ByteArrayToInt64(super.S_firts_ino[:])+inicio,io.SeekStart)
	if e!=nil{msg_error(e)}
	disco.WriteAt([]byte{'1'},prt)
	//Actualizamos el libre
	for i := int64(0); i < fin; i++ {
		ptr,e:=disco.Seek(inicio+i,io.SeekStart)
		if e!=nil{msg_error(e)}
		x := make([]byte,1)
		disco.ReadAt(x,ptr)
		if x[0] == '0'{
			copy(super.S_firts_ino[:],helpers.IntToByteArray(i))
			break
		}
	}
}
func actualizarBitmapBloque(disco *os.File,fin int64,inicio int64,super *estructuras.SuperBloque){
	//Agregamos en el último libre
	prt,e:=disco.Seek(helpers.ByteArrayToInt64(super.S_first_blo[:])+inicio,io.SeekStart)
	if e!=nil{msg_error(e)}
	disco.WriteAt([]byte{'1'},prt)
	//Actualizamos el libre
	for i := int64(0); i < fin; i++ {
		ptr,e:=disco.Seek(inicio+i,io.SeekStart)
		if e!=nil{msg_error(e)}
		x := make([]byte,1)
		disco.ReadAt(x,ptr)
		if x[0] == '0'{
			copy(super.S_first_blo[:],helpers.IntToByteArray(i))
			break
		}
	}
}
func archivoUsers(contenido string,ruta string,super *estructuras.SuperBloque,particion estructuras.Particion)string{
	disco, err := os.OpenFile(ruta,os.O_RDWR,0660);
	if err != nil{
		msg_error(err)
	}
	//Reservamos posiciones
	InicioBitmapInodo := helpers.ByteArrayToInt64(super.S_bm_inode_start[:])
	InodoLibreBM  := helpers.ByteArrayToInt64(super.S_firts_ino[:])
	InodoLibre := helpers.ByteArrayToInt64(super.S_inode_start[:])+helpers.ByteArrayToInt64(super.S_inode_size[:])*InodoLibreBM
	//Crear el Inodo
	InodoUsers := estructuras.Inodo{}
	estructuras.NuevoInodo(&InodoUsers,1,1,int64(len(contenido)),"1","777")
	//Actualizar bitmap, primer libre y cantidad libre
	actualizarBitmapInodo(disco,helpers.ByteArrayToInt64(super.S_inodes_count[:]),InicioBitmapInodo,super)
	nuevoCantInodoLibre := helpers.ByteArrayToInt64(super.S_free_inodes_count[:]) - 1
	copy(super.S_free_inodes_count[:],helpers.IntToByteArray(nuevoCantInodoLibre))
    
	//Reservamos posiciones x2
	InicioBitmapBloque := helpers.ByteArrayToInt64(super.S_bm_block_start[:])
	BloqueLibreBM  := helpers.ByteArrayToInt64(super.S_first_blo[:])
	BloqueLibre := helpers.ByteArrayToInt64(super.S_block_start[:])+helpers.ByteArrayToInt64(super.S_block_size[:])*(BloqueLibreBM)
	
	//Crear el archivo
	archivoNuevo := estructuras.BloqueArchivos{}
	copy(archivoNuevo.B_content[:],contenido)
	//Actualizar bitmap, primer libre y cantidad libre
	actualizarBitmapBloque(disco,helpers.ByteArrayToInt64(super.S_blocks_count[:]),InicioBitmapBloque,super)
	nuevoCantBloqueLibre := helpers.ByteArrayToInt64(super.S_free_blocks_count[:]) - 1
	copy(super.S_free_blocks_count[:],helpers.IntToByteArray(nuevoCantBloqueLibre))
    
	// Apuntar bloques
		//Bloque carpeta Raiz (Inicio de tabla bloques) -> InodoUsers
		data := make([]byte,helpers.HandleSizeof(estructuras.BloqueCarpeta{}))
		puntero,e:=disco.Seek(helpers.ByteArrayToInt64(super.S_block_start[:]),io.SeekStart); if e!=nil{msg_error(e)}
		disco.ReadAt(data,puntero)
		RaizCarpeta := helpers.ByteArrayToDirBlock(data)
		copy(RaizCarpeta.B_content[2].B_name[:],"users.txt")
		copy(RaizCarpeta.B_content[2].B_inodo[:],helpers.IntToByteArray(InodoLibreBM))
		puntero,e =disco.Seek(helpers.ByteArrayToInt64(super.S_block_start[:]),io.SeekStart); if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(RaizCarpeta),puntero)

		//InodoUsers -> archivoNuevo
		InodoUsers.I_block[0] = helpers.IntToByteArray(BloqueLibreBM)[0] //Revisar porque pueden ser +255
	//Escribir en archivo
	puntero,e = disco.Seek(helpers.ByteArrayToInt64(particion.Part_start[:]),io.SeekStart); if e!=nil{msg_error(e)}
	disco.WriteAt(Struct_to_bytes(super),puntero)
	puntero,e = disco.Seek(InodoLibre,io.SeekStart); if e!=nil{msg_error(e)}
	disco.WriteAt(Struct_to_bytes(InodoUsers),puntero)
	puntero,e = disco.Seek(BloqueLibre,io.SeekStart); if e!=nil{msg_error(e)}
	disco.WriteAt(Struct_to_bytes(archivoNuevo),puntero)
	disco.Close()
	return "Archivo Users.txt Creado"
}
func iniciarBitmaps(ruta string,inicioInodo int64,finInodo int64,inicioBloque int64, finBloque int64){
	disco, err := os.OpenFile(ruta,os.O_RDWR,0660);
	if err != nil{
		msg_error(err)
	}
	puntero, e := disco.Seek(inicioInodo,io.SeekStart)
	if e!=nil{
		msg_error(e)
	}
	for i := 0; i < int(finInodo); i++ {
		disco.WriteAt([]byte{'0'},puntero)
		inicioInodo ++
		puntero, e = disco.Seek(inicioInodo,io.SeekStart); if e!=nil{msg_error(e)}
	}
	puntero, e = disco.Seek(inicioBloque,io.SeekStart)
	if e!=nil{
		msg_error(e)
	}
	for i := 0; i < int(finBloque); i++ {
		disco.WriteAt([]byte{'0'},puntero)
		inicioBloque ++
		puntero, e = disco.Seek(inicioBloque,io.SeekStart); if e!=nil{msg_error(e)}
	}
}
func Login(parametros []string) string{
	cant := 0
	user := ""
	password :=""
	id :=""
	consola := "==========LOGIN==========\n"
	for i := 1; i < len(parametros); i++ {
		param := strings.ToLower(parametros[i])
		if strings.Contains(param,"usuario="){
			user = strings.Replace(param,"usuario=","",1)
			cant++
			if len(user)>10{cant--; consola+="Usuario debe ser menor o igual que 10 letras\n"}
		}else if strings.Contains(param,"password="){
			password = strings.Replace(param,"password=","",1)
			cant++
			if len(password)>10{cant--; consola+="Password debe ser menor o igual que 10 letras\n"}
		}else if strings.Contains(param,"id="){
			id = strings.Replace(param,"id=","",1)
			cant++
		}else{
			consola += "Error: Parámetro "+param+" no válido \n"
			cant = -1
		}
	}
	if cant==3{
		consola += iniciarSesion(id,user,password)
	}else{
		consola += "Error: Faltan parámetros obligatorios\n"
	}
	return consola
}
func iniciarSesion(id string,user string, password string) string{
	consola :=""
	if sesionInicida{ return "Debe cerrar la sesión actual, use el comando LOGOUT"}
	consola += getParticionMontada(id,&particionActual)
	if strings.Contains(consola,"no encontrada"){return consola}
	disco,e := os.OpenFile(particionActual.Path,os.O_RDWR,0660); if e!=nil{msg_error(e)}
	superBloque := helpers.ReadSuperBlock(disco,helpers.ByteArrayToInt64(particionActual.Particion.Part_start[:]))
	carpetas := strings.Split("/users.txt","/")
	pos := buscarArchivo(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),&superBloque,disco,carpetas,1)
	if pos == -1{consola +="Error: No se encontró el archivo/Directorio"+"/users.txt\n"}
	inodo := helpers.ReadInode(disco,int64(pos))
	contenido := contenidoArchivo(inodo,&superBloque,disco)
	encontrado,index,init := buscarUsuario(contenido,user,password,true)
	index++;init++; // No se usan xd, para que no del error ;v
	if encontrado{
		sesionInicida = true
		usuarioActual.Id = id
		usuarioActual.User = user
		usuarioActual.Password = password
		consola += "Ingreso exitoso"
	}else{consola += "Error: El usuario"+ user+" no existe"}
	disco.Close()
	return consola
}
func buscarArchivo(posInodo int64,super *estructuras.SuperBloque,disco *os.File, directorio []string,indice int) int{
	inodo := helpers.ReadInode(disco,posInodo)
	if inodo.I_type[0] == '1' {return int(posInodo)}
	for i := 0; i < len(inodo.I_block); i++ {
		directo := inodo.I_block[i]
		if (directo != helpers.IntToByteArray(-1)[0]){
			x,e:= strconv.ParseInt(string(directo),36,64); if e!=nil{msg_error(e)}
			inicio := helpers.ByteArrayToInt64(super.S_block_start[:])
			size := helpers.ByteArrayToInt64(super.S_block_size[:])
			puntero,e:=disco.Seek(inicio+x*size,io.SeekStart); if e!=nil{msg_error(e)}
			data := make([]byte,helpers.HandleSizeof(estructuras.BloqueCarpeta{}))
			disco.ReadAt(data,puntero)
			carpeta := helpers.ByteArrayToDirBlock(data)
			for j := 0; j < len(carpeta.B_content); j++ {
				dir := carpeta.B_content[j]
				nombre :=make([]byte,12)
				copy(nombre,directorio[indice])
				if bytes.Equal(dir.B_name[:],nombre){
					bit := helpers.ByteArrayToInt64(dir.B_inodo[:])
					inicio = helpers.ByteArrayToInt64(super.S_inode_start[:])
					size = helpers.ByteArrayToInt64(super.S_inode_size[:])
					return buscarArchivo(bit*size+inicio,super,disco,directorio,indice+1)
				}
			}
		}
	}
	return -1
}
func contenidoArchivo(inodo estructuras.Inodo,super *estructuras.SuperBloque, disco *os.File) string{
	contenido :=""
	for i := 0; i < len(inodo.I_block); i++ {
		directo := inodo.I_block[i]
		if directo != helpers.IntToByteArray(-1)[0]{
			inicio := helpers.ByteArrayToInt64(super.S_block_start[:])
			size := helpers.ByteArrayToInt64(super.S_block_size[:])
			bit,e:= strconv.ParseInt(string(directo),36,64); if e!=nil{msg_error(e)}
			bloque := helpers.ReadFileBlock(disco,bit*size+inicio)
			contenido += string(bloque.B_content[:])
		}
	}
	return contenido
}
func buscarUsuario(contenido string, user string, password string, comprobar bool) (bool,int,int){
	registros := strings.Split(contenido,"\n")
	offset := 0
	index := 0
	for i := 0; i < len(registros)-1; i++ {
		datos := strings.Split(registros[i],",")
		if datos[0]!="0"{
			x,e:= strconv.Atoi(datos[0]); if e!=nil{msg_error(e)}
			if x>index{index = x}
			if datos[1]=="U"{
				if datos[2] == user && !comprobar{
					return true,index,offset
				}else if datos[2] == user && datos[4] == password{
					return true,index,offset
				}
			}
		}
		offset += len(registros[i])+1
	}
	return false,index,offset
}
func buscarGrupo(contenido string,grupo string) (bool,int,int){
	registros := strings.Split(contenido,"\n")
	index := 0
	i :=0
	offset :=0
	for i = 0; i < len(registros)-1; i++ {
		datos := strings.Split(registros[i],",")
		if datos[0]!="0"{
			x,e:= strconv.Atoi(datos[0]); if e!=nil{msg_error(e)}
			if x>index{index = x}
			if datos[1]=="G"{
				if datos[2] == grupo{
					return true,index,offset
				}
			}
		}
		offset += len(registros[i])+1
	}
	return false,index,offset
}
func Logout(parametros []string)string{
	consola := "==========LOGOUT==========\n"
	if len(parametros)!=1{
		return consola + "Error: Este comando no necesita parámetros\n"
	}
	if sesionInicida{
		sesionInicida = false
		usuarioActual = estructuras.Usuario{}
		particionActual = estructuras.Pmontada{}
		return consola + "Sesión cerrada con éxito\n"
	}
	return consola + "Error: No hay una sesión activa"
}
func MkRmgrp(parametros []string,remover bool) string{
	consola := "==========MKGRP==========\n"
	if remover{consola = "==========RMGRP==========\n"}
	cant := 0
	nombreGrupo := ""
	if !sesionInicida{ return consola + "Error: Debe iniciar sesión"}
	for i := 1; i < len(parametros); i++ {
		parametro := strings.ToLower(parametros[i])
		if strings.Contains(parametro,"name="){
			nombreGrupo = strings.Split(parametros[i],"=")[1]
			nombreGrupo = strings.ReplaceAll(nombreGrupo,"\"","")
			cant++
			if len(nombreGrupo)>10{cant--; consola += "Error: El nombre de grupo es mayor que 10\n"}
		}else{
			cant = -1
			consola += "Error: Parámetro "+parametro+" no válido\n"
		}
	}
	if cant == 1{
		if usuarioActual.User !="root" { consola+= "Error: Debe ser usuario root para crear un grupo"
		}else{
			if !remover{
				consola += crearGrupo(nombreGrupo)
			}else{
				consola += removerGrupo(nombreGrupo)
			}
		}
	}else{
		consola += "Error: Faltan parámetros obligatorios\n"
	}
	return consola
}
func crearGrupo(nombreGrupo string) string{
	consola := "==========CREANDO GRUPO==========\n"
	disco,e := os.OpenFile(particionActual.Path,os.O_RDWR,0660); if e!=nil{msg_error(e)}
	superBloque := helpers.ReadSuperBlock(disco,helpers.ByteArrayToInt64(particionActual.Particion.Part_start[:]))
	carpetas := strings.Split("/users.txt","/")
	pos := buscarArchivo(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),&superBloque,disco,carpetas,1)
	if pos == -1{consola +="Error: No se encontró el archivo/Directorio"+"/users.txt\n"}
	inodo := helpers.ReadInode(disco,int64(pos))
	contenido := contenidoArchivo(inodo,&superBloque,disco)
	fmt.Println(len(contenido))
	existe,index,init:= buscarGrupo(contenido,nombreGrupo)
	init++
	if existe {consola += "Error: El grupo "+nombreGrupo+" ya existe\n"
	}else{
		data := strconv.Itoa(index+1)+",G,"+nombreGrupo+"\n"
		sizeInodo := helpers.ByteArrayToInt64(inodo.I_size[:])
		offset := sizeInodo%64
		directo := sizeInodo/64
		inicio := helpers.ByteArrayToInt64(superBloque.S_block_start[:])
		size := helpers.ByteArrayToInt64(superBloque.S_block_size[:])
		bit,e:= strconv.ParseInt(string(inodo.I_block[directo]),36,64);if e!=nil{msg_error(e)		}
		bloque := helpers.ReadFileBlock(disco,bit*size+inicio)
		copy(bloque.B_content[offset:],data)
		sizeData :=int64(len(data))
		sizeInodo+=sizeData
		copy(inodo.I_size[:],helpers.IntToByteArray(sizeInodo))
		puntero,e:= disco.Seek(bit*size+inicio,io.SeekStart);if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(bloque),puntero)
		
		sobrante:=  64-(offset+sizeData)
		if sobrante<0{
			nuevoBloque := estructuras.BloqueArchivos{}
			copy(nuevoBloque.B_content[:],[]byte(data)[(sizeData+sobrante):])
			InicioBitmapBloque := helpers.ByteArrayToInt64(superBloque.S_bm_block_start[:])
			BloqueLibreBM  := helpers.ByteArrayToInt64(superBloque.S_first_blo[:])
			BloqueLibre := helpers.ByteArrayToInt64(superBloque.S_block_start[:])+helpers.ByteArrayToInt64(superBloque.S_block_size[:])*(BloqueLibreBM)
			
			//Actualizar bitmap, primer libre y cantidad libre
			actualizarBitmapBloque(disco,helpers.ByteArrayToInt64(superBloque.S_blocks_count[:]),InicioBitmapBloque,&superBloque)
			nuevoCantBloqueLibre := helpers.ByteArrayToInt64(superBloque.S_free_blocks_count[:]) - 1
			copy(superBloque.S_free_blocks_count[:],helpers.IntToByteArray(nuevoCantBloqueLibre))
			
			inodo.I_block[directo+1] = helpers.IntToByteArray(BloqueLibreBM)[0]

			puntero,e= disco.Seek(BloqueLibre,io.SeekStart);if e!=nil{msg_error(e)}
			disco.WriteAt(Struct_to_bytes(nuevoBloque),puntero)

			puntero,e= disco.Seek(helpers.ByteArrayToInt64(particionActual.Particion.Part_start[:]),io.SeekStart);if e!=nil{msg_error(e)}
			disco.WriteAt(Struct_to_bytes(superBloque),puntero)
		}
		puntero,e= disco.Seek(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),io.SeekStart);if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(inodo),puntero)
		consola += "Grupo "+nombreGrupo+" creado con éxito"
	}
	disco.Close()
	return consola
}
func removerGrupo(nombreGrupo string) string{
	consola := "==========REMOVIENDO GRUPO==========\n"
	disco,e := os.OpenFile(particionActual.Path,os.O_RDWR,0660); if e!=nil{msg_error(e)}
	superBloque := helpers.ReadSuperBlock(disco,helpers.ByteArrayToInt64(particionActual.Particion.Part_start[:]))
	carpetas := strings.Split("/users.txt","/")
	pos := buscarArchivo(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),&superBloque,disco,carpetas,1)
	if pos == -1{consola +="Error: No se encontró el archivo/Directorio"+"/users.txt\n"}
	inodo := helpers.ReadInode(disco,int64(pos))
	contenido := contenidoArchivo(inodo,&superBloque,disco)
	fmt.Println(len(contenido))
	existe,index,init:= buscarGrupo(contenido,nombreGrupo)
	index++ // no se usa xd
	if !existe{
		consola += "El grupo "+ nombreGrupo+" no existe"
	}else{

		offset := init%64
		pos := init/64
		directo,e:= strconv.ParseInt(string(inodo.I_block[pos]),36,64);if e!=nil{msg_error(e)}
		sizeB := helpers.ByteArrayToInt64(superBloque.S_block_size[:])
		inicioBloque := helpers.ByteArrayToInt64(superBloque.S_block_start[:])
		bloque := helpers.ReadFileBlock(disco,inicioBloque+int64(directo)*sizeB)
		bloque.B_content[offset] = '0'
		puntero,e:= disco.Seek(inicioBloque+int64(directo)*sizeB,io.SeekStart); if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(bloque),puntero)
		consola+="Grupo "+nombreGrupo+" eliminado con éxito"
	}
	disco.Close()
	return consola
}
func MkRmusr(parametros []string,remover bool) string{
	consola := "==========MKUSR==========\n"
	if remover{consola = "==========RMUSR==========\n"}
	if !sesionInicida{ return consola + "Error: Debe iniciar sesión"}
	cant := 0
	user :=""
	pwd := ""
	grp := ""
	for i := 1; i < len(parametros); i++ {
		parametro := strings.ToLower(parametros[i])
		if strings.Contains(parametro,"usuario="){
			user = strings.Split(parametros[i],"=")[1]
			cant++
			if len(user)>10{cant--;consola+="Error: Carácteres de usuario mayor que 10\n"}
		}else if strings.Contains(parametro,"pwd=") && !remover{
			pwd = strings.Split(parametros[i],"=")[1]
			cant++
			if len(pwd)>10{cant--; consola+="Error: Carácteres de usuario mayor que 10\n"}

		}else if strings.Contains(parametro,"grp=") && !remover{
			grp = strings.Split(parametros[i],"=")[1]
			cant++
			if len(grp)>10{cant--; consola+="Error: Carácteres de usuario mayor que 10\n"}
			
		}else{
			consola += "Parámetro "+parametro+" no válido\n"
		}
	}
	if cant==3 && !remover{
		if usuarioActual.User != "root"{
			consola += "Error: Debe ser usuario root\n"
		}else{
			consola += crearUsuario(user,pwd,grp)
		}
	}else if cant==1 && remover{
		if usuarioActual.User != "root"{
			consola += "Error: Debe ser usuario root\n"
		}else{
			consola += removerUsuario(user)
		}
	}else{
		consola += "Error: Cantidad de parámetros no válidos"
	}
	return consola
}
func crearUsuario(user string,pwd string,grp string) string{
	consola := "==========CREANDO USUARIO==========\n"
	disco,e := os.OpenFile(particionActual.Path,os.O_RDWR,0660); if e!=nil{msg_error(e)}
	superBloque := helpers.ReadSuperBlock(disco,helpers.ByteArrayToInt64(particionActual.Particion.Part_start[:]))
	carpetas := strings.Split("/users.txt","/")
	pos := buscarArchivo(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),&superBloque,disco,carpetas,1)
	if pos == -1{consola +="Error: No se encontró el archivo/Directorio"+"/users.txt\n"}
	inodo := helpers.ReadInode(disco,int64(pos))
	contenido := contenidoArchivo(inodo,&superBloque,disco)
	existegrp,indexgrp,initgrp:= buscarGrupo(contenido,grp)
	indexgrp++;initgrp++
	if !existegrp{ consola +="Error: El grupo "+grp+" no existe\n"
	}else{
	existeuser,indexuser,inituser:= buscarUsuario(contenido,user,pwd,false)
	inituser++
	if existeuser{ consola += "Error: El usuario "+user+" ya existe\n"
		}else{
		data := strconv.Itoa(indexuser+1)+",U,"+user+","+grp+","+pwd+"\n"
		sizeInodo := helpers.ByteArrayToInt64(inodo.I_size[:])
		offset := sizeInodo%64
		directo := sizeInodo/64
		inicio := helpers.ByteArrayToInt64(superBloque.S_block_start[:])
		size := helpers.ByteArrayToInt64(superBloque.S_block_size[:])
		bit,e:= strconv.ParseInt(string(inodo.I_block[directo]),36,64);if e!=nil{msg_error(e)		}
		bloque := helpers.ReadFileBlock(disco,bit*size+inicio)
		copy(bloque.B_content[offset:],data)
		sizeData :=int64(len(data))
		sizeInodo+=sizeData
		copy(inodo.I_size[:],helpers.IntToByteArray(sizeInodo))
		puntero,e:= disco.Seek(bit*size+inicio,io.SeekStart);if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(bloque),puntero)
		
		sobrante:=  64-(offset+sizeData)
		if sobrante<0{
			nuevoBloque := estructuras.BloqueArchivos{}
			copy(nuevoBloque.B_content[:],[]byte(data)[(sizeData+sobrante):])
			InicioBitmapBloque := helpers.ByteArrayToInt64(superBloque.S_bm_block_start[:])
			BloqueLibreBM  := helpers.ByteArrayToInt64(superBloque.S_first_blo[:])
			BloqueLibre := helpers.ByteArrayToInt64(superBloque.S_block_start[:])+helpers.ByteArrayToInt64(superBloque.S_block_size[:])*(BloqueLibreBM)
			
			//Actualizar bitmap, primer libre y cantidad libre
			actualizarBitmapBloque(disco,helpers.ByteArrayToInt64(superBloque.S_blocks_count[:]),InicioBitmapBloque,&superBloque)
			nuevoCantBloqueLibre := helpers.ByteArrayToInt64(superBloque.S_free_blocks_count[:]) - 1
			copy(superBloque.S_free_blocks_count[:],helpers.IntToByteArray(nuevoCantBloqueLibre))
			
			inodo.I_block[directo+1] = helpers.IntToByteArray(BloqueLibreBM)[0]

			puntero,e= disco.Seek(BloqueLibre,io.SeekStart);if e!=nil{msg_error(e)}
			disco.WriteAt(Struct_to_bytes(nuevoBloque),puntero)

			puntero,e= disco.Seek(helpers.ByteArrayToInt64(particionActual.Particion.Part_start[:]),io.SeekStart);if e!=nil{msg_error(e)}
			disco.WriteAt(Struct_to_bytes(superBloque),puntero)
		}
		puntero,e= disco.Seek(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),io.SeekStart);if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(inodo),puntero)
		consola += "Usuario "+user+" creado con éxito"
		}
	}
	return consola
}
func removerUsuario(user string)string{
	consola := "==========REMOVIENDO USUARIO==========\n"
	disco,e := os.OpenFile(particionActual.Path,os.O_RDWR,0660); if e!=nil{msg_error(e)}
	superBloque := helpers.ReadSuperBlock(disco,helpers.ByteArrayToInt64(particionActual.Particion.Part_start[:]))
	carpetas := strings.Split("/users.txt","/")
	pos := buscarArchivo(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),&superBloque,disco,carpetas,1)
	if pos == -1{consola +="Error: No se encontró el archivo/Directorio"+"/users.txt\n"}
	inodo := helpers.ReadInode(disco,int64(pos))
	contenido := contenidoArchivo(inodo,&superBloque,disco)
	fmt.Println(len(contenido))
	existe,index,init:= buscarUsuario(contenido,user,"",false)
	index++ // no se usa xd
	if !existe{
		consola += "El usuario "+ user+" no existe"
	}else{
		offset := init%64
		pos := init/64
		directo,e:= strconv.ParseInt(string(inodo.I_block[pos]),36,64);if e!=nil{msg_error(e)}
		sizeB := helpers.ByteArrayToInt64(superBloque.S_block_size[:])
		inicioBloque := helpers.ByteArrayToInt64(superBloque.S_block_start[:])
		bloque := helpers.ReadFileBlock(disco,inicioBloque+int64(directo)*sizeB)
		bloque.B_content[offset] = '0'
		puntero,e:= disco.Seek(inicioBloque+int64(directo)*sizeB,io.SeekStart); if e!=nil{msg_error(e)}
		disco.WriteAt(Struct_to_bytes(bloque),puntero)
		consola+="Usuario "+user+" eliminado con éxito"
	}
	disco.Close()
	return consola
}
func Rep(parametros []string)string{
	consola := "==========REP==========\n"
	cant:=0
	nombre := ""
	ruta :=""
	id :=""
	path :=""
	for i := 1; i < len(parametros); i++ {
		parametro := strings.ToLower(parametros[i])
		if strings.Contains(parametro,"name="){
			nombre = strings.Replace(parametro,"name=","",1)
			cant++
		}else if strings.Contains(parametro,"path="){
			path = strings.Replace(parametro,"path=","",1)
			path = strings.ReplaceAll(path,"\"","")
			cant++
		}else if strings.Contains(parametro,"id="){
			id = strings.Replace(parametro,"id=","",1)
			cant++
		}else if strings.Contains(parametro,"ruta="){
			ruta = strings.Replace(parametro,"ruta=","",1)
			ruta = strings.ReplaceAll(ruta,"\"","")
		}else{
			consola += "Error: Parámetro "+parametro+" no válido\n"
			cant = -1
		}
	}
	if cant==3{
		if nombre == "file" && ruta ==""{consola += "Error: Falta parámetro ruta para reporte File\n"
		}else{
				switch nombre {
			case "disk":
			case "tree":
			case "file":
				consola += repFile(id,ruta)
			case "sb":
				consola += repSb(id)
			default:
				consola += "Error: Reporte "+nombre+" no existe \n"
			}
		}
	}else{
		consola += "Error: Cantidad de parámetros obligatorios no válida"
	}
	return consola
}
func repFile(id string,ruta string)string{
	consola := "==========Reporte File==========\n"
	part := estructuras.Pmontada{}
	consola += getParticionMontada(id,&part)
	if strings.Contains(consola,"no encontrada"){return consola}
	disco,e := os.OpenFile(part.Path,os.O_RDWR,0660); if e!=nil{msg_error(e)}
	superBloque := helpers.ReadSuperBlock(disco,helpers.ByteArrayToInt64(part.Particion.Part_start[:]))
	carpetas := strings.Split(ruta,"/")
	pos := buscarArchivo(helpers.ByteArrayToInt64(superBloque.S_inode_start[:]),&superBloque,disco,carpetas,1)
	if pos == -1{consola +="Error: No se encontró el archivo/Directorio"+ruta+"\n"}
	inodo := helpers.ReadInode(disco,int64(pos))
	contenido := contenidoArchivo(inodo,&superBloque,disco)
	RepFile = "digraph File{node[shape=\"rectangle\"]; \"file\"[label = \""+ruta+"\n"+contenido+" \"] \"Reporte File\"	rankdir = LR}"
	RepFile = strings.ReplaceAll(RepFile,"\n","\\n")
	RepFile = strings.Replace(RepFile, "\x00", "", -1)
	consola += "Reporte Creado"
	return consola
}
func repSb(id string)string{
	consola := "==========ReporteSB==========\n"
	part := estructuras.Pmontada{}
	consola += getParticionMontada(id,&part)
	if strings.Contains(consola,"no encontrada"){return consola}
	disco,e := os.OpenFile(part.Path,os.O_RDWR,0660); if e!=nil{msg_error(e)}
	superBloque := helpers.ReadSuperBlock(disco,helpers.ByteArrayToInt64(part.Particion.Part_start[:]))
    grafo := "digraph tabla{\n";
    grafo += "abc [shape=none, margin=0, label=<\n";
    grafo += "<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n";
    grafo +="<TR><TD colspan=\"2\">Reporte SuperBloque</TD></TR>";
    grafo +="<TR><TD>s_filesystem_type</TD><TD>"+string(superBloque.S_filesystem_type[:])+"</TD></TR>";
    grafo +="<TR><TD>s_inodes_count</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_inodes_count[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_blocks_count</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_blocks_count[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_free_blocks_count</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_free_blocks_count[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_free_inodes_count</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_free_inodes_count[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_mtime</TD><TD>"+string(superBloque.S_mtime[:])+"</TD></TR>";
    grafo +="<TR><TD>s_mnt_count</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_mnt_count[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_magic</TD><TD>"+string(superBloque.S_magic[:])+"</TD></TR>";
    grafo +="<TR><TD>s_inode_s</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_inode_size[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_block_s</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_block_size[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_firts_ino</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_firts_ino[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_first_blo</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_first_blo[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_bm_inode_start</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_bm_inode_start[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_bm_block_start</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_bm_block_start[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_inode_start</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_inode_start[:])))+"</TD></TR>";
    grafo +="<TR><TD>s_block_start</TD><TD>"+strconv.Itoa(int(helpers.ByteArrayToInt64(superBloque.S_block_start[:])))+"</TD></TR>";
    grafo += "</TABLE>>];\n}";
	RepSB = strings.Replace(grafo, "\x00", "", -1)
	consola += "Reporte Creado"
	return consola
}