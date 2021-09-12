package eval

import (
	"context"
	"encoding/csv"
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
	fileActions   = "actions.csv"
	fileRuns      = "runs.csv"
	fileTransfers = "transfers.csv"
	fileSetups    = "setups.csv"
)

type Server struct {
	pb.UnimplementedEvalServer
	scenario *pb.EvalScenario
	files    map[string]*os.File
	sync     map[string]*sync.Mutex
}

func (server *Server) log(file string, msg proto.Message) {
	server.sync[file].Lock()
	defer server.sync[file].Unlock()

	log.Printf("log: %s", simple.PrettyString(msg))
	_, values := MarshallProto(msg.ProtoReflect())

	writer := csv.NewWriter(server.files[file])
	defer writer.Flush()
	_ = writer.Write(values)
}

func (server *Server) Scenario(_ context.Context, scenario *pb.EvalScenario) (*emptypb.Empty, error) {
	server.scenario = scenario

	log.Printf("scenario: %s", simple.PrettyString(scenario))

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

func NewServer() (server *Server) {
	dir := filepath.Join(defines.CacheDir(), "eval")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	protoTypes := map[string]proto.Message{
		fileActions:   &pb.ActionEvent{},
		fileRuns:      &pb.RunEvent{},
		fileTransfers: &pb.TransferEvent{},
		fileSetups:    &pb.SetupEvent{},
	}

	server = &Server{
		files: make(map[string]*os.File),
		sync:  make(map[string]*sync.Mutex),
	}

	for name, typ := range protoTypes {
		filename := filepath.Join(dir, name)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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

		server.files[name] = file
		server.sync[name] = &sync.Mutex{}
	}

	return
}
