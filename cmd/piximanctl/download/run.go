package download

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/flagext"
)

type flags struct {
	id       *uint64
	kind     *string
	size     *uint
	path     *string
	onlyMeta *bool
}

func Run() {
	flags := flags{
		id:       flag.Uint64("id", 0, ""),
		kind:     flag.String("type", "artwork", ""),
		size:     flag.Uint("size", uint(downloader.ImageSizeDefault), ""),
		path:     flag.String("path", "", ""),
		onlyMeta: flag.Bool("onlymeta", false, ""),
	}
	flag.Usage = help.RunDownload
	flag.Parse()

	if flag.NArg() != 0 {
		flagext.BadUsage("too many arguments")
	}

	if flag.NFlag() == 0 {
		interactive(flags)
	} else {
		nonInteractive(flags)
	}
}

func interactive(flags flags) {
	// TODO: implement interactive mode
	fmt.Println("Interactive mode for download is not yet implemented")
	fmt.Println(flags)
	os.Exit(1)
}

func nonInteractive(flags flags) {
	if !flagext.Provided("id") {
		flagext.BadUsage("required flag is not set: -id")
	}
	if flagext.Provided("type") && *flags.kind != "artwork" && *flags.kind != "novel" {
		flagext.BadUsage("invalid argument value: -type")
	}
	if flagext.Provided("size") && *flags.size > 3 {
		flagext.BadUsage("invalid argument value: -size")
	}
	download(flags)
}

func download(flags flags) {
	d := downloader.New()

	var err error
	if *flags.kind == "novel" && *flags.onlyMeta {
		_, err = d.DownloadNovelMeta(*flags.id, *flags.path)
	} else if *flags.kind == "novel" {
		_, err = d.DownloadNovel(*flags.id, *flags.path)
	} else if *flags.kind == "artwork" && *flags.onlyMeta {
		_, err = d.DownloadArtworkMeta(*flags.id, *flags.path)
	} else if *flags.kind == "artwork" {
		_, err = d.DownloadArtwork(*flags.id, downloader.ImageSize(*flags.size), *flags.path)
	}
	if err != nil {
		os.Exit(1)
	}
}
