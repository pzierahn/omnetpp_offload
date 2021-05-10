package stateinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc"
)

func Workers(config gconfig.GRPCConnection, simulationIds []string) {

	conn, err := grpc.Dial(config.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	status, err := client.WorkerInfo(context.Background(), &pb.WorkerInfoRequest{})
	if err != nil {
		panic(err)
	}

	jbyt, _ := json.MarshalIndent(status, "", "  ")
	fmt.Println(string(jbyt))

	return
}
