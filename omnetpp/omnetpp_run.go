package omnetpp

import (
	"os"
	"os/exec"
)

func (project *OmnetProject) Run(config, run string) (err error) {

	//
	// Run simulation
	//

	sim := exec.Command("./"+project.simulationExe, "-c", config, "-r", run)
	sim.Dir = project.SourcePath
	sim.Stdout = os.Stdout
	sim.Stderr = os.Stderr

	err = sim.Run()

	return
}
