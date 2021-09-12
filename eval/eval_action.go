package eval

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"math/rand"
	"time"
)

const (
	ActionCompile = "Compile"
	//ActionCompress = "Compress"
	//ActionExtract  = "Extract"
)

func LogAction(provider, action string) (done func(err error) error) {

	timestamp := time.Now()
	ts, _ := timestamp.MarshalText()

	id := fmt.Sprintf("0x%00000000x", rand.Uint32())

	ctx := context.Background()
	_, _ = client.Action(ctx, &pb.ActionEvent{
		TimeStamp:  string(ts),
		ProviderId: provider,
		Step:       uint32(StepStart),
		EventId:    id,
		Action:     action,
	})

	done = func(err error) error {
		timestamp = time.Now()
		ts, _ = timestamp.MarshalText()

		if err != nil {
			_, _ = client.Action(ctx, &pb.ActionEvent{
				TimeStamp:  string(ts),
				ProviderId: provider,
				Step:       uint32(StepError),
				EventId:    id,
				Action:     action,
				Error:      err.Error(),
			})
		} else {
			_, _ = client.Action(ctx, &pb.ActionEvent{
				TimeStamp:  string(ts),
				ProviderId: provider,
				Step:       uint32(StepSuccess),
				EventId:    id,
				Action:     action,
			})
		}

		return err
	}

	return
}
