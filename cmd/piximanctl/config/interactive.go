package config

import (
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/manifoldco/promptui"
)

func interactive() {
	sessionId := promptSessionId()

	var password *string
	if len(sessionId) != 0 {
		p := promptPassword()
		password = &p
	}

	config(&options{
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
	logext.MaybeFatal(err, "failed to read session id")
	return sessionId
}

var passwordPrompt = promptui.Prompt{
	Label: "Encrypt with a password",
	Mask:  '*',
}

func promptPassword() string {
	password, err := passwordPrompt.Run()
	logext.MaybeFatal(err, "failed to read password")
	return password
}
