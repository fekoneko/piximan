package config

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/pkg/secretstorage"
)

func configSessionId(options *options) {
	if len(options.SessionId) == 0 {
		if err := secretstorage.RemoveSessionId(); err != nil {
			fmt.Printf("failed to remove session id: %v\n", err)
			os.Exit(1)
		}
	} else {
		password := ""
		if options.Password != nil {
			password = *options.Password
		}
		storage, err := secretstorage.Open(password)
		if err != nil {
			fmt.Printf("failed to open session id storage: %v\n", err)
			os.Exit(1)
		}
		if err := storage.StoreSessionId(options.SessionId); err != nil {
			fmt.Printf("failed to set session id: %v\n", err)
			os.Exit(1)
		}
	}
}
