package worker

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/defines"
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/patrickz98/project.go.omnetpp/storage"
	"os"
	"path/filepath"
	"sync"
)

var setupSync sync.Mutex

var copyIgnores = map[string]bool{
	// Don't copy results
	"results/": true,
}

func setup(job *pb.Task) (project omnetpp.OmnetProject, err error) {

	// Prevent that a simulation will be downloaded multiple times
	setupSync.Lock()
	defer setupSync.Unlock()

	// Simulation directory with simulation source code
	simulationBase := filepath.Join(defines.SimulationPath, job.SimulationId)

	// This will be the working directory, that contains the results for the job
	// A symbolic copy is created to use all configs, ned files and ini files
	simulationPath := filepath.Join(
		defines.SimulationPath,
		"mirrors",
		simple.NamedId(job.SimulationId, 8),
	)

	if _, err = os.Stat(simulationBase); err == nil {

		//
		// Simulation already downloaded and prepared
		//

		logger.Printf("simulation %s already downloaded\n", job.SimulationId)

		err = simple.SymbolicCopy(simulationBase, simulationPath, copyIgnores)
		if err != nil {
			return
		}

		oppConf := omnetpp.Config{
			OppConfig: job.Simulation,
			Path:      simulationPath,
		}

		project = omnetpp.New(&oppConf)

		return
	}

	//
	// Download and compile the simulation
	//

	logger.Printf("download %s to %s\n", job.SimulationId, simulationBase)

	byt, err := storage.Download(job.Source)
	if err != nil {
		return
	}

	logger.Printf("unzip %s\n", job.SimulationId)

	err = simple.UnTarGz(defines.SimulationPath, byt)
	if err != nil {
		_ = os.RemoveAll(simulationBase)
		return
	}

	logger.Printf("running setup %s\n", job.SimulationId)

	oppConf := omnetpp.Config{
		OppConfig: job.Simulation,
		Path:      simulationBase,
	}

	// Compile simulation source code
	srcProject := omnetpp.New(&oppConf)
	err = srcProject.Setup()
	if err != nil {
		return
	}

	// Create a new symbolic copy
	err = simple.SymbolicCopy(simulationBase, simulationPath, copyIgnores)
	if err != nil {
		return
	}

	oppConf.Path = simulationPath

	project = omnetpp.New(&oppConf)

	return
}

func (client *workerConnection) uploadResults(project omnetpp.OmnetProject, job *pb.Task) (err error) {

	buf, err := project.ZipResults()
	if err != nil {
		return
	}

	ref, err := storage.Upload(&buf, storage.FileMeta{
		Bucket:   job.SimulationId,
		Filename: fmt.Sprintf("results_%s_%s.tar.gz", job.Config, job.RunNumber),
	})
	if err != nil {
		return
	}

	results := pb.TaskResult{
		Task:    job,
		Results: ref,
	}

	aff, err := client.client.CommitResults(context.Background(), &results)
	if err != nil {
		// TODO: Delete storage upload
		// _ = storage.Delete(ref)
		return
	}

	if aff.Error != "" {
		err = fmt.Errorf(aff.Error)
	}

	return
}

func (client *workerConnection) runTasks(job *pb.Task) {

	//
	// Setup simulation environment
	// Includes downloading and compiling the simulation
	//

	opp, err := setup(job)
	if err != nil {
		logger.Fatalln(err)
	}

	//
	// Setup simulation environment
	//

	err = opp.Run(job.Config, job.RunNumber)
	if err != nil {
		logger.Fatalln(err)
	}

	//
	// Upload simulation results
	//

	err = client.uploadResults(opp, job)
	if err != nil {
		logger.Fatalln(err)
	}

	// Todo: Cleanup simulationBase

	// Cleanup symbolic mirrors
	_ = os.RemoveAll(opp.SourcePath)
}
