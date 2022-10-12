package comando

import (
	"MIA-Proyecto2_202004724/Estructuras"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"regexp"
)

/*======================MKDISK=======================*/
func Mkdisk(commandArray []string) {
	//mkdisk -Size=3000 -unit=K -path=/home/user2/Disco1.dk
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
		if strings.Contains(data, "-size=") {
			strtam := strings.Replace(data, "-size=", "", 1)
			strtam = strings.Replace(strtam, "\"", "", 2)
			strtam = strings.Replace(strtam, "\r", "", 1)
			tamano2, err := strconv.Atoi(strtam)
			tamano = tamano2
			if err != nil {
				msg_error(err)
			}
			if tamano<=0{
				fmt.Print("El tamaño debe ser mayor o igual a cero")
				seguir = false
			}
			cantidad +=1
		} else if strings.Contains(data, "-unit=") {
			dimensional = strings.Replace(data, "-unit=", "", 1)
			dimensional = strings.Replace(dimensional, "\"", "", 2)
		}else if strings.Contains(data,"-fit="){
			ajuste = strings.Replace(data,"-fit=","",1)
			ajuste = strings.Replace(ajuste, "\"", "", 2)
			if ajuste!="bf" || ajuste!="ff" || ajuste!="wf"{
				seguir = false
				fmt.Print("Ajuste no válido")
			}
		}else if strings.Contains(data,"-path="){
			ruta = strings.Replace(data,"-path=","",1)
			ruta = strings.Replace(ruta, "\"", "", 2)
			reg := regexp.MustCompile("/[a-zA-Z0-9]+.dk")
			junto := reg.ReplaceAllString(ruta,"")
			junto = "/home"+junto
			err := os.MkdirAll(junto,os.ModePerm)
			if err != nil {
				msg_error(err)
			}

			cantidad+=1
		}else{
			fmt.Print("Parámetro no válido:", data)
		}
	}
	if !seguir{
		return
	}
	if cantidad!=2{
		return
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
	ruta = "/home"+ruta
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
	puntero,er := disco.Seek(0,os.SEEK_SET)
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
	fmt.Print("Creacion de Disco:")
	fmt.Print(" Tamaño: ")
	fmt.Print(tamano)
	fmt.Print(" Dimensional: ")
	fmt.Println(dimensional)
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