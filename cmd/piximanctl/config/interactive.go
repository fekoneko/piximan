package config

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

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

	configSessionId(flags)
}
