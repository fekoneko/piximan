package config

import "os"

type options struct {
	SessionId       *string   `short:"s" long:"session-id"`
	Password        *string   `short:"P" long:"password"`
	Rules           *[]string `short:"r" long:"rules"`
	MaxPending      *uint64   `short:"m" long:"max-pending"`
	Delay           *uint64   `short:"d" long:"delay"`
	PximgMaxPending *uint64   `short:"M" long:"pximg-max-pending"`
	PximgDelay      *uint64   `short:"D" long:"pximg-delay"`
	ResetSession    *bool     `long:"reset-session"`
	ResetRules      *bool     `long:"reset-rules"`
	ResetLimits     *bool     `long:"reset-limits"`
	Reset           *bool     `long:"reset"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
