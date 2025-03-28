package main

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/config"
	"github.com/fekoneko/piximan/cmd/piximanctl/download"
	"github.com/fekoneko/piximan/cmd/piximanctl/usage"
	"github.com/joho/godotenv"
)

var version string

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load .env: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("piximanctl v%v\n", version)

	args := os.Args[1:]
	if len(args) == 0 {
		usage.RunGeneral()
		return
	}
	command := args[0]
	os.Args = args

	switch command {
	case "config":
		config.Run()
	case "download":
		download.Run()
	default:
		usage.RunGeneral()
	}
}
