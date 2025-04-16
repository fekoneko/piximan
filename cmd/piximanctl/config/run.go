package config

import (
	"flag"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/flagext"
)

type flags struct {
	sessionId *string
	password  *string
}

func Run() {
	// TODO: use different module for parsing flags that will provide good
	//       way to know if the flag was provided - all of those should be
	//       nil by default
	flags := flags{
		sessionId: flag.String("sessionid", "", ""),
		password:  flag.String("password", "", ""),
	}
	flag.Usage = help.RunConfig
	flag.Parse()

	if flag.NArg() != 0 {
		flagext.BadUsage("too many arguments")
	}

	if flag.NFlag() == 0 {
		interactive()
	} else {
		nonInteractive(flags)
	}
}
