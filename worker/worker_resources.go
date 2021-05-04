package worker

import (
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

var resourceMutex sync.RWMutex

func (client *workerConnection) OccupyResource(num int) {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	client.freeResources -= num
	logger.Println("OccupyResource", client.freeResources)

	return
}

func (client *workerConnection) FeeResource() {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	client.freeResources++
	logger.Println("FeeResource", client.freeResources)

	return
}

func (client *workerConnection) SendResourceCapacity(link pb.Broker_TaskSubscriptionClient) (err error) {
	resourceMutex.RLock()
	defer resourceMutex.RUnlock()

	logger.Printf("sending info freeResources=%d\n", client.freeResources)

	info := pb.ResourceCapacity{
		WorkerId:      client.workerId,
		Timestamp:     timestamppb.Now(),
		FreeResources: int32(client.freeResources),
	}

	err = link.Send(&info)

	return
}
