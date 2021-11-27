package provider

import (
	"context"
	"github.com/pzierahn/omnetpp_offload/eval"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (prov *provider) StartEvaluation(_ context.Context, scenario *pb.EvaluationScenario) (empty *emptypb.Empty, err error) {

	log.Printf("StartEvaluation: scenario=%v", simple.PrettyString(scenario))

	eval.Enable(scenario.ScenarioId, scenario.TrailId)
	eval.LogDevice(scenario.Connection, prov.numJobs)

	return &emptypb.Empty{}, nil
}

func (prov *provider) StopEvaluation(_ context.Context, _ *emptypb.Empty) (empty *emptypb.Empty, err error) {

	log.Printf("StopEvaluation:")
	eval.Disable()

	return &emptypb.Empty{}, nil
}
