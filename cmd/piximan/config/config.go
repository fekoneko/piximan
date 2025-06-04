package config

import (
	config1 "github.com/fekoneko/piximan/internal/config"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/termext"
)

func config(options *options) {
	termext.DisableInputEcho()
	defer termext.RestoreInputEcho()

	if len(options.SessionId) == 0 {
		err := config1.RemoveSessionId()
		logext.MaybeFatal(err, "failed to remove session id")
	} else {
		password := ""
		if options.Password != nil {
			password = *options.Password
		}
		storage, err := config1.Open(password)
		logext.MaybeFatal(err, "failed to open session id storage")

		err = storage.WriteSessionId(options.SessionId)
		logext.MaybeFatal(err, "failed to set session id")
	}
}
