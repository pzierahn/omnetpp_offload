package omnetpp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/shell"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// command
//
// Get simulation executable. This can ether be a simulationExe
// or simulationLib in conjunction with opp_run
func (project *OmnetProject) command(args ...string) (cmd *exec.Cmd, err error) {

	iniDir, iniFile := filepath.Split(project.IniFile)
	iniDir = filepath.Join(project.Path, iniDir)

	nedPaths := make([]string, len(project.NedPaths))

	for inx, nedpath := range project.NedPaths {
		nedPaths[inx], err = filepath.Rel(iniDir, filepath.Join(project.Path, nedpath))
		if err != nil {
			return
		}
	}

	args = append([]string{
		"-u", "Cmdenv",
		// Ini filepath
		"-f", iniFile,
		// Ned paths
		"-n", strings.Join(nedPaths, ":"),
	}, args...)

	if project.UseLib {

		//
		// Use simulation shared library
		//

		var lib string
		lib, err = filepath.Rel(iniDir, filepath.Join(project.Path, project.Simulation))
		if err != nil {
			return
		}

		args = append([]string{
			"-l", lib,
		}, args...)

		cmd = shell.Command("opp_run", args...)
		cmd.Dir = iniDir
	} else {

		//
		// Use simulation exe
		//

		exe := filepath.Join(project.Path, project.Simulation)
		exe, err = filepath.Abs(exe)
		if err != nil {
			return
		}

		cmd = exec.Command(exe, args...)
		cmd.Dir = iniDir
	}

	return
}

// Run the simulation with configuration (-c) and run number (-r)
func (project *OmnetProject) Run(config, run string) (err error) {

	// Todo: Add timeout, because some simulations are running indefinitely
	sim, err := project.command("-c", config, "-r", run)

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

	go func() {
		regex := regexp.MustCompile(`\(([0-9]{1,3})% total\)`)
		scanner := bufio.NewScanner(pipe)

		for scanner.Scan() {
			match := regex.FindStringSubmatch(scanner.Text())

			if len(match) == 2 {
				logger.Printf("config=%s run=%s (%s%%)\n", config, run, match[1])
			}
		}
	}()

	err = sim.Run()
	if err != nil {
		err = fmt.Errorf("err='%v'"+
			"stderr='%s' "+
			"command='%v' "+
			"dir='%v'\n", err, errBuf.String(), sim.Args, sim.Dir)
	}

	return
}
