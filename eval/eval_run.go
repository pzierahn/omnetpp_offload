package eval

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"math/rand"
	"time"
)

func LogRun(provider, config, num string) (done func(err error) error) {

	timestamp := time.Now()
	ts, _ := timestamp.MarshalText()

	id := fmt.Sprintf("0x%00000000x", rand.Uint32())

	ctx := context.Background()
	_, _ = cli.Run(ctx, &pb.RunEvent{
		TimeStamp:  string(ts),
		ProviderId: provider,
		Step:       uint32(StepStart),
		EventId:    id,
		RunConfig:  config,
		RunNumber:  num,
	})

	done = func(err error) error {
		timestamp = time.Now()
		ts, _ = timestamp.MarshalText()

		if err != nil {
			_, _ = cli.Run(ctx, &pb.RunEvent{
				TimeStamp:  string(ts),
				ProviderId: provider,
				Step:       uint32(StepError),
				EventId:    id,
				RunConfig:  config,
				RunNumber:  num,
				Error:      err.Error(),
			})
		} else {
			_, _ = cli.Run(ctx, &pb.RunEvent{
				TimeStamp:  string(ts),
				ProviderId: provider,
				Step:       uint32(StepSuccess),
				EventId:    id,
				RunConfig:  config,
				RunNumber:  num,
			})
		}

		return err
	}

	return
}
