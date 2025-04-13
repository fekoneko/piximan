package download

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/flagext"
	"github.com/fekoneko/piximan/pkg/pathext"
)

type flags struct {
	id       *uint64
	kind     *string
	size     *uint
	path     *string
	inferId  *string
	onlyMeta *bool
}

func Run() {
	flags := flags{
		id:       flag.Uint64("id", 0, ""),
		kind:     flag.String("type", "artwork", ""),
		size:     flag.Uint("size", uint(downloader.ImageSizeDefault), ""),
		path:     flag.String("path", "", ""),
		inferId:  flag.String("inferid", "", ""),
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
	fmt.Println("interactive mode for download is not yet implemented")
	fmt.Println(flags)
	os.Exit(1)
}

func nonInteractive(flags flags) {
	if !flagext.Provided("id") && !flagext.Provided("inferid") {
		flagext.BadUsage("one of these arguments is not provided: -id, -inferid")
	}
	if flagext.Provided("id") && flagext.Provided("inferid") {
		flagext.BadUsage("providing these arguments together is not supporded: -id, -inferid")
	}
	if flagext.Provided("type") && flagext.Provided("size") && *flags.kind == queue.ItemKindNovelString {
		flagext.BadUsage("cannot use -size argument with -type novel")
	}
	if flagext.Provided("type") && queue.ValidItemKindString(*flags.kind) {
		flagext.BadUsage("invalid argument value: -type")
	}
	if flagext.Provided("size") && *flags.size > 3 {
		flagext.BadUsage("invalid argument value: -size")
	}
	continueDownload(flags)
}

func continueDownload(flags flags) {
	d := downloader.New()

	if flagext.Provided("inferid") {
		result, err := pathext.InferIdsFormWorkPath(*flags.inferId)
		if err != nil {
			fmt.Printf("cannot infer work id from pattern %v: %v\n", *flags.inferId, err)
			os.Exit(1)
		}
		queue := queue.FromMap(result, queue.ItemKindFromString(*flags.kind))
		fmt.Println(queue)

		// TODO: implement download queue for this
		fmt.Println("\ndownloading multiple works is not yet implemented")
		os.Exit(1)
	}

	var err error
	if *flags.kind == queue.ItemKindNovelString && *flags.onlyMeta {
		_, err = d.DownloadNovelMeta(*flags.id, *flags.path)
	} else if *flags.kind == queue.ItemKindNovelString {
		_, err = d.DownloadNovel(*flags.id, *flags.path)
	} else if *flags.kind == queue.ItemKindArtworkString && *flags.onlyMeta {
		_, err = d.DownloadArtworkMeta(*flags.id, *flags.path)
	} else if *flags.kind == queue.ItemKindArtworkString {
		_, err = d.DownloadArtwork(*flags.id, *flags.path, downloader.ImageSize(*flags.size))
	}
	if err != nil {
		os.Exit(1)
	}
}
