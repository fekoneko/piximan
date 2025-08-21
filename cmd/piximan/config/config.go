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

	if utils.FromPtr(options.ResetSession, false) {
		err := c.ResetSessionId()
		logger.MaybeSuccess(err, "session id was removed")
		logger.MaybeFatal(err, "failed to remove session id")
	} else if options.SessionId != nil {
		err = c.WriteSessionId(*options.SessionId)
		logger.MaybeSuccess(err, "session id was set%v",
			utils.If(options.Password != nil, " and encrypted with password", ""),
		)
		logger.MaybeFatal(err, "failed to set session id")
	}

	if utils.FromPtr(options.ResetLimits, false) {
		err := c.ResetLimits()
		logger.MaybeSuccess(err, "configuration parameters were reset")
		logger.MaybeFatal(err, "failed to reset configuration parameters")
		return
	}

	changed := false

	if options.PximgMaxPending != nil {
		c.PximgMaxPending = *options.PximgMaxPending
		changed = true
	}
	if options.PximgDelay != nil {
		c.PximgDelay = time.Duration(*options.PximgDelay) * time.Second
		changed = true
	}
	if options.DefaultMaxPending != nil {
		c.DefaultMaxPending = *options.DefaultMaxPending
		changed = true
	}
	if options.DefaultDelay != nil {
		c.DefaultDelay = time.Duration(*options.DefaultDelay) * time.Second
		changed = true
	}
	if changed {
		err = c.WriteLimits()
		logger.MaybeSuccess(err, "configuration parameters were saved")
		logger.MaybeFatal(err, "failed to save configuration parameters")
	}
}
