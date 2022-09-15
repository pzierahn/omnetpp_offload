package eval

import (
	"context"
	"fmt"
	"github.com/pzierahn/omnetpp_offload/csv"
	"github.com/pzierahn/omnetpp_offload/gconfig"
	pb "github.com/pzierahn/omnetpp_offload/proto"
	"github.com/pzierahn/omnetpp_offload/stargrpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var collecting = false
var scenarioId string
var trailNum string

var events *csv.Writer
var devices *csv.Writer

func init() {
	scenario, trail := os.Getenv("SCENARIO"), os.Getenv("TRAIL")
	if scenario != "" && trail != "" {
		return
	}

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

	devices.Write([]string{
		"scenario",
		"trail",
		"providerId",
		"os",
		"arch",
		"cpus",
		"numJobs",
		"timesent",
		"timerecieved",
		"rtt",
		"connect",
	})
}

func CollectLogs(client *grpc.ClientConn, prov *pb.ProviderInfo, connect int) {
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
		scenarioId,
		trailNum,
		prov.ProviderId,
		prov.Arch.Os,
		prov.Arch.Arch,
		fmt.Sprint(prov.NumCPUs),
		fmt.Sprint(prov.NumJobs),
		fmt.Sprint(clock.Timesent.AsTime()),
		fmt.Sprint(clock.Timereceived.AsTime()),
		fmt.Sprint(rtt),
		stargrpc.ConnectionToName(connect),
	})

	stream, err := evaluation.Logs(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("CollectLogs: cloudn't read logging on provider %v: %v",
			prov.ProviderId, err)
	}

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

	events.Write([]string{
		scenarioId,
		trailNum,
		id,
		"consumer",
		time2Tex(start),
		event.Activity,
		fmt.Sprint(StateStarted),
		conf,
		runNum,
		"",
		"",
		event.Filename,
	})

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

		events.Write([]string{
			scenarioId,
			trailNum,
			id,
			"consumer",
			time2Tex(end),
			event.Activity,
			fmt.Sprint(state),
			conf,
			runNum,
			fail,
			fmt.Sprint(dlsize),
			event.Filename,
		})

		return end.Sub(start)
	}
}
