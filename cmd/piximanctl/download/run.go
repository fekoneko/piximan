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

	if flag.NArg() != 0 {
		flagext.BadUsage("too many arguments")
	}

	if flag.NFlag() == 0 {
		interactive()
	} else {
		nonInteractive(id, kind, size, path, sessionId, password)
	}
}

func interactive() {
	fmt.Println("Interactive mode for download is not yet implemented")
}

func nonInteractive(id *uint64, kind *string, size *uint, path *string, sessionId *string, password *string) {
	if !flagext.Provided("id") {
		flagext.BadUsage("required flag is not set: -id")
	}
	if flagext.Provided("type") && *kind != "artwork" && *kind != "novel" {
		flagext.BadUsage("invalid argument value: -type")
	}
	if flagext.Provided("size") && *size > 3 {
		flagext.BadUsage("invalid argument value: -size")
	}

	if !flagext.Provided("sessionid") {
		// TODO: if -password is not provided,
		//       try with empty string and then ask for password interactively
		sessionId = readSessionId(password)
	}
	download(id, kind, size, path, sessionId)
}

func readSessionId(password *string) *string {
	storage, err := secretstorage.Open(*password)
	if err != nil {
		fmt.Printf("failed to get session id: %v\n", err)
		os.Exit(1)
	}
	if err := storage.Read(); err != nil {
		fmt.Printf("failed to get session id: %v\n", err)
		os.Exit(1)
	}
	if storage.SessionId == nil {
		fmt.Println("no session id is configured or provided")
		os.Exit(1)
	}
	return storage.SessionId
}

func download(id *uint64, kind *string, size *uint, path *string, sessionId *string) {
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
