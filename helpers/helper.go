package helpers

import (
	estructuras "MIA-Proyecto2_202004724/Estructuras"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"unsafe"
	//"unsafe"
)

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
func msg_error(err error){
	fmt.Println("Error: ",err)
}
func HandleSizeof(p interface{})int64{
	return int64(len(Struct_to_bytes(p)))
}
func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	cos := strconv.FormatInt(num, 36)
	for i := 0; i < 2; i++ {
		//byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		
		//arr[i] = byt
	}
	copy(arr,cos)
	return arr
}
func ByteArrayToInt64(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		if arr[i]==0{break}
		val = val*36
		x,e:= strconv.ParseInt(string(arr[i]),36,64)
		if e!=nil{
			msg_error(e)
		}
		val += x
	}
	return val
}
func ByteArrayToFloat64(arr []byte) float64 {
	val := float64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		if arr[i]==0{break}
		val = val*36
		x,e:= strconv.ParseInt(string(arr[i]),36,64)
		if e!=nil{
			msg_error(e)
		}
		val += float64(x)
	}
	return val
}
func ByteArrayToInt(arr []byte) int{
	val := 0
	size := len(arr)
	for i := 0; i < size; i++ {
		if arr[i]==0{break}
		val = val*10
		val += int(arr[i])-48
	}
	return val
}
func ByteArrayToInode(arr []byte) estructuras.Inodo {
	// Decodificacion de [] Bytes a Struct ejemplo
	p := estructuras.Inodo{}
	dec := gob.NewDecoder(bytes.NewReader(arr))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}
func ByteArrayToDirBlock(arr []byte) estructuras.BloqueCarpeta {
	// Decodificacion de [] Bytes a Struct ejemplo
	p := estructuras.BloqueCarpeta{}
	dec := gob.NewDecoder(bytes.NewReader(arr))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}
func ByteArrayToFileBlock(arr []byte) estructuras.BloqueArchivos {
	// Decodificacion de [] Bytes a Struct ejemplo
	p := estructuras.BloqueArchivos{}
	dec := gob.NewDecoder(bytes.NewReader(arr))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}
func ByteArrayToSuperBlock(arr []byte) estructuras.SuperBloque {
	// Decodificacion de [] Bytes a Struct ejemplo
	p := estructuras.SuperBloque{}
	dec := gob.NewDecoder(bytes.NewReader(arr))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}
func ReadInode(disco *os.File,pos int64) estructuras.Inodo{
	data := make([]byte,HandleSizeof(estructuras.Inodo{}))
	puntero,e :=disco.Seek(pos,io.SeekStart); if e!=nil{msg_error(e)}
	disco.ReadAt(data,puntero)
	return ByteArrayToInode(data)
}
func ReadFileBlock(disco *os.File,pos int64) estructuras.BloqueArchivos{
	data := make([]byte,HandleSizeof(estructuras.BloqueArchivos{}))
	puntero,e :=disco.Seek(pos,io.SeekStart); if e!=nil{msg_error(e)}
	disco.ReadAt(data,puntero)
	return ByteArrayToFileBlock(data)
}
func ReadDirBlock(disco *os.File,pos int64) estructuras.BloqueCarpeta{
	data := make([]byte,HandleSizeof(estructuras.BloqueCarpeta{}))
	puntero,e :=disco.Seek(pos,io.SeekStart); if e!=nil{msg_error(e)}
	disco.ReadAt(data,puntero)
	return ByteArrayToDirBlock(data)
}
func ReadSuperBlock(disco *os.File,pos int64) estructuras.SuperBloque{
	data := make([]byte,HandleSizeof(estructuras.SuperBloque{}))
	puntero,e :=disco.Seek(pos,io.SeekStart); if e!=nil{msg_error(e)}
	disco.ReadAt(data,puntero)
	return ByteArrayToSuperBlock(data)
}
func Round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(Round(num * output)) / output
}
