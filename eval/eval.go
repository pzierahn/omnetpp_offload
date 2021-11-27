package eval

import (
	"context"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
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
	ActivityChecksum = "Checksum"
	ActivityExtract  = "Extract"
)

var enabled bool
var Scenario string
var Trail string

var cli pb.EvaluationClient
var deviceId string

func Init(broker *grpc.ClientConn, id string) {
	deviceId = id

	if deviceId == "" {
		deviceId, _ = os.Hostname()
	}

	cli = pb.NewEvaluationClient(broker)
}

func Enable(scenario, trail string) {

	if scenario == "" {
		log.Fatalf("error: scenario name is invalid")
	}

	Scenario, Trail = scenario, trail
	enabled = true

	log.Printf("start evaluation logging (scenario=%v, trail=%v)", Scenario, Trail)
}

func IsEnabled() bool {
	return enabled
}

func Disable() {

	log.Printf("stop evaluation logging (scenario=%v, trail=%v)", Scenario, Trail)

	enabled = false
	Scenario, Trail = "", ""
}

func Start(ctx context.Context, scenario, trail, simulation string) {

	Enable(scenario, trail)

	_, err := cli.Start(ctx, &pb.Scenario{
		Scenario:   scenario,
		Trail:      trail,
		Simulation: simulation,
	})
	if err != nil {
		log.Printf("couldn't start evaluation")
	}
}

func Finish() {
	defer Disable()

	if !enabled {
		return
	}

	_, _ = cli.Finish(context.Background(), &emptypb.Empty{})
}
