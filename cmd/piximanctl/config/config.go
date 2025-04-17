package config

import (
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
)

func configSessionId(options *options) {
	if len(options.SessionId) == 0 {
		err := secretstorage.RemoveSessionId()
		logext.MaybeFatal(err, "failed to remove session id")
	} else {
		password := ""
		if options.Password != nil {
			password = *options.Password
		}
		storage, err := secretstorage.Open(password)
		logext.MaybeFatal(err, "failed to open session id storage")

		err = storage.StoreSessionId(options.SessionId)
		logext.MaybeFatal(err, "failed to set session id")
	}
}
