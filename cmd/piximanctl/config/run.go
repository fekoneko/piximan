package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
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

func interactive(flags flags) {
	fmt.Print("Your session ID: ")
	fmt.Scanln(flags.sessionId)
	fmt.Print("Encrypt session ID with a password: ")
	fmt.Scanln(flags.password) // TODO: hide password with asterisks

	setSessionId(flags)
}

func nonInteractive(flags flags) {
	if flagext.Provided("password") && !flagext.Provided("sessionid") {
		flagext.BadUsage("flag -password requires -sessionid to be specified")
	}

	if flagext.Provided("sessionid") {
		setSessionId(flags)
	}
}

func setSessionId(flags flags) {
	storage, err := secretstorage.Open(*flags.password)
	if err != nil {
		fmt.Printf("failed to set session id: %v\n", err)
		os.Exit(1)
	}
	if err := storage.StoreSessionId(*flags.sessionId); err != nil {
		fmt.Printf("failed to set session id: %v\n", err)
		os.Exit(1)
	}
}
