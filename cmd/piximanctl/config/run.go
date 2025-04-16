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
		interactive(flags)
	} else {
		nonInteractive(flags)
	}
}
