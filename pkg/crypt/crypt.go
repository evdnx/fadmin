package crypt

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"github.com/evdnx/fadmin/db"
	"github.com/evdnx/fadmin/pkg/util"
	"golang.org/x/crypto/chacha20poly1305"
)

var aeadCipher cipher.AEAD = nil

func getCipher() (cipher.AEAD, error) {
	if aeadCipher != nil {
		return aeadCipher, nil
	}

	aead, err := chacha20poly1305.NewX(Key())
	if err != nil {
		return nil, err
	}

	return aead, nil
}

func Key() []byte {
	key, err := db.Read(db.AuthBucket, "crypto_key")
	if err != nil {
		return nil
	}

	return []byte(key)
}

func GenerateKey(length int) []byte {
	key := make([]byte, length)
	rand.Read(key)
	return key
}

func Reset() {
	aeadCipher = nil
}

// Encrypt encrypts a message
func Encrypt(message any) (string, error) {
	msg := util.InterfaceToByte(message)
	aead, err := getCipher()
	if err != nil {
		return "", err
	}

	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(msg)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg := aead.Seal(nonce, nonce, msg, nil)
	return string(encryptedMsg), nil
}

// Decrypt decrypts a message
func Decrypt(encryptedMessage any) (string, error) {
	aead, err := getCipher()
	if err != nil {
		return "", err
	}

	byteMsg := util.InterfaceToByte(encryptedMessage)

	if len(byteMsg) < aead.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := byteMsg[:aead.NonceSize()], byteMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
