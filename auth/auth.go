package auth

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/evdnx/unixmint/pkg/crypt"
	"github.com/google/uuid"
)

var pwdFileName string = ""

func Login(username, password string) error {
	command := fmt.Sprint("echo ", password, " | su - ", username, " -c ", `"echo 1"`)
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// encrypt password
	encryptedPassword, err := crypt.Encrypt(password)
	if err != nil {
		return err
	}

	// generate file name
	pwdFileName = uuid.NewString()

	// write encrypted password to a temporary file
	err = os.WriteFile(fmt.Sprint("/tmp/", pwdFileName), []byte(encryptedPassword), 0644)
	if err != nil {
		return err
	}

	return nil
}

func Username() string {
	return ""
}

func Password() string {
	// read password file
	pwd, err := os.ReadFile(pwdFileName)
	if err != nil {
		return ""
	}

	// decrypt password
	password, err := crypt.Decrypt(pwd)
	if err != nil {
		return ""
	}

	return password
}
