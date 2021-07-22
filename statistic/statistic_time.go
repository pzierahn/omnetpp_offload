package statistic

import (
	"fmt"
	"sync"
	"time"
)

type TimeStatistic struct {
	Events int
	Sum    string
	Avg    string
}

type Time struct {
	mu    sync.RWMutex
	items []time.Duration
}

func (stat *Time) Add(items ...time.Duration) {

	stat.mu.Lock()
	defer stat.mu.Unlock()

	stat.items = append(stat.items, items...)
}

func (stat *Time) Until(item time.Time) (duration time.Duration) {

	stat.mu.Lock()
	defer stat.mu.Unlock()

	duration = time.Now().Sub(item)
	stat.items = append(stat.items, duration)

	return
}

func (stat *Time) Avg() (sum, avg time.Duration) {

	stat.mu.RLock()
	defer stat.mu.RUnlock()

	if len(stat.items) == 0 {
		return
	}

	for _, dur := range stat.items {
		sum += dur
	}

	avg = time.Duration(sum.Nanoseconds() / int64(len(stat.items)))

	return
}

func (stat *Time) Export() (data *TimeStatistic) {

	if stat == nil {
		return
	}

	stat.mu.RLock()
	defer stat.mu.RUnlock()

	if len(stat.items) == 0 {
		return
	}

	sum, avg := stat.Avg()

	data = &TimeStatistic{
		Events: len(stat.items),
		Sum:    fmt.Sprintf("%v", sum),
		Avg:    fmt.Sprintf("%v", avg),
	}

	return
}
