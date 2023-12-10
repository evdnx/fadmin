package util

import (
	"encoding/json"
	"math/rand"
)

// RandomString generates a random string of a given length
func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[]{}|;:',.<>/?")

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

// InterfaceToByte converts an interface to byte array
func InterfaceToByte(i any) []byte {
	switch i := i.(type) {
	case []byte:
		return i
	case string:
		return []byte(i)
	default:
		b, _ := json.Marshal(i)
		return b
	}
}
