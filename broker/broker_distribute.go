package broker

import pb "github.com/patrickz98/project.go.omnetpp/proto"

func (server *broker) distribute() {
	//logger.Printf("distribute work!")

	server.providers.RLock()
	defer server.providers.RUnlock()

	for _, node := range server.providers.provider {

		// arch := osArchId(providerState.Arch)
		//logger.Printf("checking %s", node.id)

		if node.busy() {

			//
			// Provider busy
			//

			continue
		}

		node.RLock()
		build := server.simulations.pullCompile(node.arch)
		node.RUnlock()

		if build != nil {

			//logger.Printf("--> %s compile %s", node.id, build.simulationId)

			assignBuild := true

			for _, prov2 := range server.providers.provider {
				prov2.RLock()
				if prov2.building == build.simulationId {
					//
					// Build is already assigned to a provider
					//

					assignBuild = false
					prov2.RUnlock()

					break
				}
				prov2.RUnlock()
			}

			if assignBuild {
				node.assignCompile(&pb.Build{
					SimulationId: build.simulationId,
					Config:       build.oppConfig,
					Source:       build.source,
				})
				continue
			}
		}

		slots := node.freeSlots()
		for inx := 0; inx < slots; inx++ {
			task := server.simulations.pullWork(node.arch)
			if task == nil {
				// No jobs left to do for arch
				break
			}

			node.assignRun(task)
		}
	}
}
