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

func (project *OmnetProject) Setup() (err error) {
	err = project.Clean()
	if err != nil {
		return
	}

	err = project.MakeMake()
	if err != nil {
		return
	}

	err = project.Compile()
	if err != nil {
		return
	}

	return
}

func (project *OmnetProject) SetupCheck() (err error) {

	if _, err = os.Stat(project.SourcePath + "/" + project.simulationExe); err == nil {
		return
	}

	err = project.Clean()
	if err != nil {
		return
	}

	err = project.MakeMake()
	if err != nil {
		return
	}

	err = project.Compile()
	if err != nil {
		return
	}

	return
}
