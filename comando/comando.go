package comando

import (
	"MIA-Proyecto2_202004724/Estructuras"
	"MIA-Proyecto2_202004724/helpers"
	"bytes"
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
	copy(mbr.Mbr_tamano[:],strconv.Itoa(tamano_archivo*1024))
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
        c :=0
		if bytes.Equal(mbr.Dsk_fit[:],[]byte{'F'}){
			bs:=calcularTamanio(size,unit)
            copy(particion.Part_size[:],strconv.Itoa(bs));
            c = primerAjuste(&mbr,&particion,consola); // 0 == inicio disco
		}else if bytes.Equal(mbr.Dsk_fit[:],[]byte{'B'}){
			bs:=calcularTamanio(size,unit)
            copy(particion.Part_size[:],strconv.Itoa(bs));
            c = mejorAjuste(&mbr,&particion,consola); // 0 == inicio disco
		}else if bytes.Equal(mbr.Dsk_fit[:],[]byte{'W'}){
			bs:=calcularTamanio(size,unit)
            copy(particion.Part_size[:],strconv.Itoa(bs));
            c = PeorAjuste(&mbr,&particion,consola); // 0 == inicio disco
		}
            if (c!=0){
				copy(particion.Part_start[:],strconv.Itoa(c));
			}else{crear=false;}

			if helpers.ByteArrayToInt(aux[0].Part_start[:])==0{
                mbr.Mbr_partition_1 = particion;
            }else if helpers.ByteArrayToInt(aux[1].Part_start[:])==0{
                mbr.Mbr_partition_2 = particion;
            }else if helpers.ByteArrayToInt(aux[2].Part_start[:])==0{
                mbr.Mbr_partition_3 = particion;
            }else if helpers.ByteArrayToInt(aux[3].Part_start[:])==0{
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
func primerAjuste(mbr *estructuras.MBR, particion *estructuras.Particion,consola *string) int{
    var info estructuras.MBR
	comienza := int(unsafe.Sizeof(info));
    aux := [4]estructuras.Particion{};
    aux[0] = mbr.Mbr_partition_1;
    aux[1] = mbr.Mbr_partition_2;
    aux[2] = mbr.Mbr_partition_3;
    aux[3] = mbr.Mbr_partition_4;
    n := 4;
    i, j:=0,0;
    for i = 0; i < n - 1; i++{
        for j = 0; j < n - i - 1; j++{
			actual := helpers.ByteArrayToInt(aux[j].Part_start[:])
			siguiente := helpers.ByteArrayToInt(aux[j+1].Part_start[:])
            if (actual > siguiente){
                aux[j], aux[j + 1] = aux[j + 1] , aux[j] 
			}
		}
	}
    for i = 0; i < n; i++ {
		actual := helpers.ByteArrayToInt(aux[i].Part_start[:])
        espacio := actual - comienza; //int
		necesito := helpers.ByteArrayToInt(particion.Part_size[:])
        if (espacio>=necesito){
			copy(particion.Part_start[:],strconv.Itoa(comienza))
            return comienza;
        }

        if(espacio>=0){
			nuevo := helpers.ByteArrayToInt(aux[i].Part_size[:])
			actual := helpers.ByteArrayToInt(aux[i].Part_start[:])
            comienza = actual+ nuevo;
        }
    }
    if (comienza==0){comienza = int(unsafe.Sizeof(info));}
    x := helpers.ByteArrayToInt(mbr.Mbr_tamano[:]); //int
	nuevo := helpers.ByteArrayToInt(particion.Part_size[:])
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

func mejorAjuste(mbr *estructuras.MBR,particion *estructuras.Particion,consola *string) int{
    
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
            if (helpers.ByteArrayToInt(aux[j].Part_start[:]) > helpers.ByteArrayToInt(aux[j + 1].Part_start[:])){
                aux[j], aux[j + 1] = aux[j + 1] ,aux[j] 
			}
		}
	}
	var info estructuras.MBR
	comienza := int(unsafe.Sizeof(info));
    for i = 0; i < n; i++{
        espacio := helpers.ByteArrayToInt(aux[i].Part_start[:]) - comienza;
        pares[i].Inicio = comienza;
        pares[i].Tamanio = espacio;
        if (espacio>=0){
        comienza = helpers.ByteArrayToInt(aux[i].Part_start[:])+helpers.ByteArrayToInt(aux[i].Part_size[:]);
        }
    }
    if (comienza==0){comienza = int(unsafe.Sizeof(info));}
    espacio := helpers.ByteArrayToInt(mbr.Mbr_tamano[:]) - comienza;
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
    	a := helpers.ByteArrayToInt(particion.Part_size[:]);
        if(x>a){
            return pares[i].Inicio;
        }
    }
        *consola += "No hay espacio para esta particion: "+string(particion.Part_name[:]);
    return 0;
}

func PeorAjuste(mbr *estructuras.MBR, particion *estructuras.Particion,consola *string)int{
    
	aux := [4]estructuras.Particion{};
    aux[0] = mbr.Mbr_partition_1;
    aux[1] = mbr.Mbr_partition_2;
    aux[2] = mbr.Mbr_partition_3;
    aux[3] = mbr.Mbr_partition_4;

    n := 4;
    i, j := 0,0;
    for i = 0; i < n - 1; i++{
        for j = 0; j < n - i - 1; j++{
            if helpers.ByteArrayToInt(aux[j].Part_start[:]) > helpers.ByteArrayToInt(aux[j + 1].Part_start[:]){
                aux[j], aux[j + 1] = aux[j+1], aux[j]
			}
		}
	}
	var info estructuras.MBR
    tamanio := 0; comienza := int(unsafe.Sizeof(info)); start := comienza;
    for i = 0; i < n; i++{
        espacio := helpers.ByteArrayToInt(aux[i].Part_start[:]) - comienza;
        if (espacio>tamanio){
            tamanio = espacio;
            start = comienza;

        }
        if(espacio>=0){
        comienza = helpers.ByteArrayToInt(aux[i].Part_start[:])+helpers.ByteArrayToInt(aux[i].Part_size[:]);
        }
    }
    if (comienza==0){comienza = int(unsafe.Sizeof(info));}

    x := helpers.ByteArrayToInt(mbr.Mbr_tamano[:]); 
    espacio := x - comienza;
    if (espacio>tamanio){
        tamanio = espacio;
        start = comienza;
    }
    if (tamanio>=helpers.ByteArrayToInt(particion.Part_size[:])){
        return start;
    }else{
        *consola += "No hay espacio para esta particion: "+string(particion.Part_name[:]);
    }
    return 0;    
}
func addParticion(path string, size int,unit rune,nombre string,consola *string){
    tamanioBytes := calcularTamanio(size,unit);//Ha sumar o restar
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
            if (helpers.ByteArrayToInt(aux[j].Part_start[:]) > helpers.ByteArrayToInt(aux[j + 1].Part_start[:])){
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
        res := helpers.ByteArrayToInt(aux[i].Part_size[:]) + tamanioBytes; //Positivo o negativo
        if(res<=0){
            *consola += "No quedará espacio en la partición\n";
        }else{
            fin := helpers.ByteArrayToInt(mbr.Mbr_tamano[:])
            if(i+1<4){
                fin =helpers.ByteArrayToInt(aux[i+1].Part_start[:])
            }
            if ((fin-helpers.ByteArrayToInt(aux[i].Part_start[:]))>=res){//Si lo puede guardar
                vacio := make([]byte,10)
				copy(aux[i].Part_size[:],vacio);
				copy(aux[i].Part_size[:],strconv.Itoa(res));
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
        if helpers.ByteArrayToInt(aux[i].Part_start[:])!=0{
            if (aux[i].Part_type[0]=='P'){
                p := string(aux[i].Part_name[:]);
                if (p == name){
                    encontrado = true;
                    break;
                }
            }
        }
    }
    if (encontrado){ //FULL
        nuevo :=estructuras.Particion{};
        ceros := make([]byte,1)
        ceros[0] = 0;
        j:=0;
		puntero, err :=disco.Seek(int64(helpers.ByteArrayToInt(aux[i].Part_start[:])),io.SeekStart)
        if err!=nil{
			msg_error(err)
		}
		for (j!=helpers.ByteArrayToInt(aux[i].Part_size[:])){
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