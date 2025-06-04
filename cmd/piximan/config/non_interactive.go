package config

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	if _, err := flags.Parse(options); err != nil {
		os.Exit(2)
	}

	if options.Password != nil && options.SessionId == nil {
		fmt.Println("`-P, --password' can only be used with `-s, --session-id' flag")
		os.Exit(2)
	}

	config(options)
}
