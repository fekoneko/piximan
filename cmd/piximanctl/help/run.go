package help

import (
	"os"
)

func Run() {
	var section string
	if len(os.Args) > 1 {
		section = os.Args[1]
	}

	switch section {
	case "config":
		RunConfig()
	case "download":
		RunDownload()
	default:
		RunGeneral()
	}
}
