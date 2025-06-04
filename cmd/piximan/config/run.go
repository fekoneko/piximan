package config

import "os"

type options struct {
	SessionId         *string `short:"s" long:"session-id"`
	Password          *string `short:"P" long:"password"`
	PximgMaxPending   *uint64 `short:"M" long:"image-max-pending"`
	PximgDelay        *uint64 `short:"D" long:"image-delay"`
	DefaultMaxPending *uint64 `short:"m" long:"max-pending"`
	DefaultDelay      *uint64 `short:"d" long:"delay"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
