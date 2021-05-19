package broker

import pb "github.com/patrickz98/project.go.omnetpp/proto"

func (server *broker) distribute() {
	// logger.Printf("distribute work!")

	for id, node := range server.providers.provider {
		// arch := osArchId(providerState.Arch)
		// logger.Printf("%s arch=%s usage=%3.0f%%", id, arch, providerState.CpuUsage)

		if node.busy() {
			//
			// Provider busy
			//

			continue
		}

		logger.Printf("%s assignments=%d", id, len(node.assignments))

		compile := server.simulations.pullCompile(node.arch)
		if compile != nil {

			node.assignWork(&pb.Assignment{
				Do: &pb.Assignment_Build{Build: compile},
			})

			continue
		}

		task := server.simulations.pullWork(node.arch)
		if task != nil {
			node.assignWork(&pb.Assignment{
				Do: &pb.Assignment_Run{Run: task},
			})

			continue
		}
	}
}
