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

func selectMode() (bool, bool, bool, bool) {
	_, mode, err := modeSelect.Run()
	logger.MaybeFatal(err, "failed to read configuration mode")

	switch mode {
	case sessionIdOption:
		return true, false, false, false
	case requestParamsOption:
		return false, true, false, false
	case resetSessionOption:
		return false, false, true, false
	case resetConfigOption:
		return false, false, false, true
	default:
		logger.Fatal("incorrect configuration mode: %v", mode)
		panic("unreachable")
	}
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
