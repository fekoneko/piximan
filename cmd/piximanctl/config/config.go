package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/usage"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/settings"
)

func Run() {
	sessionId := flag.String("sessionid", "", "")
	flag.Usage = usage.Config
	flag.Parse()

	if flagext.Provided("sessionid") {
		err := settings.SetSessionId(*sessionId)
		if err != nil {
			fmt.Printf("failed to set sessionid: %v\n", err)
			os.Exit(1)
		}
	}
}
