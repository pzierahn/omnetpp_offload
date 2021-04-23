package worker

import (
	pb "com.github.patrickz98.omnet/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

var resourceMutex sync.Mutex

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

func (client *workerConnection) SendResourceCapacity(link pb.Broker_LinkClient) (err error) {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	logger.Printf("sending info freeResources=%d\n", client.freeResources)

	info := pb.ResourceCapacity{
		WorkerId:      client.config.WorkerId,
		Timestamp:     timestamppb.Now(),
		FreeResources: int32(client.freeResources),
	}

	err = link.Send(&info)

	return
}
