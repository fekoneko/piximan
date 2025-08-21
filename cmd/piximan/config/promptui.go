package config

import (
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

var modeSelectLabel = "What to configure?"
var sessionIdOption = "Authorization (session ID)"
var requestParamsOption = "Request delays and limits"
var resetSessionOption = "Reset configured session ID"
var resetLimitsOption = "Reset request delays and limits"

var modeSelect = promptui.Select{
	Label: modeSelectLabel,
	Items: []string{sessionIdOption, requestParamsOption, resetSessionOption, resetLimitsOption},
}

var sessionIdPrompt = promptui.Prompt{
	Label: "Your session ID",
	Mask:  '*',
}

var passwordPrompt = promptui.Prompt{
	Label: "Encrypt with a password",
	Mask:  '*',
}

var defaultMaxPendingPrompt = promptui.Prompt{
	Label:    "Maximum number of concurrent requests to pixiv.net",
	Validate: utils.ValidateNumber("value must be a number"),
}

var defaultDelayPrompt = promptui.Prompt{
	Label:    "Delay between requests to pixiv.net",
	Validate: utils.ValidateNumber("value must be a number"),
}

var pximgMaxPendingPrompt = promptui.Prompt{
	Label:    "Maximum number of concurrent requests to i.pximg.net",
	Validate: utils.ValidateNumber("value must be a number"),
}

var pximgDelayPrompt = promptui.Prompt{
	Label:    "Delay between requests to i.pximg.net",
	Validate: utils.ValidateNumber("value must be a number"),
}
