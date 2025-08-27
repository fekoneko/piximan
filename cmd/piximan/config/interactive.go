package config

import (
	"strconv"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/manifoldco/promptui"
)

func interactive() {
	withSessionId, withDefaults, withRules, withLimits,
		resetSession, resetDefaults, resetRules, resetLimits, reset := selectMode()
	sessionId := promptSessionId(withSessionId)
	password := promptPassword(withSessionId)
	size := selectSize(withDefaults)
	language := selectLanguage(withDefaults)
	rules := promptRules(withRules)
	defaultMaxPending := promptLimitWith(withLimits, &maxPendingPrompt)
	defaultDelay := promptLimitWith(withLimits, &delayPrompt)
	pximgMaxPending := promptLimitWith(withLimits, &pximgMaxPendingPrompt)
	pximgDelay := promptLimitWith(withLimits, &pximgDelayPrompt)

	config(&options{
		SessionId:       sessionId,
		Password:        password,
		Size:            size,
		Language:        language,
		Rules:           rules,
		PximgMaxPending: pximgMaxPending,
		PximgDelay:      pximgDelay,
		MaxPending:      defaultMaxPending,
		Delay:           defaultDelay,
		ResetSession:    &resetSession,
		ResetDefaults:   &resetDefaults,
		ResetRules:      &resetRules,
		ResetLimits:     &resetLimits,
		Reset:           &reset,
	})
}

func selectMode() (
	withSessionId, withDefaults, withRules, withLimits,
	resetSession, resetDefaults, resetRules, resetLimits, reset bool,
) {
	_, mode, err := modeSelect.Run()
	logger.MaybeFatal(err, "failed to read configuration mode")

	withSessionId = mode == sessionIdOption
	withDefaults = mode == defaultsOption
	withRules = mode == rulesOption
	withLimits = mode == limitsOption
	resetSession = mode == resetSessionOption
	resetDefaults = mode == resetDefaultsOption
	resetRules = mode == resetRulesOption
	resetLimits = mode == resetLimitsOption
	reset = mode == resetOption

	if !withSessionId && !withDefaults && !withRules && !withLimits &&
		!resetSession && !resetDefaults && !resetRules && !resetLimits && !reset {
		logger.Fatal("incorrect configuration mode: %v", mode)
	}
	return
}

func promptSessionId(withSessionId bool) *string {
	if !withSessionId {
		return nil
	}

	sessionId, err := sessionIdPrompt.Run()
	logger.MaybeFatal(err, "failed to read session id")
	return &sessionId
}

func promptPassword(withSessionId bool) *string {
	if !withSessionId {
		return nil
	}

	password, err := passwordPrompt.Run()
	logger.MaybeFatal(err, "failed to read password")
	if password == "" {
		return nil
	}
	return &password
}

func selectSize(withDefaults bool) *uint64 {
	if !withDefaults {
		return nil
	}

	_, size, err := sizeSelect.Run()
	logger.MaybeFatal(err, "failed to read size")

	switch size {
	case thumbnailSizeOption:
		result := imageext.SizeThumbnail.ToUint()
		return &result
	case smallSizeOption:
		result := imageext.SizeSmall.ToUint()
		return &result
	case mediumSizeOption:
		result := imageext.SizeMedium.ToUint()
		return &result
	case originalSizeOption:
		result := imageext.SizeOriginal.ToUint()
		return &result
	default:
		logger.Fatal("incorrect size: %v", size)
		panic("unreachable")
	}
}

func selectLanguage(withDefaults bool) *string {
	if !withDefaults {
		return nil
	}

	_, language, err := languageSelect.Run()
	logger.MaybeFatal(err, "failed to read language choice")

	switch language {
	case japaneseOption:
		return utils.ToPtr(work.LanguageJapaneseString)
	case englishOption:
		return utils.ToPtr(work.LanguageEnglishString)
	default:
		logger.Fatal("incorrect language: %v", language)
		panic("unreachable")
	}
}

func promptRules(withRules bool) *[]string {
	if !withRules {
		return nil
	}

	rules, err := rulesPrompt.Run()
	logger.MaybeFatal(err, "failed to read download rules")
	parsed := parseStrings(rules)
	if len(parsed) == 0 {
		return nil
	}
	return &parsed
}

func promptLimitWith(withLimits bool, prompt *promptui.Prompt) *uint64 {
	if !withLimits {
		return nil
	}

	valueStr, err := prompt.Run()
	logger.MaybeFatal(err, "failed to read request parameter")

	value, err := strconv.ParseUint(valueStr, 10, 64)
	logger.MaybeFatal(err, "cannot parse request parameter")
	return &value
}
