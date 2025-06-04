package config

import "os"

type options struct {
	SessionId string  `short:"s" long:"sessionid" required:"true"`
	Password  *string `short:"p" long:"password"`
}

func Run() {
	if len(os.Args) <= 1 {
		interactive()
	} else {
		nonInteractive()
	}
}
