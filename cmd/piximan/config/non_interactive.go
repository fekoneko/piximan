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

	if options.SessionId != nil && options.ResetSession != nil {
		fmt.Println("`-s, --session-id' cannot be used with `--no-session' flag")
		os.Exit(2)
	}

	if (options.PximgMaxPending != nil || options.PximgDelay != nil ||
		options.DefaultMaxPending != nil || options.DefaultDelay != nil) &&
		options.ResetConfig != nil {
		fmt.Println("no configuration parameters can be used with `--default' flag")
		os.Exit(2)
	}

	config(options)
}
