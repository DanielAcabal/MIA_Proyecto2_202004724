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
				Split_comand(comando,consola)
			}
		}

		fmt.Print(*consola)
	}
}
func Split_comand(entrada string,consola *string) {
	var commandArray []string

	command1 := strings.Split(entrada,"\n") //For multi commands

	for i := 0; i < len(command1); i++ {
		comando := command1[i]
		if comando == "" { // Command without paramaters
			continue
		} else if strings.Contains(comando,"#") || strings.Contains(comando,"pause") {
			*consola += comando + "\n"
			continue
		}else{
			r :=regexp.MustCompile("[ 	]*-")
			commandArray = r.Split(comando,15) //[id,param1,param2,param3]
		}
		execComand(commandArray,consola)
		*consola += "\n"
	}
	
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
	case "mount":
		*consola += comando.Mount(command)
	case "unmount":
		*consola += comando.Unmount(command)
	case "mkfs":
		*consola += comando.Mkfs(command)
	case "login":
		*consola += comando.Login(command)
	case "logout":
		*consola += comando.Logout(command)
	case "mkgrp":
		*consola += comando.Mkgrp(command)
	default:
		*consola += "Error: No existe el comando" + command[0]
	}
}