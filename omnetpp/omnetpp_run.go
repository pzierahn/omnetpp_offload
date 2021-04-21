package omnetpp

import (
	"com.github.patrickz98.omnet/defines"
	"log"
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

func (project *OmnetProject) RunLog(config, run string) (err error) {

	//
	// Run simulation
	//

	logPath := defines.DataPath + "/simulation-logs"
	_ = os.MkdirAll(logPath, 0755)

	file, err := os.OpenFile(logPath+"/"+config+"."+run+".log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	sim := exec.Command("./"+project.simulationExe, "-c", config, "-r", run)
	sim.Dir = project.SourcePath
	sim.Stdout = file
	sim.Stderr = file

	err = sim.Run()

	return
}
