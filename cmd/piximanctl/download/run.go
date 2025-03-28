package download

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/secretstorage"
)

func Run() {
	if len(os.Args) == 1 {
		help.RunDownload()
		os.Exit(0)
	}

	id := flag.Uint64("id", 0, "")
	kind := flag.String("type", "artwork", "")
	size := flag.Uint("size", uint(downloader.ImageSizeDefault), "")
	path := flag.String("path", "", "")
	sessionId := flag.String("sessionid", "", "")
	password := flag.String("password", "", "")
	flag.Usage = help.RunDownload
	flag.Parse()

	if len(flag.Args()) != 0 {
		fmt.Println("too many arguments")
		help.RunDownload()
		os.Exit(2)
	}

	if !flagext.Provided("id") {
		fmt.Println("required flag is not set: -id")
		help.RunDownload()
		os.Exit(2)
	}

	if flagext.Provided("type") && *kind != "artwork" && *kind != "novel" {
		fmt.Println("invalid argument value: -type")
		help.RunDownload()
		os.Exit(2)
	}

	if flagext.Provided("size") && *size > 3 {
		fmt.Println("invalid argument value: -size")
		help.RunDownload()
		os.Exit(2)
	}

	if !flagext.Provided("sessionid") {
		storage, err := secretstorage.New(*password)
		if err != nil {
			fmt.Printf("failed to get session id: %v\n", err)
			os.Exit(1)
		}
		if storage.SessionId == nil {
			fmt.Println("no session id is configured or provided")
			os.Exit(1)
		}
		sessionId = storage.SessionId
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
