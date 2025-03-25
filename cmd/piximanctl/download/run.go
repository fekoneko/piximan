package download

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/usage"
	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/storage"
)

func Run() {
	if len(os.Args) == 1 {
		usage.RunDownload()
		os.Exit(0)
	}

	id := flag.Uint64("id", 0, "")
	path := flag.String("path", "", "")
	sessionId := flag.String("sessionid", "", "")
	flag.Usage = usage.RunDownload
	flag.Parse()

	if !flagext.Provided("id") {
		fmt.Println("required flag is not set: -id")
		usage.RunDownload()
		os.Exit(2)
	}

	if !flagext.Provided("sessionid") {
		var err error
		*sessionId, err = storage.StoredSessionId()
		if err != nil {
			fmt.Printf("failed to get session id: %v\n", err)
			os.Exit(1)
		}
	}

	d := downloader.New(*sessionId)
	// TODO: add size flag
	if _, err := d.DownloadWork(*id, downloader.ImageSizeDefault, *path); err != nil {
		fmt.Printf("failed to download work: %v\n", err)
		os.Exit(1)
	}
}
