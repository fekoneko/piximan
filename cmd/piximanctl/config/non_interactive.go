package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	_, err := flags.Parse(options)
	if err != nil {
		os.Exit(2)
	}

	configSessionId(options)
}
