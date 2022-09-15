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

func (server *Server) ClockSync(_ context.Context, in *pb.Clock) (out *pb.Clock, _ error) {
	out = &pb.Clock{
		Timesent:     in.Timesent,
		Timereceived: timestamppb.New(time.Now()),
	}

	return
}

func (server *Server) Logs(_ *emptypb.Empty, stream pb.Evaluation_LogsServer) (_ error) {

	log.Printf("Starting evaluation")

	enabled.Store(true)
	defer enabled.Store(false)

	for {
		select {
		case event := <-eventChannel:
			err := stream.Send(event)
			if err != nil {
				log.Fatalf("Logs: couldn't send log event to consumer: %v", err)
			}
		case <-stream.Context().Done():
			log.Printf("Stopping evaluation")
			return
		}
	}

	return
}
