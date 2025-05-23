package download

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	if _, err := flags.Parse(options); err != nil {
		os.Exit(2)
	}

	if !utils.ExactlyOneDefined(
		options.Ids, options.Bookmarks, options.InferIdPath, options.QueuePath,
	) {
		fmt.Println("provide exactly one download source: `-i, --id', `-b, --bookmarks' `-I, --inferid' or `-l, --list'")
		os.Exit(2)
	}
	withBookmarksUserId := utils.ParseUint64Ptr(options.Bookmarks) == nil
	if options.Bookmarks != nil && *options.Bookmarks != "my" && withBookmarksUserId {
		fmt.Println("invalid argument for flag `-b, --bookmarks'")
		os.Exit(2)
	}
	if options.Kind != nil && !queue.ValidItemKindString(*options.Kind) {
		fmt.Println("invalid argument for flag `-t, --type'")
		os.Exit(2)
	}
	if options.Kind != nil && options.Size != nil && *options.Kind == queue.ItemKindNovelString {
		fmt.Println("cannot use `-s, --size' flag with `-t, --type' novel")
		os.Exit(2)
	}
	if options.Size != nil && *options.Size > 3 {
		fmt.Println("invalid argument for flag `-s, --size'")
		os.Exit(2)
	}
	if options.FromOffset != nil && options.Bookmarks == nil {
		fmt.Println("`-T, --from' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.ToOffset != nil && options.Bookmarks == nil {
		fmt.Println("`-T, --to' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.FromOffset != nil && options.ToOffset != nil &&
		*options.FromOffset >= *options.ToOffset {
		fmt.Println("argument `-T, --from' must be less than `-T, --to'")
		os.Exit(2)
	}
	if options.LowMeta != nil && options.Bookmarks == nil {
		fmt.Println("`-M, --lowmeta' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.LowMeta != nil && options.Kind != nil && *options.Kind == queue.ItemKindNovelString &&
		(options.OnlyMeta == nil || !*options.OnlyMeta) {
		fmt.Println("`-M, --lowmeta' can be removed for novels without `-m, --onlymeta'")
		os.Exit(2)
	}

	download(options)
}
