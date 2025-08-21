package config

import (
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

var modeSelectLabel = "What to configure?"
var sessionIdOption = "Authorization (session ID)"
var rulesOption = "Global download rules"
var limitsOption = "Request delays and limits"
var resetSessionOption = "Reset configured session ID"
var resetRulesOption = "Reset global download rules"
var resetLimitsOption = "Reset request delays and limits"
var resetOption = "Reset all configuration"

var modeSelect = promptui.Select{
	Label: modeSelectLabel,
	Items: []string{
		sessionIdOption, rulesOption, limitsOption,
		resetSessionOption, resetRulesOption, resetLimitsOption, resetOption,
	},
}

var sessionIdPrompt = promptui.Prompt{
	Label: "Your session ID",
	Mask:  '*',
}

var passwordPrompt = promptui.Prompt{
	Label: "Encrypt with a password",
	Mask:  '*',
}

var rulesPrompt = promptui.Prompt{
	Label: "Paths to download rules files (comma-separated)",
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
