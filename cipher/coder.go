package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base32"
	"encoding/hex"
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

func HexEncode(code string) string {
	contenido := hex.EncodeToString([]byte(code))
	return contenido
}

func newAead(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aead, nil
}

func AESEncrypt(plain []byte, key, nonce []byte) []byte {
	aead, err := newAead(key)
	if err != nil {
		println(err.Error())
		return nil
	}

	return aead.Seal(plain[:0], nonce, plain, nil)
}

func AESDecrypt(cipher []byte, key, nonce []byte) []byte {
	aead, err := newAead(key)
	if err != nil {
		println(err.Error())
		return nil
	}

	output, err := aead.Open(cipher[:0], nonce, cipher, nil)
	if err != nil {
		println(err.Error())
		return nil
	}

	return output
}
