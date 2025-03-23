package download

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/usage"
	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/settings"
)

func Run() {
	if len(os.Args) == 1 {
		usage.Download()
		os.Exit(0)
	}

	id := flag.Uint64("id", 0, "")
	path := flag.String("path", "", "")
	flag.Usage = usage.Download
	flag.Parse()

	if !flagext.Provided("id") {
		fmt.Println("required flag is not set: -id")
		usage.Download()
		os.Exit(2)
	}

	sessionId, err := settings.SessionId()
	if err != nil {
		fmt.Printf("failed to get session id: %v\n", err)
		os.Exit(1)
	}

	d := downloader.New(sessionId)
	err = d.DownloadWork(*id, *path)
	if err != nil {
		fmt.Printf("failed to download work: %v\n", err)
		os.Exit(1)
	}
}
