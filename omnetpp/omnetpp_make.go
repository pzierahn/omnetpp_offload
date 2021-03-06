package omnetpp

import (
	"context"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/simple"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// MakeMake creates a Makefile in the designated project by using opp_makemake.
func (project *OmnetProject) MakeMake(ctx context.Context) (err error) {

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

	args := []string{"-f", "--deep", "-u", "Cmdenv", "-o", obj}

	if project.UseLib {
		args = append(args, "--make-so")
	}

	makemake := simple.ShellCommandContext(ctx, "opp_makemake", args...)

	makemake.Dir = filepath.Join(project.Path, src)
	makemake.Stdout = os.Stdout
	makemake.Stderr = os.Stderr

	err = makemake.Run()

	return
}

// Compile compiles the simulation. The simulation can
// be built in two ways. If the config contains a build
// script it will be used to compile the simulation.
// Otherwise, a Makefile will be used to compile the simulation.
func (project *OmnetProject) Compile(ctx context.Context) (err error) {

	if project.BuildScript != "" {

		//
		// Compile simulation using the buildscript
		//

		dir, script := filepath.Split(project.BuildScript)

		logger.Printf("running %s\n", project.BuildScript)

		build := exec.CommandContext(ctx, "sh", script)
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

	src, _ := filepath.Split(project.Simulation)

	makeCmd := simple.ShellCommandContext(ctx, "make", "-j", fmt.Sprint(runtime.NumCPU()), "MODE=release")
	makeCmd.Dir = filepath.Join(project.Path, src)
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr

	err = makeCmd.Run()

	return
}

// Clean cleans the simulation source by calling make clean in the SourcePath.
func (project *OmnetProject) Clean(ctx context.Context) (err error) {

	//
	// Clean simulation
	//

	logger.Printf("cleaning %s\n", project.SourcePath)

	makeCmd := simple.ShellCommandContext(ctx, "make", "cleanall")
	makeCmd.Dir = filepath.Join(project.Path, project.SourcePath)
	//makeCmd.Stdout = os.Stdout
	//makeCmd.Stderr = os.Stderr

	err = makeCmd.Run()

	return
}

// Setup prepares the simulation for execution. It cleans, creates a makefile and compiles the simulation.
func (project *OmnetProject) Setup(ctx context.Context, clean bool) (err error) {

	if clean {
		err = project.Clean(ctx)
		if err != nil {
			return
		}
	}

	err = project.MakeMake(ctx)
	if err != nil {
		return
	}

	err = project.Compile(ctx)
	if err != nil {
		return
	}

	return
}
