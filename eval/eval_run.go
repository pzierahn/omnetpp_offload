package eval

import (
	"math/rand"
	"sync"
	"time"
)

const (
	CommandCompile   = "Compile"
	CommandCheckout  = "Checkout"
	CommandExecution = "Execution"
)

type Run struct {
	ScenarioId   string    `csv:"scenario_id"`
	SimulationId string    `csv:"simulation_id"`
	TrailId      string    `csv:"trail_id"`
	TimeStamp    time.Time `csv:"time_stamp"`
	ProviderId   string    `csv:"provider_id"`
	Command      string    `csv:"command"`
	Step         int       `csv:"step"`
	MatchId      xuint32   `csv:"match_id"`
	RunConfig    string    `csv:"run_config"`
	RunNumber    string    `csv:"run_number"`
	Error        error     `csv:"error"`
}

type RFeedback struct {
	Success func() time.Duration
	Error   func(err error) error
}

var rmu sync.Mutex
var rrecords []Run

func LogRun(base Run) (feedback *RFeedback) {

	id := xuint32(rand.Uint32())

	start := Run{
		ScenarioId:   ScenarioId,
		SimulationId: SimulationId,
		TrailId:      TrailId,
		TimeStamp:    time.Now(),
		MatchId:      id,
		Step:         StepStart,
		Command:      base.Command,
		ProviderId:   base.ProviderId,
		RunConfig:    base.RunConfig,
		RunNumber:    base.RunNumber,
	}

	rmu.Lock()
	rrecords = append(rrecords, start)
	rmu.Unlock()

	feedback = &RFeedback{
		Success: func() time.Duration {
			end := Run{
				ScenarioId:   ScenarioId,
				SimulationId: SimulationId,
				TrailId:      TrailId,
				TimeStamp:    time.Now(),
				MatchId:      id,
				Step:         StepSuccess,
				Command:      base.Command,
				ProviderId:   base.ProviderId,
				RunConfig:    base.RunConfig,
				RunNumber:    base.RunNumber,
			}

			rmu.Lock()
			rrecords = append(rrecords, end)
			rmu.Unlock()

			return end.TimeStamp.Sub(start.TimeStamp)
		},
		Error: func(err error) error {
			end := Run{
				ScenarioId:   ScenarioId,
				SimulationId: SimulationId,
				TrailId:      TrailId,
				TimeStamp:    time.Now(),
				MatchId:      id,
				Step:         StepError,
				Command:      base.Command,
				ProviderId:   base.ProviderId,
				RunConfig:    base.RunConfig,
				RunNumber:    base.RunNumber,
				Error:        err,
			}

			rmu.Lock()
			rrecords = append(rrecords, end)
			rmu.Unlock()

			return err
		},
	}

	return
}
