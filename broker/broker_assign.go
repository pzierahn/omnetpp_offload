package broker

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	prov "github.com/patrickz98/project.go.omnetpp/provider"
	"google.golang.org/grpc/metadata"
)

func (server *broker) Assignments(stream pb.Broker_AssignmentsServer) (err error) {

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		logger.Printf("metadata missing")
		err = fmt.Errorf("metadata missing")
		return
	}

	var meta prov.Meta
	meta.UnMarshallMeta(md)

	var node *provider
	node, err = newProvider(meta)
	if err != nil {
		return
	}
	defer func() {

		logger.Printf("%s: reassign %d tasks", node.id, len(node.assignments))

		for _, assignment := range node.assignments {

			switch left := assignment.Do.(type) {
			case *pb.Assignment_Build:
				//
				// TODO: Reassign compiling
				//

			case *pb.Assignment_Run:
				//
				// Reassign simulation runs
				//

				sState := server.simulations.getSimulationState(left.Run.SimulationId)
				sState.write(func() {
					id := tId(left.Run)
					sState.queue[id] = true
					sState.runs[id] = left.Run
				})
			}
		}

		node.close()
	}()

	logger.Printf("connected %s", node.id)

	server.providers.add(node)
	defer server.providers.remove(node)

	go func() {
		for assignment := range node.assign {

			logger.Printf("%s assigned '%v'", node.id, assignment)

			err = stream.Send(assignment)
			if err != nil {
				logger.Printf("error sending assignment: %v", err)
				break
			}
		}
	}()

	var utilization *pb.Utilization

	for {
		utilization, err = stream.Recv()
		if err != nil {
			break
		}

		node.setUtilization(utilization)
	}

	logger.Printf("disconnect %s", node.id)

	return
}
