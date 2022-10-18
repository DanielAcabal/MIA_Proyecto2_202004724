package analizador

import (
	"MIA-Proyecto2_202004724/comando"
	"bufio"
	"os"
	"regexp"
	"strings"
	"fmt"
)

func Analizar(consola *string) {
	finalizar := false
	*consola += "=====MIA-Proyecto2_202004724=====\n"
	reader := bufio.NewReader(os.Stdin)// Buffer de entrada
	for !finalizar { 
		*consola = ""
		fmt.Print("MIA-Proyecto2:$ ")
		comando, _ := reader.ReadString('\n')// Obtenemos la entrada hasta encontrar ,\n
											 // Se puede usar fmt.Scann pero solo obtiene una palabra
		if strings.Contains(comando, "exit") {
			finalizar = true
		} else {
			if comando != "" && comando != "exit\n" {
				//  Separacion de comando y parametros
				split_comand(comando,consola)
			}
		}

		fmt.Print(*consola)
	}
}
func split_comand(command string,consola *string) {
	var commandArray []string
	// Clean command
	command = strings.Replace(command, "\n", "", 1) //reemplazamos una vez, n<0 para todos
	command = strings.Replace(command, "\r", "", 1)
	// Save id and params
	if strings.Contains(command, "mostrar") { // Commands without paramaters
		commandArray = append(commandArray, command)
	} else {
		r :=regexp.MustCompile("[ 	]*-")
		commandArray = r.Split(command,15)
		//commandArray = strings.Split(command, " -") //[id,param1,param2,param3]
	}
	execComand(commandArray,consola)
}
func execComand(command []string,consola *string){
	switch command[0] {
	case "exec":
	case "mkdisk":
		*consola += comando.Mkdisk(command)
	case "rmdisk":
		*consola += comando.Rmdisk(command)
	case "fdisk":
		*consola += comando.Fdisk(command)
	}
}