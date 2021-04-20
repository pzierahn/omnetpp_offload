package shell

import (
	"os/exec"
	"strings"
)

func Command(name string, arg ...string) (cmd *exec.Cmd) {

	parts := []string{
		name,
	}
	parts = append(parts, arg...)

	command := strings.Join(parts, " ")

	cmd = exec.Command("sh", "-c", command)

	return
}
