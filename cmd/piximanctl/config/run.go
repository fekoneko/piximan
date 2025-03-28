package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
)

func Run() {
	sessionId := flag.String("sessionid", "", "")
	password := flag.String("password", "", "")
	flag.Usage = help.RunConfig
	flag.Parse()

	if flag.NArg() != 0 {
		flagext.BadUsage("too many arguments")
	}

	if flag.NFlag() == 0 {
		interactive(sessionId, password)
	} else {
		nonInteractive(sessionId, password)
	}
}

func interactive(sessionId *string, password *string) {
	fmt.Print("Your session ID: ")
	fmt.Scanln(sessionId)
	fmt.Print("Encrypt session ID with a password: ")
	fmt.Scanln(password) // TODO: hide password with asterisks

	setSessionId(sessionId, password)
}

func nonInteractive(sessionId *string, password *string) {
	if flagext.Provided("password") && !flagext.Provided("sessionid") {
		flagext.BadUsage("flag -password requires -sessionid to be specified")
	}

	if flagext.Provided("sessionid") {
		setSessionId(sessionId, password)
	}
}

func setSessionId(sessionId *string, password *string) {
	storage, err := secretstorage.Open(*password)
	if err != nil {
		fmt.Printf("failed to set session id: %v\n", err)
		os.Exit(1)
	}
	if err := storage.StoreSessionId(*sessionId); err != nil {
		fmt.Printf("failed to set session id: %v\n", err)
		os.Exit(1)
	}
}
