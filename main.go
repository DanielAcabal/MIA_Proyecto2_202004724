package main

import (
	"MIA-Proyecto2_202004724/analizador"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)
type Contenido struct{
	Entrada string
	Consola string
}
func enableCors(w *http.ResponseWriter){
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func main(){
	http.HandleFunc("/Ejecutar",Analizar)
	fmt.Println("Servidor en puerto 5000!")
	log.Fatal(http.ListenAndServe(":5000",nil))
}
func Analizar(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	reqBody,err  := ioutil.ReadAll(r.Body)
	if err!=nil{
		fmt.Print(err)
	}
	var data Contenido
	json.Unmarshal(reqBody,&data)
	consola :=""
	analizador.Split_comand(data.Entrada,&consola)
	data.Entrada = ""
	data.Consola = consola
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(data)
}
