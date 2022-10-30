package comando

import (
	"MIA-Proyecto2_202004724/Estructuras"
	"MIA-Proyecto2_202004724/helpers"
	"bytes"
	"container/list"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)
var particionesMontadas = list.New()
/*======================MKDISK=======================*/
func Mkdisk(commandArray []string) string{
	//mkdisk -Size=3000 -unit=K -path=/home/user/Disco4.dk
	// mkdisk -size=5 -unit=M -path="/home/mis discos/Disco3.dk"
	consola :=""
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
		aux := make([]byte,50)
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
	consola := ""
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
	return "Partición "+id+" no encontrada"
}
func formatear(id string,tipo string) string{
	particion := estructuras.Pmontada{}
	consola := "Formateando partición "+id+"\n"
	consola += getParticionMontada(id,&particion)
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
		tamanioBloque := helpers.HandleSizeof(estructuras.BloqueArchivos{})
		n := (tamanioParticion-tamanioSuper)/(4+3*tamanioBloque+tamanioInodo);
    	numero := int64(math.Floor(float64(n))); //Cantidad de inodos
		fmt.Print(numero)
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
		fmt.Print()
		carpetaRaiz(superbloque,aux[i],particion.Path)
		archivoUsers(id,"1,G,root\n 1,U,root,root,123,\n",disco,&superbloque,&mbr)
	}else{
		consola += "Partición no encontrada en disco"
	}
	return consola
}
func carpetaRaiz(super estructuras.SuperBloque, particion estructuras.Particion,ruta string){
	disco, err := os.OpenFile(ruta,os.O_RDWR,0660);
	if err != nil{
		msg_error(err)
	}
	//Posiciones para saber dónde escribir
	IniciaBitmapInodo := super.S_bm_inode_start[:]
	InodoLibreBM := super.S_firts_ino[:]
    Inodo_libreI := helpers.ByteArrayToInt64(super.S_inode_start[:])+(helpers.ByteArrayToInt64(InodoLibreBM)*helpers.ByteArrayToInt64(super.S_inode_size[:]));//posicion archivo del inodo libre
	

	inodoRaiz := estructuras.Inodo{}
	estructuras.NuevoInodo(&inodoRaiz,1,1,0,"0","777")
	//Actualizamos Bitmap, inodos libres, primer inodo libre
	actualizarBitmapInodo(disco,helpers.ByteArrayToInt64(super.S_inodes_count[:]),helpers.ByteArrayToInt64(IniciaBitmapInodo),&super)
	nuevoLibre := helpers.ByteArrayToInt64(super.S_free_inodes_count[:])-1
	copy(super.S_free_inodes_count[:],helpers.IntToByteArray(nuevoLibre))

	// Bloque carpetaRaiz
	IniciaBitmapBloque := super.S_bm_block_start[:]
	PrimerBloqueLibre := super.S_first_blo[:]
	BloqueLibre := helpers.ByteArrayToInt64(super.S_block_start[:]) + helpers.ByteArrayToInt64(super.S_block_size[:])*helpers.ByteArrayToInt64(PrimerBloqueLibre)
	
	carpetaR := estructuras.BloqueCarpeta{}
	estructuras.NuevoBloqueCarpeta(&carpetaR)//Constructor xd
	// Los primeros 2 registros del primer apuntador directo del Inodo son la carpeta y carpeta padre
	copy(carpetaR.B_content[0].B_name[:],"/")
	copy(carpetaR.B_content[0].B_inodo[:],InodoLibreBM[:]) // Apuntan al inodo creado antes
	copy(carpetaR.B_content[1].B_name[:],"/")
	copy(carpetaR.B_content[1].B_inodo[:],InodoLibreBM[:])
	//Primer apuntador directo 
	inodoRaiz.I_block[0] = InodoLibreBM[0]

	actualizarBitmapBloque(disco,helpers.ByteArrayToInt64(super.S_blocks_count[:]),
								helpers.ByteArrayToInt64(IniciaBitmapBloque[:]),&super)

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
func archivoUsers(id string,contenido string,disco *os.File,super *estructuras.SuperBloque,mbr *estructuras.MBR)string{
	Inodo := estructuras.Inodo{}
	estructuras.NuevoInodo(&Inodo,1,1,0,"1","")
	carpeta := estructuras.BloqueCarpeta{}
	copy(carpeta.B_content[0].B_name[:],"/")
    copy(carpeta.B_content[0].B_inodo[:],"0") //modificar
	copy(carpeta.B_content[1].B_name[:],"/")
    copy(carpeta.B_content[1].B_inodo[:],"0") //modificar
	copy(carpeta.B_content[2].B_name[:],"users.txt")
    copy(carpeta.B_content[2].B_inodo[:],"1") //modificar
    InodoArchivo := estructuras.Inodo{}
	estructuras.NuevoInodo(&InodoArchivo,1,1,0,"0","") // Ver Permisos y size (se actualiza creo)
	for i := 0; i < len(InodoArchivo.I_block); i++ {
		apuntador := InodoArchivo.I_block[i]
		if apuntador == 0{
			
		}
	}

	fmt.Print()
		
	return ""
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
		puntero, e = disco.Seek(inicioInodo,io.SeekStart)
	}
	puntero, e = disco.Seek(inicioBloque,io.SeekStart)
	if e!=nil{
		msg_error(e)
	}
	for i := 0; i < int(finBloque); i++ {
		disco.WriteAt([]byte{'0'},puntero)
		inicioBloque ++
		puntero, e = disco.Seek(inicioBloque,io.SeekStart)
	}
}/*
func buscarCarpeta(disco *os.File,super estructuras.SuperBloque,ruta string) int{
	inicioInodo := helpers.ByteArrayToInt64(super.S_bm_inode_start[:])
	inicioBloques := helpers.ByteArrayToInt64(super.S_bm_block_start[:])
	puntero , err :=disco.Seek(int64(inicioInodo),io.SeekStart)
	if err!=nil{
		msg_error(err)
	}
	data := Struct_to_bytes(estructuras.Inodo{})
	_,err = disco.ReadAt(data,puntero)
	if err!=nil{
		msg_error(err)
	}
	Inodo := helpers.ByteArrayToInode(data)
	pdirecto := Inodo.I_block[0]
	tipo := Inodo.I_type[0]
	index := 0
	existe := false
	for pdirecto != 0{
		if tipo== 1{
			puntero , err :=disco.Seek(int64(inicioBloques+pdirecto),io.SeekStart)
			if err!=nil{
				msg_error(err)
			}
			data := Struct_to_bytes(estructuras.BloqueCarpeta{})
			_,err = disco.ReadAt(data,puntero)
			if err!=nil{
				msg_error(err)
			}
			bCarpeta := helpers.ByteArrayToDirBlock(data)
			for i := 0; i < 4; i++ {
				content := bCarpeta.B_content[i]
				if content.B_inodo!=0{
					if content.B_name == "" {//&& index == len -1
						//Index==len-1{existe =true; break}
						//Dirección del siguiente Inodo
						//Decodificarlo
						//actualizar pInodo
						// Actualizar index

					}
				}else{
					break
				}
			}
		}
		index++
		pdirecto = Inodo.I_block[index]
	}
	if !existe{

	}
	
}*/