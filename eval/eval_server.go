package eval

import (
	"context"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/csv"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Server struct {
	pb.UnimplementedEvaluationServer
	mu       sync.Mutex
	scenario *pb.Scenario
	events   *csv.Writer
	devices  *csv.Writer
}

func (server *Server) Start(_ context.Context, scenario *pb.Scenario) (empty *emptypb.Empty, err error) {
	log.Printf("Start: evaluation logging %v", simple.PrettyString(scenario))

	server.mu.Lock()
	defer server.mu.Unlock()

	server.scenario = scenario

	dir := filepath.Join(gconfig.CacheDir(), "evaluation", scenario.Scenario)
	_ = os.MkdirAll(dir, 0755)

	eventsLogFile := fmt.Sprintf("%s_%03s_events.csv", scenario.Scenario, scenario.Trail)
	server.events = csv.NewWriter(dir, eventsLogFile)

	devicesLogFile := fmt.Sprintf("%s_%03s_devices.csv", scenario.Scenario, scenario.Trail)
	server.devices = csv.NewWriter(dir, devicesLogFile)

	return &emptypb.Empty{}, nil
}

func (server *Server) Finish(_ context.Context, _ *emptypb.Empty) (empty *emptypb.Empty, err error) {

	log.Printf("Finish: finishing evaluation logging")

	server.mu.Lock()
	defer server.mu.Unlock()

	server.scenario = nil

	server.events.Close()
	server.events = nil

	server.devices.Close()
	server.devices = nil

	return &emptypb.Empty{}, nil
}

func (server *Server) Init(_ context.Context, device *pb.Device) (*emptypb.Empty, error) {

	device.Timereceived = time2Tex(time.Now())

	log.Printf("Init: device=%s", simple.PrettyString(device))

	server.devices.RecordProtos(server.scenario.ProtoReflect(), device.ProtoReflect())

	return &emptypb.Empty{}, nil
}

func (server *Server) Log(_ context.Context, event *pb.Event) (*emptypb.Empty, error) {

	server.events.RecordProtos(server.scenario.ProtoReflect(), event.ProtoReflect())

	return &emptypb.Empty{}, nil
}
