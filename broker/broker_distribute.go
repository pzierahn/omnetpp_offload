package broker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"sync"
)

var compileAssignmentsMu sync.Mutex
var compileAssignments = make(map[string]map[osArch]string) // simulationId --> osArch --> providerId

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

		// TODO: Check if simulation compilation is already assigned to an provider for an node.arch

		simulation := server.simulations.pullCompile(node.arch)
		if simulation != nil {

			compileAssignmentsMu.Lock()

			if compileAssignments[simulation.simulationId] == nil {
				compileAssignments[simulation.simulationId] = make(map[osArch]string)
			}

			tmpId, assigned := compileAssignments[simulation.simulationId][osArchId(node.arch)]

			logger.Printf("#### Check compile assignments: %s>%s>%s assigned=%v",
				simulation.simulationId, osArchId(node.arch), tmpId, assigned)

			if !assigned {

				compileAssignments[simulation.simulationId][osArchId(node.arch)] = node.id

				build := &pb.Build{
					SimulationId: simulation.simulationId,
					OppConfig:    simulation.oppConfig,
					Source:       simulation.source,
				}

				node.assignWork(&pb.Assignment{
					Do: &pb.Assignment_Build{Build: build},
				})

				compileAssignmentsMu.Unlock()
				continue
			}

			compileAssignmentsMu.Unlock()
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
