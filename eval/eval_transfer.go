package eval

import (
	"math/rand"
	"sync"
	"time"
)

const (
	TransferDirectionUpload   = "Upload"
	TransferDirectionDownload = "Download"
)

type Transfer struct {
	ScenarioId       string    `csv:"scenario_id"`
	SimulationId     string    `csv:"simulation_id"`
	TrailId          string    `csv:"trail_id"`
	TimeStamp        time.Time `csv:"time_stamp"`
	ProviderId       string    `csv:"provider_id"`
	Step             int       `csv:"step"`
	MatchId          xuint32   `csv:"match_id"`
	Direction        string    `csv:"direction"`
	BytesTransferred uint64    `csv:"bytes_transferred"`
	Error            error     `csv:"error"`
}

type TFeedback struct {
	Success func(bytes uint64) time.Duration
	Error   func(err error) error
}

var tmu sync.Mutex
var trecords []Transfer

func LogTransfer(base Transfer) (feedback *TFeedback) {

	id := xuint32(rand.Uint32())

	start := Transfer{
		ScenarioId:       ScenarioId,
		SimulationId:     SimulationId,
		TrailId:          TrailId,
		TimeStamp:        time.Now(),
		MatchId:          id,
		Step:             StepStart,
		ProviderId:       base.ProviderId,
		Direction:        base.Direction,
		BytesTransferred: base.BytesTransferred,
	}

	tmu.Lock()
	trecords = append(trecords, start)
	tmu.Unlock()

	feedback = &TFeedback{
		Success: func(bytes uint64) time.Duration {
			end := Transfer{
				ScenarioId:       ScenarioId,
				SimulationId:     SimulationId,
				TrailId:          TrailId,
				TimeStamp:        time.Now(),
				MatchId:          id,
				Step:             StepSuccess,
				ProviderId:       base.ProviderId,
				Direction:        base.Direction,
				BytesTransferred: bytes,
			}

			tmu.Lock()
			trecords = append(trecords, end)
			tmu.Unlock()

			return end.TimeStamp.Sub(start.TimeStamp)
		},
		Error: func(err error) error {
			end := Transfer{
				ScenarioId:   ScenarioId,
				SimulationId: SimulationId,
				TrailId:      TrailId,
				TimeStamp:    time.Now(),
				MatchId:      id,
				Step:         StepError,
				ProviderId:   base.ProviderId,
				Direction:    base.Direction,
				Error:        err,
			}

			tmu.Lock()
			trecords = append(trecords, end)
			tmu.Unlock()

			return err
		},
	}

	return
}
