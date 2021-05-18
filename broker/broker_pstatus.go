package broker

import (
	"fmt"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"github.com/patrickz98/project.go.omnetpp/worker"
	"google.golang.org/grpc/metadata"
)

func (server *broker) WorkSubscription(stream pb.Broker_WorkSubscriptionServer) (err error) {

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		logger.Printf("metadata missing")
		err = fmt.Errorf("metadata missing")
		return
	}

	var workerInfo worker.DeviceInfo
	workerInfo.UnMarshallMeta(md)

	workerId := providerId(workerInfo.WorkerId)
	if workerId == "" {
		err = fmt.Errorf("missing workerId in metadata")
		return
	}

	workStream := make(chan *pb.Work)
	defer close(workStream)

	server.providers.newProvider(workerId, workStream)

	go func() {
		for work := range workStream {

			logger.Printf("sending %s to %s", work, workerId)

			err = stream.Send(work)
			if err != nil {
				logger.Printf("error sending work: %v", err)
				break
			}
		}
	}()

	var state *pb.ProviderState

	for {
		state, err = stream.Recv()
		if err != nil {
			break
		}

		// jsonBytes, _ := json.MarshalIndent(state, "", "  ")
		// logger.Printf("UProviderLoad: stream=%s", string(jsonBytes))

		server.providers.update(state)
	}

	server.providers.remove(state)

	return
}

func (server *broker) ProviderLoad(pId *pb.ProviderId, stream pb.Broker_ProviderLoadServer) (err error) {
	logger.Printf("ProviderLoad: providerId.Id='%s'", pId.Id)

	//id := providerId(pId.Id)
	//
	//ch := make(chan *pb.ProviderState)
	//defer close(ch)
	//
	//server.providers.newListener(id, ch)
	//defer server.providers.removeListener(id, ch)
	//
	//for event := range ch {
	//
	//	logger.Printf("ProviderLoad: send event %f", event.CpuUsage)
	//
	//	err = stream.Send(event)
	//	if err != nil {
	//		break
	//	}
	//}

	return
}
