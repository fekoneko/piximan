package config

import (
	"time"

	appconfig "github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/termext"
	"github.com/fekoneko/piximan/internal/utils"
)

func config(options *options) {
	termext.DisableInputEcho()
	defer termext.RestoreInputEcho()

	storage, err := appconfig.Open(options.Password)
	logext.MaybeFatal(err, "failed to open config storage")

	if utils.FromPtr(options.ResetSession, false) {
		err := storage.RemoveSessionId()
		logext.MaybeSuccess(err, "session id was removed")
		logext.MaybeFatal(err, "failed to remove session id")
	} else if options.SessionId != nil {
		err = storage.WriteSessionId(*options.SessionId)
		logext.MaybeSuccess(err, "session id was set%v",
			utils.If(options.Password != nil, " and encrypted with password", ""),
		)
		logext.MaybeFatal(err, "failed to set session id")
	}

	if utils.FromPtr(options.ResetConfig, false) {
		err := storage.Reset()
		logext.MaybeSuccess(err, "configuration parameters were reset")
		logext.MaybeFatal(err, "failed to reset configuration parameters")
		return
	}

	changed := false

	if options.PximgMaxPending != nil {
		storage.PximgMaxPending = *options.PximgMaxPending
		changed = true
	}
	if options.PximgDelay != nil {
		storage.PximgDelay = time.Duration(*options.PximgDelay) * time.Second
		changed = true
	}
	if options.DefaultMaxPending != nil {
		storage.DefaultMaxPending = *options.DefaultMaxPending
		changed = true
	}
	if options.DefaultDelay != nil {
		storage.DefaultDelay = time.Duration(*options.DefaultDelay) * time.Second
		changed = true
	}
	if changed {
		err = storage.Write()
		logext.MaybeSuccess(err, "configuration parameters were saved")
		logext.MaybeFatal(err, "failed to save configuration parameters")
	}
}
