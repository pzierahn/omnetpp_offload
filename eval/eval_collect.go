package eval

import (
	"context"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/csv"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"path/filepath"
	"time"
)

var collecting = false
var scenarioId string
var trailNum string

var events *csv.Writer
var devices *csv.Writer

func StartCollecting(scenario, trail string) {
	dir := filepath.Join(gconfig.CacheDir(), "evaluation", scenario)

	log.Printf("################################################################")
	log.Printf("Enable evaluation logging: scenario=%s trail=%s", scenario, trail)
	log.Printf("Logging to %s", dir)
	log.Printf("################################################################")

	eventsLogFile := fmt.Sprintf("%s_%03s_events.csv", scenario, trail)
	events = csv.NewWriter(dir, eventsLogFile)

	devicesLogFile := fmt.Sprintf("%s_%03s_devices.csv", scenario, trail)
	devices = csv.NewWriter(dir, devicesLogFile)

	collecting = true
	scenarioId = scenario
	trailNum = trail
}

func CollectLogs(prov *pb.ProviderInfo, client *grpc.ClientConn) {
	if !collecting {
		return
	}

	ctx := context.Background()
	evaluation := pb.NewEvaluationClient(client)

	_, err := evaluation.EnableLog(ctx, &pb.Enable{Enable: true})
	if err != nil {
		log.Fatalf("CollectLogs: cloudn't enable evaluation logging on provider %v: %v",
			prov.ProviderId, err)
	}
	defer func() {
		_, _ = evaluation.EnableLog(ctx, &pb.Enable{Enable: false})
	}()

	clock, err := evaluation.ClockSync(ctx, &pb.Clock{
		Timesent: timestamppb.New(time.Now()),
	})
	if err != nil {
		log.Fatalf("CollectLogs: cloudn't compare clock time with provider %v: %v",
			prov.ProviderId, err)
	}

	rtt := time.Now().Sub(clock.Timesent.AsTime())
	log.Printf("ClockSync: timesent=%s timerecieved=%s RTT=%s",
		clock.Timesent.AsTime(), clock.Timereceived.AsTime(), rtt)

	devices.Write([]string{
		"scenario",
		"trail",
		"providerId",
		"os",
		"arch",
		"cpus",
		"numJobs",
	}, []string{
		scenarioId,
		trailNum,
		prov.ProviderId,
		prov.Arch.Os,
		prov.Arch.Arch,
		fmt.Sprint(prov.NumCPUs),
		fmt.Sprint(prov.NumJobs),
	})

	stream, err := evaluation.Logs(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("CollectLogs: cloudn't read logging on provider %v: %v",
			prov.ProviderId, err)
	}

	events.Write([]string{
		"scenario",
		"trail",
		"eventId",
		"providerId",
		"timestamp",
		"activity",
		"state",
		"opp-config",
		"opp-runNum",
		"error",
		"byteSize",
		"filename",
	})

	for {
		event, err := stream.Recv()
		if err != nil {
			break
		}

		events.Write([]string{
			scenarioId,
			trailNum,
			event.EventId,
			prov.ProviderId,
			event.Timestamp,
			event.Activity,
			fmt.Sprint(event.State),
			event.Config,
			event.RunNum,
			event.Error,
			fmt.Sprint(event.ByteSize),
			event.Filename,
		})
	}
}

func LogLocal(event Event) (finish func(err error, dlsize uint64) (duration time.Duration)) {

	start := time.Now()

	if !collecting {
		return func(_ error, _ uint64) time.Duration {
			return time.Now().Sub(start)
		}
	}

	id := fmt.Sprintf("%00000000x", rand.Uint32())

	conf, runNum := event.runId()

	evn := &pb.Event{
		Timestamp: time2Tex(start),
		State:     StateStarted,
		EventId:   id,
		Activity:  event.Activity,
		Filename:  event.Filename,
		Config:    conf,
		RunNum:    runNum,
	}

	events.RecordProtos(evn.ProtoReflect())

	return func(err error, dlsize uint64) time.Duration {
		var end = time.Now()
		var state uint32
		var fail string

		if err != nil {
			state = StateFailed
			fail = err.Error()
		} else {
			state = StateFinished
		}

		endEvent := &pb.Event{
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

		events.RecordProtos(endEvent.ProtoReflect())

		return end.Sub(start)
	}
}
