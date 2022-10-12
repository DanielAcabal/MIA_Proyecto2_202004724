package analizador

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"MIA-Proyecto2_202004724/comando"	
)

func Analizar() {
	finalizar := false
	fmt.Println("=====MIA-Proyecto2_202004724=====")
	reader := bufio.NewReader(os.Stdin)// Buffer de entrada
	for !finalizar { 
		fmt.Print("MIA-Proyecto2:$ ")
		comando, _ := reader.ReadString('\n')// Obtenemos la entrada hasta encontrar ,\n
											 // Se puede usar fmt.Scann pero solo obtiene una palabra
		if strings.Contains(comando, "exit") {
			finalizar = true
		} else {
			if comando != "" && comando != "exit\n" {
				//  Separacion de comando y parametros
				split_comand(comando)
			}
		}
	}
}
func split_comand(command string) {
	var commandArray []string
	// Clean command
	command = strings.Replace(command, "\n", "", 1) //reemplazamos una vez, n<0 para todos
	command = strings.Replace(command, "\r", "", 1)
	// Save id and params
	if strings.Contains(command, "mostrar") { // Commands without paramaters
		commandArray = append(commandArray, command)
	} else {
		commandArray = strings.Split(command, " ") //[id,param1,param2,param3]
	}
	execComand(commandArray)
}
func execComand(command []string){
	switch command[0] {
	case "exec":
	case "mkdisk":
		comando.Mkdisk(command)
	}
}