package simple

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
)

func ShellCommand(name string, arg ...string) (cmd *exec.Cmd) {

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

func ShellCommandContext(ctx context.Context, name string, arg ...string) (cmd *exec.Cmd) {

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, name, arg...)
		return
	}

	parts := []string{
		name,
	}
	parts = append(parts, arg...)

	command := strings.Join(parts, " ")

	cmd = exec.CommandContext(ctx, "sh", "-c", command)

	return
}
