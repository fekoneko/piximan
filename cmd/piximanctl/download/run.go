package download

import (
	"flag"
	"fmt"
	"os"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/downloader"
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
		flagext.BadUsage("one of these flags is not provided: -id, -inferid")
	}
	if flagext.Provided("id") && flagext.Provided("inferid") {
		flagext.BadUsage("providing these flags together is not supporded: -id, -inferid")
	}
	if flagext.Provided("type") && *flags.kind != "artwork" && *flags.kind != "novel" {
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
			fmt.Printf("cannot infer work id from path %v: %v\n", *flags.path, err)
			os.Exit(1)
		}

		count := 0
		fmt.Println("ids inferred from the pattern:")
		for id, paths := range result {
			if count >= 10 {
				fmt.Printf("... and %v more\n", len(result)-count)
				break
			}
			fmt.Printf("- id: %-10v paths: %v\n", id, paths)
			count++
		}

		// TODO: implement download queue for this
		fmt.Println("\ndownloading multiple works is not yet implemented")
		os.Exit(1)
	}

	var err error
	if *flags.kind == "novel" && *flags.onlyMeta {
		_, err = d.DownloadNovelMeta(*flags.id, *flags.path)
	} else if *flags.kind == "novel" {
		_, err = d.DownloadNovel(*flags.id, *flags.path)
	} else if *flags.kind == "artwork" && *flags.onlyMeta {
		_, err = d.DownloadArtworkMeta(*flags.id, *flags.path)
	} else if *flags.kind == "artwork" {
		_, err = d.DownloadArtwork(*flags.id, *flags.path, downloader.ImageSize(*flags.size))
	}
	if err != nil {
		os.Exit(1)
	}
}
