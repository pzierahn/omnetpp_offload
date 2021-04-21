package worker

import (
	pb "com.github.patrickz98.omnet/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"runtime"
	"sync"
)

type workerClient struct {
	sync.Mutex
	config        Config
	link          pb.Broker_LinkClient
	freeResources int
}

func (client *workerClient) OccupyResource(num int) (err error) {
	client.Lock()
	client.freeResources -= num
	logger.Println("OccupyResource", client.freeResources)
	client.Unlock()

	//err = client.SendClientInfo()

	return
}

func (client *workerClient) FeeResource() (err error) {
	client.Lock()
	client.freeResources++
	logger.Println("FeeResource", client.freeResources)
	client.Unlock()

	err = client.SendClientInfo()

	return
}

func (client *workerClient) SendClientInfo() (err error) {
	client.Lock()

	logger.Printf("sending info freeResources=%d\n", client.freeResources)

	info := pb.ClientInfo{
		WorkerId:      client.config.WorkerId,
		Os:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		NumCPU:        int32(runtime.NumCPU()),
		Timestamp:     timestamppb.Now(),
		FreeResources: int32(client.freeResources),
	}

	err = client.link.Send(&info)

	client.Unlock()

	return
}
