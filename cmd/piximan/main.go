package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var version string

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load .env: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("piximan v%v\n", version)
}
