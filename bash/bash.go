package bash

import (
	"fmt"
	"os/exec"

	"github.com/evdnx/unixmint/auth"
)

func Cmd(command string) *exec.Cmd {
	c := fmt.Sprint("echo ", auth.Password(), " | su - ", auth.Username(), " -c ", command)
	return exec.Command("bash", "-c", c)
}
