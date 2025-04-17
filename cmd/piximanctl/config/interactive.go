package config

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func interactive() {
	sessionId := promptSessionId()

	var password *string
	if len(sessionId) != 0 {
		p := promptPassword()
		password = &p
	}

	configSessionId(&options{
		SessionId: sessionId,
		Password:  password,
	})
}

var sessionIdPrompt = promptui.Prompt{
	Label: "Your session ID",
	Mask:  '*',
}

func promptSessionId() string {
	sessionId, err := sessionIdPrompt.Run()
	if err != nil {
		fmt.Printf("failed to read session id: %v\n", err)
		os.Exit(1)
	}
	return sessionId
}

var passwordPrompt = promptui.Prompt{
	Label: "Encrypt with a password",
	Mask:  '*',
}

func promptPassword() string {
	password, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("failed to read password: %v\n", err)
		os.Exit(1)
	}
	return password
}
