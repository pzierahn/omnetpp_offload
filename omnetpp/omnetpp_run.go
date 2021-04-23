package omnetpp

import (
	"github.com/patrickz98/project.go.omnetpp/defines"
	"os"
	"os/exec"
	"path/filepath"
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

func (project *OmnetProject) RunLog(config, run string) (err error) {

	//
	// Run simulation
	//

	logDir := filepath.Join(defines.DataPath, "simulation-logs")
	_ = os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, config+"."+run+".log")

	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}

	sim := exec.Command("./"+project.simulationExe, "-c", config, "-r", run)
	sim.Dir = project.SourcePath
	sim.Stdout = file
	sim.Stderr = file

	err = sim.Run()

	return
}
