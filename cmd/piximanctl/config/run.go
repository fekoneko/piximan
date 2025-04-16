package config

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
	"golang.org/x/term"
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
	if len(*flags.sessionId) != 0 {
		fmt.Print("Encrypt session ID with a password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Printf("failed to read password: %v\n", err)
			os.Exit(1)
		}
		*flags.password = string(password)
	}

	continueSessionId(flags)
}

func nonInteractive(flags flags) {
	if flagext.Provided("password") && !flagext.Provided("sessionid") {
		flagext.BadUsage("flag -password requires -sessionid to be specified")
	}

	if flagext.Provided("sessionid") {
		continueSessionId(flags)
	}
}

func continueSessionId(flags flags) {
	storage, err := secretstorage.Open(*flags.password)
	if err != nil {
		fmt.Printf("failed to set session id: %v\n", err)
		os.Exit(1)
	}
	if len(*flags.sessionId) == 0 {
		if err := storage.RemoveSessionId(); err != nil {
			fmt.Printf("failed to remove session id: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := storage.StoreSessionId(*flags.sessionId); err != nil {
			fmt.Printf("failed to set session id: %v\n", err)
			os.Exit(1)
		}
	}
}
