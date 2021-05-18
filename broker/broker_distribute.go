package broker

import pb "github.com/patrickz98/project.go.omnetpp/proto"

func (server *broker) distribute() {
	logger.Printf("distribute work!")

	for id, providerState := range server.providers.provider {
		arch := osArchId(providerState.Arch)
		logger.Printf("%s arch=%s usage=%3.0f%%", id, arch, providerState.CpuUsage)

		if providerState.CpuUsage >= 50.0 {
			//
			// Provider busy
			//

			continue
		}

		compile := server.simulations.pullCompile(providerState.Arch)

		if compile != nil {

			logger.Printf("--> compile='%s'", compile)

			server.providers.work[id] <- &pb.Work{
				Work: &pb.Work_Compile{Compile: compile},
			}

			continue
		}

		task := server.simulations.pullWork()

		if task == nil {
			//
			// No work
			//

			break
		}

		logger.Printf("--> task='%v'", task)

		server.providers.work[id] <- &pb.Work{
			Work: &pb.Work_Task{Task: task},
		}
	}
}
