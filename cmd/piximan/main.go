package main

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximan/app"
	"github.com/fekoneko/piximan/cmd/piximan/config"
	"github.com/fekoneko/piximan/cmd/piximan/download"
	"github.com/fekoneko/piximan/cmd/piximan/help"
)

var version string

func main() {
	fmt.Printf("piximan %v\n\n", version)

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
		os.Args = os.Args[1:]
	}

	switch command {
	case "":
		app.Run()
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
