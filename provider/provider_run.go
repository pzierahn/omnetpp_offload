package provider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"os"
	"path/filepath"
)

func (prov *provider) run(ctx context.Context, run *pb.SimulationRun) (ref *pb.StorageRef, err error) {

	//
	// Setup mirror simulation
	//

	// Simulation directory with simulation source code
	simulationBase := filepath.Join(cachePath, run.SimulationId)

	// This will be the working directory, that contains the results for the job
	// A symbolic copy is created to use all configs, ned files and ini files
	simulationPath := filepath.Join(
		cachePath,
		"mirrors",
		simple.NamedId(run.SimulationId, 8),
	)

	defer func() {
		_ = os.RemoveAll(simulationPath)
	}()

	err = simple.SymbolicCopy(simulationBase, simulationPath)
	if err != nil {
		return
	}

	var results map[string]bool
	var filesBefore map[string]string
	var filesAfter map[string]string

	filesBefore, err = simple.ListDir(simulationPath)
	if err != nil {
		return
	}

	//
	// Run simulation run
	//

	oppConf := omnetpp.Config{
		OppConfig: run.OppConfig,
		Path:      simulationPath,
	}

	opp := omnetpp.New(&oppConf)
	err = opp.RunContext(ctx, run.Config, run.RunNum)
	if err != nil {
		return
	}

	filesAfter, err = simple.ListDir(simulationPath)
	if err != nil {
		return
	}

	//
	// Collect and upload results
	//

	var buf bytes.Buffer
	results = simple.DirDiff(filesBefore, filesAfter)
	buf, err = simple.TarGzFiles(simulationPath, run.SimulationId, results)
	if err != nil {
		return
	}

	ref = &pb.StorageRef{
		Bucket:   run.SimulationId,
		Filename: fmt.Sprintf("results/%s_%s.tgz", run.Config, run.RunNum),
	}

	err = prov.store.Put(&buf, ref)

	return
}
