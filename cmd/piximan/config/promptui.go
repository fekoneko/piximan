package config

import (
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

var modeSelectLabel = "What to configure?"
var sessionIdOption = "Authorization (session ID)"
var defaultsOption = "Downloader defaults (image size, language)"
var rulesOption = "Global download rules"
var limitsOption = "Request delays and limits"
var resetSessionOption = "Reset configured session ID"
var resetDefaultsOption = "Reset downloader defaults (image size, language)"
var resetRulesOption = "Reset global download rules"
var resetLimitsOption = "Reset request delays and limits"
var resetOption = "Reset all configuration"

var modeSelect = promptui.Select{
	Label: modeSelectLabel,
	Items: []string{
		sessionIdOption, defaultsOption, rulesOption, limitsOption,
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

var thumbnailSizeOption = "Thumbnail"
var smallSizeOption = "Small"
var mediumSizeOption = "Medium"
var originalSizeOption = "Original"

var sizeSelect = promptui.Select{
	Label:     "Default size of downloaded images",
	Items:     []string{thumbnailSizeOption, smallSizeOption, mediumSizeOption, originalSizeOption},
	CursorPos: 3,
}

var japaneseOption = "Japanese (or other original language)"
var englishOption = "English"

var languageSelect = promptui.Select{
	Label: "Default language of work titles and descriptions",
	Items: []string{japaneseOption, englishOption},
}

var rulesPrompt = promptui.Prompt{
	Label: "Paths to download rules files (comma-separated)",
}

var maxPendingPrompt = promptui.Prompt{
	Label:    "Maximum number of concurrent requests to pixiv.net",
	Validate: utils.ValidateUint("value must be a number"),
}

var delayPrompt = promptui.Prompt{
	Label:    "Delay between requests to pixiv.net",
	Validate: utils.ValidateUint("value must be a number"),
}

var pximgMaxPendingPrompt = promptui.Prompt{
	Label:    "Maximum number of concurrent requests to i.pximg.net",
	Validate: utils.ValidateUint("value must be a number"),
}

var pximgDelayPrompt = promptui.Prompt{
	Label:    "Delay between requests to i.pximg.net",
	Validate: utils.ValidateUint("value must be a number"),
}
