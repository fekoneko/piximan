package syncext

import (
	"sync"
	"time"
)

type RequestGroup struct {
	maxPending uint64
	delay      time.Duration

	numPending        uint64
	numPendingCond    *sync.Cond
	previousTime      time.Time
	previousTimeMutex *sync.Mutex
}

func NewRequestGroup(maxPending uint64, delay time.Duration) *RequestGroup {
	return &RequestGroup{
		maxPending:        maxPending,
		delay:             delay,
		numPending:        0,
		numPendingCond:    sync.NewCond(&sync.Mutex{}),
		previousTime:      time.Time{},
		previousTimeMutex: &sync.Mutex{},
	}
}

func (rg *RequestGroup) Start() {
	rg.previousTimeMutex.Lock()
	defer rg.previousTimeMutex.Unlock()

	duration := time.Until(rg.previousTime.Add(rg.delay))
	time.Sleep(duration)

	rg.numPendingCond.L.Lock()
	defer rg.numPendingCond.L.Unlock()

	for rg.numPending >= rg.maxPending {
		rg.numPendingCond.Wait()
	}

	rg.numPending++
	rg.previousTime = time.Now()
}

func (rg *RequestGroup) Done() {
	rg.numPendingCond.L.Lock()
	rg.numPending--
	rg.numPendingCond.Broadcast()
	rg.numPendingCond.L.Unlock()
}
