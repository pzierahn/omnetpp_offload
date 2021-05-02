package distribute

import (
	"context"
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/gconfig"
	pb "github.com/patrickz98/project.go.omnetpp/proto"
	"google.golang.org/grpc"
	"sort"
)

func Status(config gconfig.GRPCConnection, simulationId string) (err error) {

	conn, err := grpc.Dial(config.DialAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err = fmt.Errorf("did not connect: %v", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewBrokerClient(conn)

	status, err := client.Status(context.Background(), &pb.ResultsRequest{SimulationId: simulationId})
	if err != nil {
		panic(err)
	}

	sort.Slice(status.Queue, func(i, j int) bool {
		return fmt.Sprintf("%s.%6s", status.Queue[i].Config, status.Queue[i].RunNum) <
			fmt.Sprintf("%s.%6s", status.Queue[j].Config, status.Queue[j].RunNum)
	})

	sort.Slice(status.Assigned, func(i, j int) bool {
		return fmt.Sprintf("%s.%6s", status.Assigned[i].Config.Config, status.Assigned[i].Config.RunNum) <
			fmt.Sprintf("%s.%6s", status.Assigned[j].Config.Config, status.Assigned[j].Config.RunNum)
	})

	sort.Slice(status.Finished, func(i, j int) bool {
		return fmt.Sprintf("%s.%6s", status.Finished[i].Config, status.Finished[i].RunNum) <
			fmt.Sprintf("%s.%6s", status.Finished[j].Config, status.Finished[j].RunNum)
	})

	fmt.Println("Queue", len(status.Queue))
	for _, elem := range status.Queue {
		fmt.Println("  ", elem.Config, elem.RunNum)
	}

	fmt.Println("\nAssigned", len(status.Assigned))
	for _, elem := range status.Assigned {
		fmt.Println("  ", elem.Config.Config, elem.Config.RunNum, elem.WorkerId)
	}

	fmt.Println("\nFinished", len(status.Finished))
	for _, elem := range status.Finished {
		fmt.Println("  ", elem.Config, elem.RunNum)
	}

	return
}
