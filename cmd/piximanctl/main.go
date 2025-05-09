package main

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/config"
	"github.com/fekoneko/piximan/cmd/piximanctl/download"
	"github.com/fekoneko/piximan/cmd/piximanctl/help"
)

var version string

func main() {
	fmt.Printf("piximanctl v%v\n\n", version)

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
		os.Args = os.Args[1:]
	}

	switch command {
	case "config":
		config.Run()
	case "download":
		download.Run()
	case "help":
		help.Run()
	default:
		help.RunGeneral()
	}
}
