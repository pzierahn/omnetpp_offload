package broker

import pb "github.com/patrickz98/project.go.omnetpp/proto"

func (server *broker) distribute() {
	// logger.Printf("distribute work!")

	for _, node := range server.providers.provider {
		// arch := osArchId(providerState.Arch)
		// logger.Printf("%s arch=%s usage=%3.0f%%", id, arch, providerState.CpuUsage)

		if node.busy() {
			//
			// Provider busy
			//

			continue
		}

		compile := server.simulations.pullCompile(node.arch)
		if compile != nil {

			node.assignWork(&pb.Assignment{
				Do: &pb.Assignment_Build{Build: compile},
			})

			continue
		}

		slots := node.freeSlots()
		for inx := 0; inx < slots; inx++ {
			task := server.simulations.pullWork(node.arch)
			if task == nil {
				// No jobs left
				break
			}

			node.assignWork(&pb.Assignment{
				Do: &pb.Assignment_Run{Run: task},
			})
		}
	}
}
