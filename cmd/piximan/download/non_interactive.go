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
	if options.Tags != nil && options.Bookmarks == nil {
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
	if options.Fresh != nil && options.Skip == nil {
		fmt.Println("`-f, --fresh' flag can only be used when `-S, --skip' was provided")
		os.Exit(2)
	}
	if options.Skip != nil && fsext.CanBeInferIdPath(*options.Skip) {
		if err := fsext.InferIdPathValid(*options.Skip); err != nil {
			fmt.Printf("invalid argument for flag `-S, --skip': "+
				"infer id pattern found but it's invalid: %v\n", err)
			os.Exit(2)
		}
	}
	if options.Path == nil && options.InferId == nil {
		fmt.Println("`-p, --path' flag is required unless `-I, --infer-id' is provided")
		os.Exit(2)
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
