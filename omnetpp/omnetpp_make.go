package omnetpp

import (
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/shell"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func (project *OmnetProject) MakeMake() (err error) {

	if project.BuildScript != "" {

		//
		// Buildscript provided: nothing to do here
		//

		return
	}

	//
	// Create Makefile
	//

	src, obj := filepath.Split(project.Simulation)

	logger.Printf("creating makefile in %s\n", src)

	makemake := shell.Command("opp_makemake",
		"-f", "--deep", "-u", "Cmdenv", "-o", obj)

	makemake.Dir = filepath.Join(project.Path, src)
	makemake.Stdout = os.Stdout
	makemake.Stderr = os.Stderr

	err = makemake.Run()

	return
}

func (project *OmnetProject) Compile() (err error) {

	if project.BuildScript != "" {

		//
		// Compile simulation using the buildscript
		//

		dir, script := filepath.Split(project.BuildScript)

		logger.Printf("running %s\n", project.BuildScript)

		build := exec.Command("sh", script)
		build.Dir = filepath.Join(project.Path, dir)

		//logger.Printf("############ build.Dir %s\n", build.Dir)
		//build.Stdout = os.Stdout
		//build.Stderr = os.Stderr

		err = build.Run()

		return
	}

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
	// Clean simulation
	//

	logger.Printf("cleaning %s\n", project.SourcePath)

	makeCmd := shell.Command("make", "cleanall")
	makeCmd.Dir = filepath.Join(project.Path, project.SourcePath)
	//makeCmd.Stdout = os.Stdout
	//makeCmd.Stderr = os.Stderr

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
