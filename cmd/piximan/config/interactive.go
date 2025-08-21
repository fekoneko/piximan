package config

import (
	"strconv"

	"github.com/fekoneko/piximan/internal/logger"
	"github.com/manifoldco/promptui"
)

func interactive() {
	withSessionId, withRequestParams, resetSession, resetLimits := selectMode()
	sessionId := promptSessionId(withSessionId)
	password := promptPassword(withSessionId)
	defaultMaxPending := promptRequestParamWith(withRequestParams, &defaultMaxPendingPrompt)
	defaultDelay := promptRequestParamWith(withRequestParams, &defaultDelayPrompt)
	pximgMaxPending := promptRequestParamWith(withRequestParams, &pximgMaxPendingPrompt)
	pximgDelay := promptRequestParamWith(withRequestParams, &pximgDelayPrompt)

	config(&options{
		SessionId:       sessionId,
		Password:        password,
		PximgMaxPending: pximgMaxPending,
		PximgDelay:      pximgDelay,
		MaxPending:      defaultMaxPending,
		Delay:           defaultDelay,
		ResetSession:    &resetSession,
		ResetLimits:     &resetLimits,
	})
}

func selectMode() (withSessionId bool, withRequestParams bool, resetSession bool, resetLimits bool) {
	_, mode, err := modeSelect.Run()
	logger.MaybeFatal(err, "failed to read configuration mode")

	withSessionId = mode == sessionIdOption
	withRequestParams = mode == requestParamsOption
	resetSession = mode == resetSessionOption
	resetLimits = mode == resetLimitsOption

	if !withSessionId && !withRequestParams && !resetSession && !resetLimits {
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

func promptRequestParamWith(withRequestParams bool, prompt *promptui.Prompt) *uint64 {
	if !withRequestParams {
		return nil
	}

	valueStr, err := prompt.Run()
	logger.MaybeFatal(err, "failed to read request parameter")

	value, err := strconv.ParseUint(valueStr, 10, 64)
	logger.MaybeFatal(err, "cannot parse request parameter")
	return &value
}
