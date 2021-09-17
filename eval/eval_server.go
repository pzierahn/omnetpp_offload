package eval

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/defines"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	fileActions = iota
	fileRuns
	fileTransfers
	fileSetups
)

//const (
//	fileActions   = "opp-edge-eval-actions.csv"
//	fileRuns      = "opp-edge-eval-runs.csv"
//	fileTransfers = "opp-edge-eval-setup.csv"
//	fileSetups    = "opp-edge-eval-transfers.csv"
//)

type Server struct {
	pb.UnimplementedEvalServer
	scenario *pb.EvalScenario
	files    map[int]*os.File
	sync     map[int]*sync.Mutex
}

func (server *Server) log(file int, msg proto.Message) {

	if server.scenario.ScenarioId == "" || server.scenario.TrailId == "" {
		return
	}

	server.sync[file].Lock()
	defer server.sync[file].Unlock()

	log.Printf("log: %s", simple.PrettyString(msg))
	_, values := MarshallProto(msg.ProtoReflect())

	writer := csv.NewWriter(server.files[file])
	defer writer.Flush()
	_ = writer.Write(values)
}

func (server *Server) Scenario(_ context.Context, scenario *pb.EvalScenario) (*emptypb.Empty, error) {

	log.Printf("scenario: %s", simple.PrettyString(scenario))

	server.scenario = scenario

	if scenario.ScenarioId == "" {
		return &emptypb.Empty{}, nil
	}

	dir := filepath.Join(defines.CacheDir(), "evaluation")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	protoTypes := map[int]proto.Message{
		fileActions:   &pb.ActionEvent{},
		fileRuns:      &pb.RunEvent{},
		fileTransfers: &pb.TransferEvent{},
		fileSetups:    &pb.SetupEvent{},
	}

	server.files = make(map[int]*os.File)
	server.sync = make(map[int]*sync.Mutex)

	id := fmt.Sprintf("s%s-t%s", scenario.ScenarioId, scenario.TrailId)

	for val, typ := range protoTypes {
		var name string

		switch val {
		case fileActions:
			name = "actions-" + id + ".csv"
		case fileRuns:
			name = "runs-" + id + ".csv"
		case fileTransfers:
			name = "transfers-" + id + ".csv"
		case fileSetups:
			name = "setup-" + id + ".csv"
		}

		filename := filepath.Join(dir, name)
		//file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}

		stat, err := file.Stat()
		if stat.Size() == 0 {
			writer := csv.NewWriter(file)
			headers, _ := MarshallProto(typ.ProtoReflect())
			if err = writer.Write(headers); err != nil {
				panic(err)
			}

			writer.Flush()
		}

		server.files[val] = file
		server.sync[val] = &sync.Mutex{}
	}

	return &emptypb.Empty{}, nil
}

func (server *Server) Action(_ context.Context, event *pb.ActionEvent) (*emptypb.Empty, error) {
	event.ScenarioId = server.scenario.ScenarioId
	event.TrailId = server.scenario.TrailId
	event.SimulationId = server.scenario.SimulationId

	server.log(fileActions, event)
	return &emptypb.Empty{}, nil
}

func (server *Server) Run(_ context.Context, event *pb.RunEvent) (*emptypb.Empty, error) {
	event.ScenarioId = server.scenario.ScenarioId
	event.TrailId = server.scenario.TrailId
	event.SimulationId = server.scenario.SimulationId

	server.log(fileRuns, event)
	return &emptypb.Empty{}, nil
}

func (server *Server) Transfer(_ context.Context, event *pb.TransferEvent) (*emptypb.Empty, error) {
	event.ScenarioId = server.scenario.ScenarioId
	event.TrailId = server.scenario.TrailId
	event.SimulationId = server.scenario.SimulationId

	server.log(fileTransfers, event)
	return &emptypb.Empty{}, nil
}

func (server *Server) Setup(_ context.Context, event *pb.SetupEvent) (*emptypb.Empty, error) {
	event.ScenarioId = server.scenario.ScenarioId
	event.TrailId = server.scenario.TrailId
	event.SimulationId = server.scenario.SimulationId

	server.log(fileSetups, event)
	return &emptypb.Empty{}, nil
}
