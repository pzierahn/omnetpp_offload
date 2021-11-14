package eval

import (
	"context"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/grpc"
	"os"
)

const (
	_ = iota
	StepStart
	StepSuccess
	StepError
)

var DeviceId string

var cli pb.EvalClient

func Init(conn *grpc.ClientConn) {
	cli = pb.NewEvalClient(conn)
}

func SetScenario(simulationId string) {
	_, _ = cli.Scenario(context.Background(), &pb.EvalScenario{
		ScenarioId:   os.Getenv("SCENARIOID"),
		TrailId:      os.Getenv("TRAILID"),
		SimulationId: simulationId,
	})
}
