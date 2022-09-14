package eval

import (
	"fmt"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"math/rand"
	"time"
)

type Event struct {
	Activity      string
	SimulationRun *pb.SimulationRun
	Filename      string
}

func (event Event) runId() (conf string, num string) {
	if event.SimulationRun == nil {
		return "", ""
	}

	return event.SimulationRun.Config, event.SimulationRun.RunNum
}

func time2Tex(date time.Time) (ts string) {
	data, _ := date.MarshalText()
	return string(data)
}

func tex2Time(ts string) (parse time.Time) {
	parse, _ = time.Parse(time.Layout, ts)
	return
}

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
		Timestamp: time2Tex(start),
		State:     StateStarted,
		EventId:   id,
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
			Timestamp: time2Tex(end),
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
