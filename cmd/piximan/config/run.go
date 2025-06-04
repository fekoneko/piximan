package config

import "os"

type options struct {
	SessionId string  `short:"s" long:"session-id" required:"true"`
	Password  *string `short:"P" long:"password"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
