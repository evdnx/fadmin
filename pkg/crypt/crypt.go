package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// Encrypt encrypts a message
func Encrypt(key []byte, message string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainText := []byte(message)
	cipherText := make([]byte, len(plainText))

	iv := cipherText[:aes.BlockSize]
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts a message
func Decrypt(key []byte, message string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cfb := cipher.NewCFBDecrypter(block, iv)

	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
