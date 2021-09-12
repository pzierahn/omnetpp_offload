package eval

import (
	"context"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"google.golang.org/grpc"
	"os"
)

const (
	_ = iota
	StepStart
	StepSuccess
	StepError
)

var client pb.EvalClient

func Init(conn *grpc.ClientConn) {
	client = pb.NewEvalClient(conn)
}

func SetScenario(simulationId string) {
	_, _ = client.Scenario(context.Background(), &pb.EvalScenario{
		ScenarioId:   os.Getenv("SCENARIOID"),
		TrailId:      os.Getenv("TRAILID"),
		SimulationId: simulationId,
	})
}
