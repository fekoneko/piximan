package config

import (
	"time"

	appconfig "github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/termext"
	"github.com/fekoneko/piximan/internal/utils"
)

func config(options *options) {
	termext.DisableInputEcho()
	defer termext.RestoreInputEcho()

	c, err := appconfig.New(options.Password)
	logger.MaybeFatal(err, "failed to open config storage")

	if options.ResetSession != nil && *options.ResetSession {
		err := c.ResetSessionId()
		logger.MaybeSuccess(err, "session id was reset")
		logger.MaybeFatal(err, "failed to reset session id")

	} else if options.SessionId != nil {
		err = c.SetSessionId(*options.SessionId)
		logger.MaybeSuccess(err, "session id was set%v",
			utils.If(options.Password != nil, " and encrypted with password", ""),
		)
		logger.MaybeFatal(err, "failed to set session id")
	}

	if options.ResetLimits != nil && *options.ResetLimits {
		err := c.ResetLimits()
		logger.MaybeSuccess(err, "request delays and limits were reset")
		logger.MaybeFatal(err, "failed to reset request delays and limits")

	} else if options.MaxPending != nil || options.Delay != nil ||
		options.PximgMaxPending != nil || options.PximgDelay != nil {
		l, err := c.Limits()
		logger.MaybeFatal(err, "failed to read request delays and limits")

		if options.MaxPending != nil {
			l.MaxPending = *options.MaxPending
		}
		if options.Delay != nil {
			l.Delay = time.Duration(*options.Delay) * time.Second
		}
		if options.PximgMaxPending != nil {
			l.PximgMaxPending = *options.PximgMaxPending
		}
		if options.PximgDelay != nil {
			l.PximgDelay = time.Duration(*options.PximgDelay) * time.Second
		}

		err = c.SetLimits(l)
		logger.MaybeSuccess(err, "request delays and limits were configured")
		logger.MaybeFatal(err, "failed to configure request delays and limits")
	}
}
