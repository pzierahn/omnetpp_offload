package provider

import (
	"bytes"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/omnetpp"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"os"
	"path/filepath"
)

func (prov *provider) run(simulation *pb.Simulation) (ref *pb.StorageRef, err error) {

	//
	// Setup mirror simulation
	//

	// Simulation directory with simulation source code
	simulationBase := filepath.Join(cachePath, simulation.Id)

	// This will be the working directory, that contains the results for the job
	// A symbolic copy is created to use all configs, ned files and ini files
	simulationPath := filepath.Join(
		cachePath,
		"mirrors",
		simple.NamedId(simulation.Id, 8),
	)

	defer func() {
		_ = os.RemoveAll(simulationPath)
	}()

	err = simple.SymbolicCopy(simulationBase, simulationPath, simulation.OppConfig.ResultsPath)
	if err != nil {
		return
	}

	var results map[string]bool
	var filesBefor map[string]bool
	var filesAfter map[string]bool

	filesBefor, err = simple.ListDir(simulationPath)
	if err != nil {
		return
	}

	//
	// Run simulation run
	//

	oppConf := omnetpp.Config{
		OppConfig: simulation.OppConfig,
		Path:      simulationPath,
	}

	opp := omnetpp.New(&oppConf)
	err = opp.Run(simulation.Config, simulation.RunNum)
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
	results = simple.DirDiff(filesBefor, filesAfter)
	buf, err = simple.TarGzFiles(simulationPath, simulation.Id, results)
	if err != nil {
		return
	}

	ref = &pb.StorageRef{
		Bucket:   simulation.Id,
		Filename: fmt.Sprintf("results/%s_%s.tgz", simulation.Config, simulation.RunNum),
	}

	err = prov.store.Put(&buf, ref)

	return
}
