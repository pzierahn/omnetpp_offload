package provider

import (
	"context"
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var setupSync sync.Mutex

var configs = make(map[string]*pb.OppConfig)

//var sources = make(map[string]*pb.StorageRef)
//var binaries = make(map[string]*pb.StorageRef)

func (client *workerConnection) setup(job *pb.SimulationRun) (project omnetpp.OmnetProject, err error) {

	// Prevent that a simulation will be downloaded multiple times
	setupSync.Lock()
	defer setupSync.Unlock()

	// Simulation directory with simulation source code
	simulationBase := filepath.Join(cachePath, job.SimulationId)

	// This will be the working directory, that contains the results for the job
	// A symbolic copy is created to use all configs, ned files and ini files
	simulationPath := filepath.Join(
		cachePath,
		"mirrors",
		simple.NamedId(job.SimulationId, 8),
	)

	if _, err = os.Stat(simulationBase); err == nil {

		//
		// Simulation already downloaded and prepared
		//

		logger.Printf("simulation %s already downloaded\n", job.SimulationId)

		cachedOppConf := configs[job.SimulationId]

		err = simple.SymbolicCopy(simulationBase, simulationPath, cachedOppConf.ResultsPath)
		if err != nil {
			return
		}

		oppConf := omnetpp.Config{
			OppConfig: cachedOppConf,
			Path:      simulationPath,
		}

		project = omnetpp.New(&oppConf)

		return
	}

	//
	// Download simulation source
	//

	logger.Printf("get source for %s\n", job.SimulationId)
	src, err := client.broker.GetSource(
		context.Background(),
		&pb.SimulationId{Id: job.SimulationId})
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Printf("download %s to %s\n", job.SimulationId, simulationBase)
	byt, err := client.storage.Download(src.Source)
	if err != nil {
		return
	}

	logger.Printf("unzip %s\n", job.SimulationId)
	err = simple.UnTarGz(cachePath, &byt)
	if err != nil {
		_ = os.RemoveAll(simulationBase)
		return
	}

	//
	// Download simulation binary
	//

	// TODO: Check for local binary before downloading

	arch := &pb.Arch{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	binary, err := client.broker.GetBinary(
		context.Background(),
		&pb.SimulationBinaryRequest{
			SimulationId: job.SimulationId,
			Arch:         arch,
		})
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Printf("download %s to %s\n", job.SimulationId, simulationBase)
	bByt, err := client.storage.Download(binary.Binary)
	if err != nil {
		return
	}

	logger.Printf("unzip %s binary (%v)\n", job.SimulationId, arch)
	err = simple.UnTarGz(cachePath, &bByt)
	if err != nil {
		_ = os.RemoveAll(simulationBase)
		return
	}

	//
	// Setup
	//

	oppConfig, err := client.broker.GetOppConfig(
		context.Background(),
		&pb.SimulationId{Id: job.SimulationId})
	if err != nil {
		logger.Fatalln(err)
	}

	configs[job.SimulationId] = oppConfig

	logger.Printf("running setup %s\n", job.SimulationId)
	oppConf := omnetpp.Config{
		OppConfig: oppConfig,
		Path:      simulationBase,
	}

	// Create a new symbolic copy to get
	// results for each individual simulation run
	err = simple.SymbolicCopy(simulationBase, simulationPath, oppConfig.ResultsPath)
	if err != nil {
		return
	}

	oppConf.Path = simulationPath

	project = omnetpp.New(&oppConf)

	return
}

func (client *workerConnection) uploadResults(project omnetpp.OmnetProject, job *pb.SimulationRun) (err error) {

	//buf, err := project.ZipResults()
	//if err != nil {
	//	return
	//}
	//
	//ref, err := client.storage.Upload(&buf, storage.FileMeta{
	//	Bucket:   job.SimulationId,
	//	Filename: fmt.Sprintf("results_%s_%s.tar.gz", job.Config, job.RunNumber),
	//})
	//if err != nil {
	//	return
	//}
	//
	//results := pb.TaskResult{
	//	Task:    job,
	//	Results: ref,
	//}
	//
	//aff, err := client.broker.PutResults(context.Background(), &results)
	//if err != nil {
	//	_, _ = client.storage.Delete(ref)
	//	return
	//}
	//
	//if aff.Error != "" {
	//	err = fmt.Errorf(aff.Error)
	//}

	return
}

func (client *workerConnection) runTasks(job *pb.SimulationRun) {

	//
	// Setup simulation environment
	// Includes downloading and compiling the simulation
	//

	opp, err := client.setup(job)
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

	//err = client.uploadResults(opp, job)
	//if err != nil {
	//	logger.Fatalln(err)
	//}

	// Todo: Cleanup simulationBase

	// Cleanup symbolic mirrors
	_ = os.RemoveAll(opp.Path)
}
