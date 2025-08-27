package config

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
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
	if options.Size != nil && !imageext.ValidSizeUint(*options.Size) {
		fmt.Println("invalid argument for flag `-s, --size'")
		os.Exit(2)
	}
	if options.Language != nil && !work.ValidLanguageString(*options.Language) {
		fmt.Println("invalid argument for flag `-L, --language'")
		os.Exit(2)
	}
	if utils.SomeDefined(options.Size, options.Language) && options.ResetDefaults != nil {
		fmt.Println("`-s, --size' and `-L, --language' cannot be used with `--reset-defaults' flag")
		os.Exit(2)
	}
	if options.Rules != nil && options.ResetRules != nil {
		fmt.Println("`-r, --rules' cannot be used with `--reset-rules' flag")
		os.Exit(2)
	}
	if utils.SomeDefined(options.MaxPending, options.Delay, options.PximgMaxPending, options.PximgDelay) &&
		options.ResetLimits != nil {
		fmt.Println("request delays and limits parameters cannot be used with `--reset-limits' flag")
		os.Exit(2)
	}
	if utils.SomeDefined(
		options.SessionId, options.Password, options.Size, options.Language, options.Rules,
		options.MaxPending, options.Delay, options.PximgMaxPending, options.PximgDelay,
		options.ResetSession, options.ResetRules, options.ResetLimits,
	) && options.Reset != nil {
		fmt.Println(("no other flags can be used when `--reset' is provided"))
		os.Exit(2)
	}

	config(options)
}
