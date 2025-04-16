package download

import (
	"flag"

	"github.com/fekoneko/piximan/cmd/piximanctl/help"
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/flagext"
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
		size:     flag.Uint("size", uint(image.SizeDefault), ""),
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
		interactive()
	} else {
		nonInteractive(flags)
	}
}
