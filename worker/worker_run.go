package worker

import (
	"com.github.patrickz98.omnet/omnetpp"
	pb "com.github.patrickz98.omnet/proto"
)

func runTasks(client *workerClient, tasks *pb.Tasks) {
	done := make(chan bool, len(tasks.Jobs))
	defer close(done)

	opp := omnetpp.New("/Users/patrick/Desktop/tictoc")
	err := opp.SetupCheck()
	if err != nil {
		logger.Fatalln(err)
	}

	for _, job := range tasks.Jobs {
		go func(job *pb.Work) {
			err := opp.Run(job.Config, job.RunNumber)
			if err != nil {
				logger.Fatalln(err)
			}

			err = client.FeeResource()
			if err != nil {
				logger.Fatalln(err)
			}
		}(job)
	}

	for inx := 0; inx < len(tasks.Jobs); inx++ {
		<-done
	}
}