package broker

import (
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/utils"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/metadata"
)

func (server *broker) Link(stream pb.Broker_LinkServer) (err error) {

	md, ok := metadata.FromIncomingContext(stream.Context())

	if !ok {
		logger.Printf("metadata missing")
		err = fmt.Errorf("metadata missing")
		return
	}

	var id string
	id, err = utils.MetaString(md, "workerId")
	if err != nil {
		return
	}

	logger.Println("link", id)

	//
	// Send work to clients
	//

	go func() {
		for {
			work, ok := <-server.work

			if !ok {
				logger.Println("exit work mode for", id)
			}

			logger.Println("send work to ", id)

			err := stream.Send(&pb.Work{
				SimulationId: work.SimulationId,
				ConfigId:     "config-xxx",
				Source:       work.Source,
				Config:       "Config-XXX",
				RunNumber:    "1",
			})

			if err != nil {
				logger.Println(err)
			}
		}
	}()

	//
	// Receive Client Info
	//

	for {
		var info *pb.ClientInfo
		info, err = stream.Recv()
		if err != nil {
			break
		}

		if id == "" {
			id = info.Id
		}

		server.workers[info.Id] = info

		jsonBytes, _ := json.MarshalIndent(info, "", "    ")
		logger.Println("link", string(jsonBytes))
	}

	logger.Println("unlink", id)

	return
}
