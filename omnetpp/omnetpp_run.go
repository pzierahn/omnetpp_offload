package omnetpp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
)

// Run the simulation with configuration (-c) and run number (-r)
func (project *OmnetProject) Run(ctx context.Context, config, run string) (err error) {
	sim, err := project.command(ctx, "-c", config, "-r", run)

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
