package eval

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/project.go.omnetpp/proto"
	"math/rand"
	"time"
)

const (
	TransferDirectionUpload   = "Upload"
	TransferDirectionDownload = "Download"
)

func LogTransfer(provider, direction, file string) (done func(dlsize uint64, err error) error) {

	timestamp := time.Now()
	ts, _ := timestamp.MarshalText()

	id := fmt.Sprintf("0x%00000000x", rand.Uint32())

	ctx := context.Background()
	_, _ = client.Transfer(ctx, &pb.TransferEvent{
		TimeStamp:  string(ts),
		ProviderId: provider,
		Step:       uint32(StepStart),
		EventId:    id,
		Direction:  direction,
		File:       file,
	})

	done = func(size uint64, err error) error {
		timestamp = time.Now()
		ts, _ = timestamp.MarshalText()

		if err != nil {
			_, _ = client.Transfer(ctx, &pb.TransferEvent{
				TimeStamp:   string(ts),
				ProviderId:  provider,
				Step:        uint32(StepError),
				EventId:     id,
				Direction:   direction,
				File:        file,
				Transferred: 0,
				Error:       err.Error(),
			})
		} else {
			_, _ = client.Transfer(ctx, &pb.TransferEvent{
				TimeStamp:   string(ts),
				ProviderId:  provider,
				Step:        uint32(StepSuccess),
				EventId:     id,
				Direction:   direction,
				File:        file,
				Transferred: size,
			})
		}

		return err
	}

	return
}
