package eval

import (
	"context"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sync/atomic"
	"time"
)

var enabled atomic.Bool
var eventChannel = make(chan *pb.Event, 64)

type Server struct {
	pb.UnimplementedEvaluationServer
}

func (server *Server) EnableLog(_ context.Context, state *pb.Enable) (*emptypb.Empty, error) {
	log.Printf("EnableLog %v", state.Enable)

	enabled.Store(state.Enable)
	return &emptypb.Empty{}, nil
}

func (server *Server) ClockSync(_ context.Context, in *pb.Clock) (out *pb.Clock, _ error) {
	out = &pb.Clock{
		Timesent:     in.Timesent,
		Timereceived: timestamppb.New(time.Now()),
	}

	return
}

func (server *Server) Logs(_ *emptypb.Empty, stream pb.Evaluation_LogsServer) (_ error) {

	if !enabled.Load() {
		return
	}

	for event := range eventChannel {
		err := stream.Send(event)
		if err != nil {
			log.Fatalf("Logs: couldn't send log event to consumer: %v", err)
		}
	}

	return
}
