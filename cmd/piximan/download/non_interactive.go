package download

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	if _, err := flags.Parse(options); err != nil {
		os.Exit(2)
	}

	if !utils.ExactlyOneDefined(
		options.Ids, options.Bookmarks, options.InferId, options.List,
	) {
		fmt.Println("provide exactly one download source: `-i, --id', `-b, --bookmarks' `-I, --infer-id' or `-l, --list'")
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
	if options.Size != nil && *options.Size > 3 {
		fmt.Println("invalid argument for flag `-s, --size'")
		os.Exit(2)
	}
	if options.Tag != nil && options.Bookmarks == nil {
		fmt.Println("`-G, --tag' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.FromOffset != nil && options.Bookmarks == nil {
		fmt.Println("`-F, --from' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.ToOffset != nil && options.Bookmarks == nil {
		fmt.Println("`-T, --to' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.FromOffset != nil && options.ToOffset != nil &&
		*options.FromOffset >= *options.ToOffset {
		fmt.Println("argument `-F, --from' must be less than `-T, --to'")
		os.Exit(2)
	}
	if options.LowMeta != nil && options.Bookmarks == nil {
		fmt.Println("`-M, --low-meta' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.LowMeta != nil && options.Kind != nil && *options.Kind == queue.ItemKindNovelString &&
		(options.OnlyMeta == nil || !*options.OnlyMeta) {
		fmt.Println("`-M, --low-meta' can be removed for novels without `-m, --only-meta'")
		os.Exit(2)
	}
	if options.Private != nil && options.Bookmarks == nil {
		fmt.Println("`-R, --private' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.Fresh != nil && options.Bookmarks == nil {
		fmt.Println("`-f, --fresh' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.Fresh != nil && options.Collection == nil {
		fmt.Println("`-f, --fresh' flag can only be used when `-c, --collection' was provided")
		os.Exit(2)
	}
	if options.Collection != nil && fsext.CanBeInferIdPath(*options.Collection) {
		if err := fsext.InferIdPathValid(*options.Collection); err != nil {
			fmt.Printf("invalid argument for flag `-c, --collection': "+
				"infer id pattern found but it's invalid: %v\n", err)
			os.Exit(2)
		}
	}
	if options.Path != nil {
		if err := fsext.WorkPathValid(*options.Path); err != nil {
			fmt.Printf("invalid argument for flag `-p, --path': %v\n", err)
			os.Exit(2)
		}
	}
	if options.InferId != nil {
		if err := fsext.InferIdPathValid(*options.InferId); err != nil {
			fmt.Printf("invalid argument for flag `-I, --infer-id': %v\n", err)
			os.Exit(2)
		}
	}

	download(options)
}
