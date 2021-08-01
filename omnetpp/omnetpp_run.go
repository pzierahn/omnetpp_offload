package omnetpp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// command
//
// Get simulation executable. This can ether be a simulationExe
// or simulationLib in conjunction with opp_run
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

// command
//
// Get simulation executable. This can ether be a simulationExe
// or simulationLib in conjunction with opp_run
func (project *OmnetProject) command(args ...string) (cmd *exec.Cmd, err error) {
	return project.commandContext(context.Background(), args...)
}

// RunContext the simulation with configuration (-c) and run number (-r)
func (project *OmnetProject) RunContext(ctx context.Context, config, run string) (err error) {
	sim, err := project.commandContext(ctx, "-c", config, "-r", run)

	if err != nil {
		return
	}

	// Debug
	//sim.Stdout = os.Stdout
	//sim.Stderr = os.Stderr

	var errBuf bytes.Buffer
	sim.Stderr = &errBuf

	pipe, err := sim.StdoutPipe()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(pipe)

	go func() {
		regex := regexp.MustCompile(`\(([0-9]{1,3})% total\)`)

		for scanner.Scan() {
			match := regex.FindStringSubmatch(scanner.Text())

			if len(match) == 2 {
				logger.Printf("base=%s config=%s run=%-3s (%s%%)\n",
					filepath.Base(project.Path), config, run, match[1])
			}
		}
	}()

	err = sim.Run()
	if err != nil {
		err = fmt.Errorf("err='%v' "+
			"stderr='%s' "+
			"command='%v' "+
			"dir='%v'\n", err, errBuf.String(), sim.Args, sim.Dir)
	}

	return
}

// Run the simulation with configuration (-c) and run number (-r)
func (project *OmnetProject) Run(config, run string) (err error) {
	return project.RunContext(context.Background(), config, run)
}
