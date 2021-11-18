package eval

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"math/rand"
	"time"
)

type Event struct {
	DeviceId      string
	Activity      string
	SimulationRun *pb.SimulationRun
	Filename      string
}

func (event Event) runId() (id string) {
	if event.SimulationRun != nil {
		id = fmt.Sprintf("%s_%03s", event.SimulationRun.Config, event.SimulationRun.RunNum)
	}

	return
}

func timestampNow() (ts string) {
	data, _ := time.Now().MarshalText()
	return string(data)
}

func Log(event Event) (done func(err error, dlsize uint64)) {

	id := fmt.Sprintf("%00000000x", rand.Uint32())

	ctx := context.Background()
	_, _ = cli.Log(ctx, &pb.Event{
		Timestamp:     timestampNow(),
		DeviceId:      event.DeviceId,
		State:         StateStarted,
		EventId:       id,
		Activity:      event.Activity,
		Filename:      event.Filename,
		SimulationRun: event.runId(),
	})

	done = func(err error, dlsize uint64) {

		var state uint32
		var fail string

		if err != nil {
			state = StateFailed
			fail = err.Error()
		} else {
			state = StateFinished
		}

		_, _ = cli.Log(ctx, &pb.Event{
			EventId:       id,
			DeviceId:      event.DeviceId,
			Timestamp:     timestampNow(),
			State:         state,
			Error:         fail,
			Activity:      event.Activity,
			Filename:      event.Filename,
			SimulationRun: event.runId(),
			ByteSize:      dlsize,
		})

		return
	}

	return
}
