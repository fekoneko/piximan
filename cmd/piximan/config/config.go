package config

import (
	appconfig "github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/termext"
)

func config(options *options) {
	termext.DisableInputEcho()
	defer termext.RestoreInputEcho()

	if len(options.SessionId) == 0 {
		err := appconfig.RemoveSessionId()
		logext.MaybeFatal(err, "failed to remove session id")
	} else {
		storage, err := appconfig.Open(options.Password)
		logext.MaybeFatal(err, "failed to open config storage")

		err = storage.WriteSessionId(options.SessionId)
		logext.MaybeFatal(err, "failed to set session id")
	}
}
