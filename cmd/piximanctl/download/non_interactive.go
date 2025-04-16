package download

import (
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/flagext"
)

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
	download(flags, flagext.Provided("inferid"), flagext.Provided("path"))
}
