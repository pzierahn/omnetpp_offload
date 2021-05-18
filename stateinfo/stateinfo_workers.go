package stateinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc"
)

func Workers(config gconfig.GRPCConnection) {

	logger.Printf("connect to %v", config.DialAddr())

	conn, err := grpc.Dial(config.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	stream, err := client.ProviderLoad(context.Background(), &pb.ProviderId{Id: "patricks-mbp-ea11aab0"})
	if err != nil {
		logger.Fatalln(err)
	}

	for {
		pstate, err := stream.Recv()
		if err != nil {
			logger.Fatalln(err)
		}

		jsonBytes, _ := json.MarshalIndent(pstate, "", "  ")
		logger.Printf("%s", string(jsonBytes))
	}
}
