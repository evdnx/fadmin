package auth_svc

import (
	"fmt"
	"os/exec"
)

func Login(username, password string) bool {
	command := fmt.Sprint("echo ", password, " | su - ", username, " -c ", `"echo 1"`)
	cmd := exec.Command("bash", "-c", command)
	return cmd.Run() == nil
}
