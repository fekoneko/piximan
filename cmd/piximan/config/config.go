package config

import (
	"path/filepath"
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
	appconfig "github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/downloader/rules"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
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

	if options.ResetDefaults != nil && *options.ResetDefaults {
		err := c.ResetDefaults()
		logger.MaybeSuccess(err, "downloader defaults were reset")
		logger.MaybeFatal(err, "failed to reset downloader defaults")

	} else if options.Size != nil || options.Language != nil {
		d, warning, err := c.Defaults()
		logger.MaybeWarning(warning, "while reading downloader defaults")
		logger.MaybeFatal(err, "failed to read downloader defaults")

		if options.Size != nil {
			d.Size = imageext.SizeFromUint(*options.Size)
		}
		if options.Language != nil {
			d.Language = work.LanguageFromString(*options.Language)
		}

		err = c.SetDefaults(d)
		logger.MaybeSuccess(err, "downloader defaults were configured")
		logger.MaybeFatal(err, "failed to configure downloader defaults")
	}

	if options.ResetRules != nil && *options.ResetRules {
		err := c.ResetRules()
		logger.MaybeSuccess(err, "global download rules were reset")
		logger.MaybeFatal(err, "failed to reset global download rules")

	} else if options.Rules != nil {
		rs := make([]rules.Rules, 0, len(*options.Rules))
		seen := make(map[string]bool, len(*options.Rules))
		for _, rawRulesPath := range *options.Rules {
			rulesPath := filepath.Clean(rawRulesPath)
			if seen[rulesPath] {
				continue
			}
			seen[rulesPath] = true
			r, warning, err := fsext.ReadRules(rulesPath)
			logger.MaybeWarning(warning, "while reading download rules from %v", rulesPath)
			logger.MaybeFatal(err, "cannot read download rules from %v", rulesPath)
			rs = append(rs, *r)
		}
		err := c.SetRules(rs)
		logger.MaybeSuccess(err, "global download rules were set")
		logger.MaybeFatal(err, "failed to set global download rules")
	}

	if options.ResetLimits != nil && *options.ResetLimits {
		err := c.ResetLimits()
		logger.MaybeSuccess(err, "request delays and limits were reset")
		logger.MaybeFatal(err, "failed to reset request delays and limits")

	} else if options.MaxPending != nil || options.Delay != nil ||
		options.PximgMaxPending != nil || options.PximgDelay != nil {
		l, warning, err := c.Limits()
		logger.MaybeWarning(warning, "while reading request delays and limits")
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
