package aes

import (
    "encoding/base64"
    "fmt"
    "testing"
)

var key = "9553B3B5FC4A2269"

func TestEncryptDecryptBase64(t *testing.T) {
    pass := "99900000000"
    //pass := []byte(strings.Repeat("9", (11+1)*100))
    ciphertext, err := CBCEncryptBase64(pass, key)
    if err != nil {
        fmt.Println(err)
        return
    }
    println("ciphertext: ", ciphertext)

    // Decrypt
    plaintext, err := CBCDecryptBase64(ciphertext, key)
    if err != nil {
        fmt.Println(err)
        return
    }
    println("plaintext: ", plaintext)
}

func TestDecryptBase64(t *testing.T) {
    ciphertext := "h8oOMFjtw72HgFXTOGciXA=="

    // Decrypt
    plaintext, err := CBCDecryptBase64(ciphertext, key)
    if err != nil {
        fmt.Println(err)
        return
    }
    println("plaintext: ", plaintext)
}

func TestEncryptDecrypt(t *testing.T) {
    pass := []byte("99900000000")
    //pass := []byte(strings.Repeat("9", (11+1)*100))
    ciphertext, err := CBCEncrypt(pass, []byte(key))
    if err != nil {
        fmt.Println(err)
        return
    }

    ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)
    fmt.Printf("ciphertext:%s\n", ciphertextBase64)

    bytesPass, err := base64.StdEncoding.DecodeString(ciphertextBase64)
    if err != nil {
        fmt.Println(err)
        return
    }

    plaintext, err := CBCDecrypt(bytesPass, []byte(key))
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("plaintext:%s\n", plaintext)
}
