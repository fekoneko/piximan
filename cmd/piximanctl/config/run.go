package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/usage"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/storage"
)

func Run() {
	sessionId := *flag.String("sessionid", "", "")
	flag.Usage = usage.RunConfig
	flag.Parse()

	if flagext.Provided("sessionid") {
		if err := storage.StoreSessionId(sessionId); err != nil {
			fmt.Printf("failed to set sessionid: %v\n", err)
			os.Exit(1)
		}
	}
}
