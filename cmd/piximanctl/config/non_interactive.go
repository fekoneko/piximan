package config

import "github.com/fekoneko/piximan/pkg/flagext"

func nonInteractive(flags flags) {
	if flagext.Provided("password") && !flagext.Provided("sessionid") {
		flagext.BadUsage("flag -password requires -sessionid to be specified")
	}

	if flagext.Provided("sessionid") {
		configSessionId(flags)
	}
}
