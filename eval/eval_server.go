package eval

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/simple"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Server struct {
	pb.UnimplementedEvaluationServer
	scenario *pb.Scenario
	logFile  *os.File
	writer   *csv.Writer
	mu       sync.Mutex
}

func (server *Server) Init(_ context.Context, scenario *pb.Scenario) (*emptypb.Empty, error) {

	server.mu.Lock()
	defer server.mu.Unlock()

	log.Printf("init: %s", simple.PrettyString(scenario))

	if scenario.ScenarioId == "" {
		return &emptypb.Empty{}, nil
	}

	server.scenario = scenario

	dir := filepath.Join(gconfig.CacheDir(), "evaluation", scenario.ScenarioId)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalln(err)
	}

	filename := fmt.Sprintf("%s_%03s", scenario.ScenarioId, scenario.TrailId)
	path := filepath.Join(dir, filename+".csv")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	server.logFile = file
	server.writer = csv.NewWriter(file)

	return &emptypb.Empty{}, nil
}

func (server *Server) Finish(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {

	server.mu.Lock()
	defer server.mu.Unlock()

	server.scenario = nil
	_ = server.logFile.Close()
	server.logFile = nil

	return &emptypb.Empty{}, nil
}

func (server *Server) Log(_ context.Context, event *pb.Event) (*emptypb.Empty, error) {

	if server.logFile == nil {
		return &emptypb.Empty{}, nil
	}

	server.mu.Lock()
	defer server.mu.Unlock()

	headers, record := MarshallProto(server.scenario.ProtoReflect())

	eh, er := MarshallProto(event.ProtoReflect())
	headers = append(headers, eh...)
	record = append(record, er...)

	stat, err := server.logFile.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if stat.Size() == 0 {
		if err = server.writer.Write(headers); err != nil {
			log.Fatalln(err)
		}
	}

	if err = server.writer.Write(record); err != nil {
		log.Fatalln(err)
	}

	server.writer.Flush()

	return &emptypb.Empty{}, nil
}
