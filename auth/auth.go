package auth

import (
	"fmt"
	"os/exec"
)

var Username string = ""
var Password string = ""

func Login(username, password string) bool {
	command := fmt.Sprint("echo ", password, " | su - ", username, " -c ", `"echo 1"`)
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		return false
	}

	Username = username
	Password = password

	return true
}
