package download

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	if args, err := flags.Parse(options); err != nil {
		os.Exit(1)
	} else if len(args) > 0 {
		fmt.Println("extra arguments provided")
		os.Exit(2)
	}

	if !utils.ExactlyOneDefined(
		options.Ids, options.Bookmarks, options.InferIds, options.Lists,
	) {
		fmt.Println("provide exactly one download source: `-i, --id', `-b, --bookmarks' `-I, --infer-id' or `-l, --list'")
		os.Exit(2)
	}
	withBookmarksUserId := utils.ParseUintPtr(options.Bookmarks) == nil
	if options.Bookmarks != nil && *options.Bookmarks != "my" && withBookmarksUserId {
		fmt.Println("invalid argument for flag `-b, --bookmarks'")
		os.Exit(2)
	}
	if options.Kind != nil && !queue.ValidItemKindString(*options.Kind) {
		fmt.Println("invalid argument for flag `-t, --type'")
		os.Exit(2)
	}
	if options.Size != nil && !imageext.ValidSizeUint(*options.Size) {
		fmt.Println("invalid argument for flag `-s, --size'")
		os.Exit(2)
	}
	if options.Language != nil && options.Kind != nil && *options.Kind != queue.ItemKindArtworkString {
		fmt.Println("`-L, --language' flag can only be used with `-t, --type' artwork")
		os.Exit(2)
	}
	if options.Language != nil && !work.ValidLanguageString(*options.Language) {
		fmt.Println("invalid argument for flag `-L, --language'")
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
	if options.UntilSkip != nil && options.Bookmarks == nil {
		fmt.Println("`-U, --until-skip' flag can only be used with `-b, --bookmarks' source")
		os.Exit(2)
	}
	if options.UntilSkip != nil && options.Skips == nil {
		fmt.Println("`-U, --until-skip' flag can only be used when `-S, --skip' was provided")
		os.Exit(2)
	}
	if options.Skips != nil {
		for _, s := range *options.Skips {
			if !fsext.IsInferIdPattern(s) {
				continue
			}
			if err := fsext.InferIdPatternValid(s); err != nil {
				fmt.Printf("invalid argument for flag `-S, --skip': "+
					"invalid infer id pattern %v: %v\n", s, err)
				os.Exit(2)
			}
		}
	}
	if options.Paths != nil {
		for _, s := range *options.Paths {
			if err := fsext.WorkPathPatternValid(s); err != nil {
				fmt.Printf("invalid argument for flag `-p, --path': %v\n", err)
				os.Exit(2)
			}
		}
	}
	if options.InferIds != nil {
		for _, s := range *options.InferIds {
			if !fsext.IsInferIdPattern(s) {
				continue
			}
			if err := fsext.InferIdPatternValid(s); err != nil {
				fmt.Printf("invalid argument for flag `-I, --infer-id': "+
					"invalid infer id pattern %v: %v\n", s, err)
				os.Exit(2)
			}
		}
	}

	download(options)
}
