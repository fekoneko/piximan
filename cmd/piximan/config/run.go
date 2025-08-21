package config

import "os"

type options struct {
	SessionId       *string `short:"s" long:"session-id"`
	Password        *string `short:"P" long:"password"`
	MaxPending      *uint64 `short:"m" long:"max-pending"`
	Delay           *uint64 `short:"d" long:"delay"`
	PximgMaxPending *uint64 `short:"M" long:"pximg-max-pending"`
	PximgDelay      *uint64 `short:"D" long:"pximg-delay"`
	ResetSession    *bool   `long:"reset-session"`
	ResetLimits     *bool   `long:"reset-limits"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
