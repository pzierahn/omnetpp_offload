package omnetpp

import (
	"context"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"os/exec"
	"path/filepath"
	"strings"
)

// Returns simulation cmd with a context. This can ether be a simulationExe
// or simulationLib in conjunction with opp_run.
func (project *OmnetProject) commandContext(ctx context.Context, args ...string) (cmd *exec.Cmd, err error) {

	base := filepath.Join(project.Path, project.BasePath)

	args = append(args, "-u", "Cmdenv")

	for _, ini := range project.IniFiles {
		ini = filepath.Join(project.Path, ini)
		ini, err = filepath.Rel(base, ini)
		if err != nil {
			return
		}

		args = append(args, "-f", ini)
	}

	nedPaths := make([]string, len(project.NedPaths))

	for inx, nedpath := range project.NedPaths {
		nedpath = filepath.Join(project.Path, nedpath)
		nedPaths[inx], err = filepath.Rel(base, nedpath)
		if err != nil {
			return
		}
	}

	if len(nedPaths) > 0 {
		args = append(args, "-n", strings.Join(nedPaths, ":"))
	}

	if project.UseLib {

		//
		// Use simulation library
		//

		lib := filepath.Join(project.Path, project.Simulation)
		lib, err = filepath.Rel(base, lib)
		if err != nil {
			return
		}

		args = append(args, "-l", lib)

		cmd = simple.ShellCommandContext(ctx, "opp_run", args...)
		cmd.Dir = base
	} else {

		//
		// Use simulation exe
		//

		exe := filepath.Join(project.Path, project.Simulation)
		exe, err = filepath.Abs(exe)
		if err != nil {
			return
		}

		cmd = exec.CommandContext(ctx, exe, args...)
		cmd.Dir = base
	}

	return
}

// command returns simulation executable. This can ether be a simulationExe
// or simulationLib in conjunction with opp_run
func (project *OmnetProject) command(args ...string) (cmd *exec.Cmd, err error) {
	return project.commandContext(context.Background(), args...)
}
