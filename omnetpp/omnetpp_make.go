package omnetpp

import (
	"com.github.patrickz98.omnet/shell"
	"fmt"
	"os"
	"runtime"
)

func (project *OmnetProject) MakeMake() (err error) {

	//
	// Create Makefile
	//

	makemake := shell.Command("opp_makemake",
		"-f", "--deep", "-u", "Cmdenv", "-o", project.simulationExe)

	makemake.Dir = project.SourcePath
	makemake.Stdout = os.Stdout
	makemake.Stderr = os.Stderr

	err = makemake.Run()

	return
}

func (project *OmnetProject) Compile() (err error) {

	//
	// Compile simulation
	//

	makeCmd := shell.Command("make", "-j", fmt.Sprint(runtime.NumCPU()), "MODE=release")
	makeCmd.Dir = project.SourcePath
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr

	err = makeCmd.Run()

	return
}

func (project *OmnetProject) Clean() (err error) {

	//
	// Compile simulation
	//

	makeCmd := shell.Command("make", "cleanall")
	makeCmd.Dir = project.SourcePath
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr

	err = makeCmd.Run()

	return
}
