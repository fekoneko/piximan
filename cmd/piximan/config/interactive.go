package config

import (
	"strconv"

	"github.com/fekoneko/piximan/internal/logger"
	"github.com/manifoldco/promptui"
)

func interactive() {
	withSessionId, withRequestParams, resetSession, resetConfig := selectMode()
	sessionId := promptSessionId(withSessionId)
	password := promptPassword(withSessionId)
	defaultMaxPending := promptRequestParamWith(withRequestParams, &defaultMaxPendingPrompt)
	defaultDelay := promptRequestParamWith(withRequestParams, &defaultDelayPrompt)
	pximgMaxPending := promptRequestParamWith(withRequestParams, &pximgMaxPendingPrompt)
	pximgDelay := promptRequestParamWith(withRequestParams, &pximgDelayPrompt)

	config(&options{
		SessionId:         sessionId,
		Password:          password,
		PximgMaxPending:   pximgMaxPending,
		PximgDelay:        pximgDelay,
		DefaultMaxPending: defaultMaxPending,
		DefaultDelay:      defaultDelay,
		ResetSession:      &resetSession,
		ResetConfig:       &resetConfig,
	})
}

func selectMode() (withSessionId bool, withRequestParams bool, resetSession bool, resetConfig bool) {
	_, mode, err := modeSelect.Run()
	logger.MaybeFatal(err, "failed to read configuration mode")

	withSessionId = mode == sessionIdOption
	withRequestParams = mode == requestParamsOption
	resetSession = mode == resetSessionOption
	resetConfig = mode == resetConfigOption

	if !withSessionId && !withRequestParams && !resetSession && !resetConfig {
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
