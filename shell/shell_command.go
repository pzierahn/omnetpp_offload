package shell

import (
	"os/exec"
	"runtime"
	"strings"
)

func Command(name string, arg ...string) (cmd *exec.Cmd) {

	if runtime.GOOS == "windows" {
		cmd = exec.Command(name, arg...)
		return
	}

	parts := []string{
		name,
	}
	parts = append(parts, arg...)

	command := strings.Join(parts, " ")

	cmd = exec.Command("sh", "-c", command)

	return
}
