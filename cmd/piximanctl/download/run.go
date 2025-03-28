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
	kind := flag.String("type", "artwork", "")
	size := flag.Uint("size", uint(downloader.ImageSizeDefault), "")
	path := flag.String("path", "", "")
	sessionId := flag.String("sessionid", "", "")
	flag.Usage = usage.RunDownload
	flag.Parse()

	if len(flag.Args()) != 0 {
		fmt.Println("too many arguments")
		usage.RunDownload()
		os.Exit(2)
	}

	if !flagext.Provided("id") {
		fmt.Println("required flag is not set: -id")
		usage.RunDownload()
		os.Exit(2)
	}

	if flagext.Provided("type") && *kind != "artwork" && *kind != "novel" {
		fmt.Println("invalid argument value: -type")
		usage.RunDownload()
		os.Exit(2)
	}

	if flagext.Provided("size") && *size > 3 {
		fmt.Println("invalid argument value: -size")
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
	var err error
	if *kind == "novel" {
		_, err = d.DownloadNovel(*id, *path)
	} else {
		_, err = d.DownloadArtwork(*id, downloader.ImageSize(*size), *path)
	}
	if err != nil {
		os.Exit(1)
	}
}
