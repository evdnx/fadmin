package cmd

import (
	"fmt"
	"os/exec"

	"github.com/evdnx/unixmint/auth"
)

func Exec(command string) *exec.Cmd {
	c := fmt.Sprint("echo ", auth.Password(), " | su - ", auth.Username(), " -c ", command)
	return exec.Command("sh", "-c", c)
}
