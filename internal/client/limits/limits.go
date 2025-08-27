package limits

import (
	"time"
)

type Limits struct {
	MaxPending      uint64
	Delay           time.Duration
	PximgMaxPending uint64
	PximgDelay      time.Duration
}

const (
	DefaultMaxPending      = 1
	DefaultDelay           = time.Second * 2
	DefaultPximgMaxPending = 5
	DefaultPximgDelay      = time.Second * 1
)

func Default() *Limits {
	return &Limits{
		MaxPending:      DefaultMaxPending,
		Delay:           DefaultDelay,
		PximgMaxPending: DefaultPximgMaxPending,
		PximgDelay:      DefaultPximgDelay,
	}
}

func (l *Limits) IsDefault() bool {
	return l.MaxPending == DefaultMaxPending &&
		l.Delay == DefaultDelay &&
		l.PximgMaxPending == DefaultPximgMaxPending &&
		l.PximgDelay == DefaultPximgDelay
}
