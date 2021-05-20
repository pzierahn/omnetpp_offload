package broker

func (server *broker) distribute() {
	// logger.Printf("distribute work!")

	server.providers.RLock()
	defer server.providers.RUnlock()

	for _, node := range server.providers.provider {
		// arch := osArchId(providerState.Arch)
		// logger.Printf("%s arch=%s usage=%3.0f%%", id, arch, providerState.CpuUsage)

		if node.busy() {

			//
			// Provider busy
			//

			continue
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
