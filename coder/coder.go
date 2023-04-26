package coder

import (
	"encoding/base32"
	"strings"
)

func XorEncoder(buffer []byte, xorkey byte) []byte {
	resultado := make([]byte, len(buffer))
	for i := 0; i < len(buffer); i++ {
		resultado[i] = xorkey ^ buffer[i]
	}

	return resultado
}

func Base32CustomEncoder(codigo string) string {
	pay := base32.StdEncoding.EncodeToString([]byte(codigo))
	rep := strings.Replace(pay, "Q", "[/___[*]]", -1)
	return rep
}
