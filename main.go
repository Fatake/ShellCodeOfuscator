package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
	"os"
	"time"

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
	color.Cyan(banner)
	menu := opciones()
	if menu.inputFile == "" {
		log.Fatal("[!] Requiere binario, user -h para ver la ayuda")
	}

	dataFile, err := os.ReadFile(menu.inputFile)
	if err != nil {
		log.Fatal(err)
	}

	////
	// XOR Part
	////
	color.Blue("[i] Realizando proceso de cifrado Xor")
	cifrar := coder.XorEncoder(dataFile, 60)

	// obs es la variable que contiene el resultado final de la obfuscacion
	obs := coder.Base32CustomEncoder(string(cifrar))

	os.WriteFile("obs.txt", []byte(obs), 0777)
	color.Yellow("[+] Archivo Ofuscado Exitosamente")

	////
	// AES Part
	////
	color.Blue("[i] Realizando proceso de cifrado AES")
	key := make([]byte, 32)
	nonce := make([]byte, 12)

	rand.Read(key)
	rand.Read(nonce)
	time.Sleep(2 * time.Second)
	color.Green("[+] Key: " + base64.StdEncoding.EncodeToString(key))
	color.Green("[+] nonce: " + base64.StdEncoding.EncodeToString(nonce))

}
