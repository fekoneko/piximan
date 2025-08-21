package config

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	if args, err := flags.Parse(options); err != nil {
		os.Exit(1)
	} else if len(args) > 0 {
		fmt.Println("extra arguments provided")
		os.Exit(2)
	}

	if options.Password != nil && options.SessionId == nil {
		fmt.Println("`-P, --password' can only be used with `-s, --session-id' flag")
		os.Exit(2)
	}

	if options.SessionId != nil && options.ResetSession != nil {
		fmt.Println("`-s, --session-id' cannot be used with `--reset-session' flag")
		os.Exit(2)
	}

	if (options.PximgMaxPending != nil || options.PximgDelay != nil ||
		options.DefaultMaxPending != nil || options.DefaultDelay != nil) &&
		options.ResetLimits != nil {
		fmt.Println("request delays and limits parameters cannot be used with `--reset-limits' flag")
		os.Exit(2)
	}

	config(options)
}
