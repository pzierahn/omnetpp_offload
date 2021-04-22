package worker

import (
	"com.github.patrickz98.omnet/defines"
	"com.github.patrickz98.omnet/omnetpp"
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/simple"
	"com.github.patrickz98.omnet/storage"
	"os"
	"sync"
)

var setupSync sync.Mutex

func setup(work *pb.Work) (project omnetpp.OmnetProject, err error) {

	setupSync.Lock()
	defer setupSync.Unlock()

	simPath := defines.Simulation + "/" + work.SimulationId

	if _, err = os.Stat(simPath); err == nil {
		//
		// Simulation already downloaded and prepared
		//

		logger.Printf("simulation %s already downloaded\n", work.SimulationId)
		project = omnetpp.New(simPath)

		return
	}

	logger.Printf("download %s to %s\n", work.SimulationId, simPath)

	byt, err := storage.Download(work.Source)
	if err != nil {
		return
	}

	err = simple.UnTarGz(defines.Simulation, byt)
	if err != nil {
		_ = os.RemoveAll(simPath)
		return
	}

	logger.Printf("setup %s\n", work.SimulationId)

	project = omnetpp.New(simPath)
	err = project.Setup()

	return
}

func runTasks(client *workerClient, tasks *pb.Tasks) {
	for _, job := range tasks.Jobs {
		go func(job *pb.Work) {
			opp, err := setup(job)
			if err != nil {
				logger.Fatalln(err)
			}

			err = opp.Run(job.Config, job.RunNumber)
			if err != nil {
				logger.Fatalln(err)
			}

			err = client.FeeResource()
			if err != nil {
				logger.Fatalln(err)
			}
		}(job)
	}
}
