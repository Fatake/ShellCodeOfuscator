package main

import (
	"flag"
	"log"
	"os"

	"github.com/Fatake/ShellCodeOfuscator/coder"
	"github.com/fatih/color"
)

type FlagsType struct {
	inputFile string
}

var banner = `

░█████╗░██████╗░███████╗██╗░░░██╗░██████╗░█████╗░░█████╗░████████╗░█████╗░██████╗░
██╔══██╗██╔══██╗██╔════╝██║░░░██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██╔══██╗██╔══██╗
██║░░██║██████╦╝█████╗░░██║░░░██║╚█████╗░██║░░╚═╝███████║░░░██║░░░██║░░██║██████╔╝
██║░░██║██╔══██╗██╔══╝░░██║░░░██║░╚═══██╗██║░░██╗██╔══██║░░░██║░░░██║░░██║██╔══██╗
╚█████╔╝██████╦╝██║░░░░░╚██████╔╝██████╔╝╚█████╔╝██║░░██║░░░██║░░░╚█████╔╝██║░░██║
░╚════╝░╚═════╝░╚═╝░░░░░░╚═════╝░╚═════╝░░╚════╝░╚═╝░░╚═╝░░░╚═╝░░░░╚════╝░╚═╝░░╚═╝
by Fatake

`

func opciones() *FlagsType {
	inputFile := flag.String("f", "", "Ruta del binario a ofuscar")
	flag.Parse()
	return &FlagsType{*inputFile}
}

func main() {
	color.Green(banner)
	menu := opciones()
	if menu.inputFile == "" {
		log.Fatal("[!] Requiere binario, user -h para ver la ayuda")
	}

	dataFile, err := os.ReadFile(menu.inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// se cifra la carga parametrizada
	cifrar := coder.XorEncoder(dataFile, 60)

	// obs es la variable que contiene el resultado final de la obfuscacion
	obs := coder.Base32CustomEncoder(string(cifrar))

	os.WriteFile("obs.txt", []byte(obs), 0777)
	color.Yellow("[+] Archivo Ofuscado Exitosamente")

}
