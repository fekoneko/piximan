package config

import (
	"path/filepath"
	"time"

	appconfig "github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/termext"
	"github.com/fekoneko/piximan/internal/utils"
)

func config(options *options) {
	termext.DisableInputEcho()
	defer termext.RestoreInputEcho()

	c, err := appconfig.New(options.Password)
	logger.MaybeFatal(err, "failed to open config storage")

	if options.Reset != nil && *options.Reset {
		err := c.Reset()
		logger.MaybeSuccess(err, "configuration was reset")
		logger.MaybeFatal(err, "failed to reset configuration")
		return
	}

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

	if options.ResetRules != nil && *options.ResetRules {
		err := c.ResetRules()
		logger.MaybeSuccess(err, "global download rules were reset")
		logger.MaybeFatal(err, "failed to reset global download rules")

	} else if options.Rules != nil {
		rules := make([]rules.Rules, 0, len(*options.Rules))
		seen := make(map[string]bool, len(*options.Rules))
		for _, rawRulesPath := range *options.Rules {
			rulesPath := filepath.Clean(rawRulesPath)
			if seen[rulesPath] {
				continue
			}
			seen[rulesPath] = true
			r, err := fsext.ReadRules(rulesPath)
			logger.MaybeFatal(err, "cannot read download rules from %v", rulesPath)
			rules = append(rules, *r)
		}
		err := c.SetRules(rules)
		logger.MaybeSuccess(err, "global download rules were set")
		logger.MaybeFatal(err, "failed to set global download rules")
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
