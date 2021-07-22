package statistic

import (
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"sync"
	"time"
)

type LoadStatistic struct {
	Transfers   int
	Transferred string
	Duration    string
	Speed       string
}

type Transfer struct {
	Duration time.Duration
	Size     uint64
}

type Load struct {
	mu    sync.RWMutex
	items []Transfer
}

func (stat *Load) Add(items ...Transfer) {

	stat.mu.Lock()
	defer stat.mu.Unlock()

	stat.items = append(stat.items, items...)
}

func (stat *Load) Export() (data *LoadStatistic) {

	if stat == nil {
		return
	}

	stat.mu.RLock()
	defer stat.mu.RUnlock()

	if len(stat.items) == 0 {
		return
	}

	var transfers int
	var transferred uint64
	var duration time.Duration

	for _, item := range stat.items {
		transfers++
		transferred += item.Size
		duration += item.Duration
	}

	sec := uint64(duration.Seconds())
	if sec == 0 {
		sec = 1
	}

	speed := transferred / sec

	data = &LoadStatistic{
		Transfers:   transfers,
		Transferred: simple.ByteSize(transferred),
		Duration:    fmt.Sprintf("%v", duration),
		Speed:       fmt.Sprintf("%v/s", simple.ByteSize(speed)),
	}

	return
}
