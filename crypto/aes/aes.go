package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func CBCEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	padding := blockSize - len(plaintext)%blockSize
	padBytes := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext = append(plaintext, padBytes...)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func CBCEncryptBase64(plaintext, key string) (string, error) {
	ciphertext, err := CBCEncrypt([]byte(plaintext), []byte(key))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func CBCDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	plaintext := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(plaintext, ciphertext)

	length := len(plaintext)
	unPadding := int(plaintext[length-1])
	plaintext = plaintext[:(length - unPadding)]
	return plaintext, nil
}

func CBCDecryptBase64(ciphertext, key string) (string, error) {
	decodeCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := CBCDecrypt(decodeCiphertext, []byte(key))
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
