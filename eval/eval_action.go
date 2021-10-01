package eval

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"math/rand"
	"time"
)

const (
	ActionCompile  = "Compile"
	ActionCompress = "Compress"
	ActionExtract  = "Extract"
)

func LogAction(action, meta string) (done func(err error) error) {

	timestamp := time.Now()
	ts, _ := timestamp.MarshalText()

	id := fmt.Sprintf("0x%00000000x", rand.Uint32())

	ctx := context.Background()
	_, _ = cli.Action(ctx, &pb.ActionEvent{
		TimeStamp: string(ts),
		DeviceId:  DeviceId,
		Step:      uint32(StepStart),
		EventId:   id,
		Action:    action,
		Meta:      meta,
	})

	done = func(err error) error {
		timestamp = time.Now()
		ts, _ = timestamp.MarshalText()

		if err != nil {
			_, _ = cli.Action(ctx, &pb.ActionEvent{
				TimeStamp: string(ts),
				DeviceId:  DeviceId,
				Step:      uint32(StepError),
				EventId:   id,
				Action:    action,
				Meta:      meta,
				Error:     err.Error(),
			})
		} else {
			_, _ = cli.Action(ctx, &pb.ActionEvent{
				TimeStamp: string(ts),
				DeviceId:  DeviceId,
				Step:      uint32(StepSuccess),
				EventId:   id,
				Action:    action,
				Meta:      meta,
			})
		}

		return err
	}

	return
}
