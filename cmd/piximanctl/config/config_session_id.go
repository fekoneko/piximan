package config

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/pkg/secretstorage"
)

func configSessionId(flags flags) {
	storage, err := secretstorage.Open(*flags.password)
	if err != nil {
		fmt.Printf("failed to set session id: %v\n", err)
		os.Exit(1)
	}
	if len(*flags.sessionId) == 0 {
		if err := storage.RemoveSessionId(); err != nil {
			fmt.Printf("failed to remove session id: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := storage.StoreSessionId(*flags.sessionId); err != nil {
			fmt.Printf("failed to set session id: %v\n", err)
			os.Exit(1)
		}
	}
}
