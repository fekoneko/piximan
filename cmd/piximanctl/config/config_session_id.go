package config

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/pkg/secretstorage"
)

func configSessionId(flags flags) {
	if len(*flags.sessionId) == 0 {
		if err := secretstorage.RemoveSessionId(); err != nil {
			fmt.Printf("failed to remove session id: %v\n", err)
			os.Exit(1)
		}
	} else {
		storage, err := secretstorage.Open(*flags.password)
		if err != nil {
			fmt.Printf("failed to open session id storage: %v\n", err)
			os.Exit(1)
		}
		if err := storage.StoreSessionId(*flags.sessionId); err != nil {
			fmt.Printf("failed to set session id: %v\n", err)
			os.Exit(1)
		}
	}
}
