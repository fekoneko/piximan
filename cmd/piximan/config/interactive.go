package config

import (
	"strconv"

	"github.com/fekoneko/piximan/internal/logger"
	"github.com/manifoldco/promptui"
)

func interactive() {
	withSessionId, withRules, withLimits, resetSession, resetRules, resetLimits, reset := selectMode()
	sessionId := promptSessionId(withSessionId)
	password := promptPassword(withSessionId)
	rules := promptRules(withRules)
	defaultMaxPending := promptLimitWith(withLimits, &defaultMaxPendingPrompt)
	defaultDelay := promptLimitWith(withLimits, &defaultDelayPrompt)
	pximgMaxPending := promptLimitWith(withLimits, &pximgMaxPendingPrompt)
	pximgDelay := promptLimitWith(withLimits, &pximgDelayPrompt)

	config(&options{
		SessionId:       sessionId,
		Password:        password,
		Rules:           rules,
		PximgMaxPending: pximgMaxPending,
		PximgDelay:      pximgDelay,
		MaxPending:      defaultMaxPending,
		Delay:           defaultDelay,
		ResetSession:    &resetSession,
		ResetRules:      &resetRules,
		ResetLimits:     &resetLimits,
		Reset:           &reset,
	})
}

func selectMode() (
	withSessionId, withRules, withLimits,
	resetSession, resetRules, resetLimits, reset bool,
) {
	_, mode, err := modeSelect.Run()
	logger.MaybeFatal(err, "failed to read configuration mode")

	withSessionId = mode == sessionIdOption
	withRules = mode == rulesOption
	withLimits = mode == limitsOption
	resetSession = mode == resetSessionOption
	resetRules = mode == resetRulesOption
	resetLimits = mode == resetLimitsOption
	reset = mode == resetOption

	if !withSessionId && !withRules && !withLimits &&
		!resetSession && !resetRules && !resetLimits && !reset {
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
