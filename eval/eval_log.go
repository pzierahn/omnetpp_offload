package eval

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"log"
	"math/rand"
	"os"
	"runtime"
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

func timestampNow() (ts string) {
	return time2Tex(time.Now())
}

func Log(event Event) (finish func(err error, dlsize uint64) (duration time.Duration)) {

	start := time.Now()

	if !enabled {
		return func(_ error, _ uint64) time.Duration {
			return time.Now().Sub(start)
		}
	}

	id := fmt.Sprintf("%00000000x", rand.Uint32())

	conf, runNum := event.runId()

	ctx := context.Background()
	_, _ = cli.Log(ctx, &pb.Event{
		Timestamp: time2Tex(start),
		DeviceId:  deviceId,
		State:     StateStarted,
		EventId:   id,
		Activity:  event.Activity,
		Filename:  event.Filename,
		Config:    conf,
		RunNum:    runNum,
	})

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

		_, _ = cli.Log(ctx, &pb.Event{
			EventId:   id,
			DeviceId:  deviceId,
			Timestamp: time2Tex(end),
			State:     state,
			Error:     fail,
			Activity:  event.Activity,
			Filename:  event.Filename,
			Config:    conf,
			RunNum:    runNum,
			ByteSize:  dlsize,
		})

		return end.Sub(start)
	}

	return
}

func LogDevice(connect string, numJobs int) {

	if !enabled {
		return
	}

	host, _ := os.Hostname()

	_, err := cli.Init(context.Background(), &pb.Device{
		DeviceId: deviceId,
		Hostname: host,
		Timesent: timestampNow(),
		Os:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		NumCPUs:  uint32(runtime.NumCPU()),
		NumJobs:  uint32(numJobs),
		Connect:  connect,
	})

	if err != nil {
		log.Fatalf("couldn't initialize evaluation service: %v", err)
	}
}
