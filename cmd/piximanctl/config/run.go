package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
)

func Run() {
	sessionId := flag.String("sessionid", "", "")
	password := flag.String("password", "", "")
	flag.Usage = help.RunConfig
	flag.Parse()

	if flag.NArg() != 0 {
		flagext.BadUsage("too many arguments")
	}

	if flagext.Provided("sessionid") {
		storage, err := secretstorage.New(*password)
		if err != nil {
			fmt.Printf("failed to set session id: %v\n", err)
			os.Exit(1)
		}
		if err := storage.StoreSessionId(*sessionId); err != nil {
			fmt.Printf("failed to set sessionid: %v\n", err)
			os.Exit(1)
		}
	}
}
