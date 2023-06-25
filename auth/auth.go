package auth

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/evdnx/unixmint/pkg/crypt"
	"github.com/evdnx/unixmint/pkg/util"
	"github.com/google/uuid"
)

var keyFileName string = ""
var pwdFileName string = ""

func Login(username, password string) error {
	command := fmt.Sprint("echo ", password, " | su - ", username, " -c ", `"echo 1"`)
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// generate encryption key
	key := util.RandomString(64)

	// encrypt password
	encryptedPassword, err := crypt.Encrypt([]byte(key), password)
	if err != nil {
		return err
	}

	// generate file names
	keyFileName = uuid.NewString()
	pwdFileName = uuid.NewString()

	// write encryption key to a temporary file
	err = os.WriteFile(fmt.Sprint("/tmp/", keyFileName), []byte(key), 0644)
	if err != nil {
		return err
	}

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
	// read key file
	key, err := os.ReadFile(keyFileName)
	if err != nil {
		return ""
	}

	// read password file
	pwd, err := os.ReadFile(pwdFileName)
	if err != nil {
		return ""
	}

	// decrypt password
	password, err := crypt.Decrypt(key, string(pwd))
	if err != nil {
		return ""
	}

	return password
}
