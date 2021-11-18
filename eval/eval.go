package eval

import (
	"context"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
)

const (
	_ = uint32(iota)
	StateStarted
	StateFinished
	StateFailed
)

const (
	ActivityCompile  = "Compile"
	ActivityRun      = "Run"
	ActivityUpload   = "Upload"
	ActivityDownload = "Download"
	ActivityCompress = "Compress"
	ActivityExtract  = "Extract"
)

var DeviceId string

var cli pb.EvaluationClient

func Init(conn *grpc.ClientConn) {
	cli = pb.NewEvaluationClient(conn)
}

func Close() {
	_, _ = cli.Finish(context.Background(), &emptypb.Empty{})
}

func SetScenario(simulationId string) {
	_, _ = cli.Init(context.Background(), &pb.Scenario{
		ScenarioId:   os.Getenv("SCENARIOID"),
		TrailId:      os.Getenv("TRAILID"),
		SimulationId: simulationId,
	})
}
