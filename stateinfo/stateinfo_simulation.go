package stateinfo

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc"
	"time"
)

func Status(config gconfig.GRPCConnection, simulationIds []string) {

	conn, err := grpc.Dial(config.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	status, err := client.SimulationStatus(context.Background(), &pb.StatusRequest{
		SimulationIds: simulationIds,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Now())

	for _, item := range status.Items {
		fmt.Println(item.SimulationId)
		fmt.Printf("  Queue:    %5d\n", len(item.Queue))
		fmt.Printf("  Finished: %5d\n", len(item.Finished))

		fmt.Printf("  Assigned: %5d\n\n", len(item.Assigned))
		for _, elem := range item.Assigned {
			fmt.Printf("    %s %-3s %s\n", elem.Config.Config, elem.Config.RunNum, elem.WorkerId)
		}
	}

	return
}
