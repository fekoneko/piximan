package main

import (
	"fmt"

	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/joho/godotenv"
)

var version string

func main() {
	err := godotenv.Load()
	logext.MaybeFatal(err, "failed to load .env")

	fmt.Printf("piximan v%v\n", version)
}
