package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
	"github.com/manifoldco/promptui"
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

// TODO: write descriptions
var sessionIdPrompt = promptui.Prompt{
	Label: "Your session ID",
	Mask:  '*',
}
var passwordPrompt = promptui.Prompt{
	Label: "Encrypt with a password",
	Mask:  '*',
}

func interactive(flags flags) {
	sessionId, err := sessionIdPrompt.Run()
	if err != nil {
		fmt.Printf("failed to read session id: %v\n", err)
		os.Exit(1)
	}
	*flags.sessionId = sessionId

	if len(*flags.sessionId) != 0 {
		password, err := passwordPrompt.Run()
		if err != nil {
			fmt.Printf("failed to read password: %v\n", err)
			os.Exit(1)
		}
		*flags.password = password
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
