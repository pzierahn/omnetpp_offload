package eval

import (
	"fmt"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"time"
)

func Log(event Event) (finish func(err error, dlsize uint64) (duration time.Duration)) {

	start := time.Now()

	if !enabled.Load() {
		return func(_ error, _ uint64) time.Duration {
			return time.Now().Sub(start)
		}
	}

	id := fmt.Sprintf("%00000000x", rand.Uint32())

	conf, runNum := event.runId()

	eventChannel <- &pb.Event{
		EventId:   id,
		Timestamp: timestamppb.New(start),
		State:     StateStarted,
		Activity:  event.Activity,
		Filename:  event.Filename,
		Config:    conf,
		RunNum:    runNum,
	}

	finish = func(err error, dlsize uint64) time.Duration {

		var end = time.Now()
		var state uint32
		var fail string

		if err != nil {
			state = StateFailed
			fail = err.Error()
		} else {
			state = StateFinished
		}

		eventChannel <- &pb.Event{
			EventId:   id,
			Timestamp: timestamppb.New(end),
			State:     state,
			Error:     fail,
			Activity:  event.Activity,
			Filename:  event.Filename,
			Config:    conf,
			RunNum:    runNum,
			ByteSize:  dlsize,
		}

		return end.Sub(start)
	}

	return
}
