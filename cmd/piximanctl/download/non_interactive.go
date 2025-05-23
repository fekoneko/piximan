package download

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	if _, err := flags.Parse(options); err != nil {
		os.Exit(2)
	}

	if options.Ids == nil && options.InferIdPath == nil && options.QueuePath == nil {
		fmt.Println("one of these flags is not provided: `-i, --id', `-I, --inferid' or `-l, --list'")
		os.Exit(2)
	}
	if options.Ids != nil && options.InferIdPath != nil {
		fmt.Println("providing these flags together is not supported: `-i, --id' and `-I, --inferid'")
		os.Exit(2)
	}
	if options.QueuePath != nil && options.Ids != nil {
		fmt.Println("providing these flags together is not supported: `-l, --list' and `-i, --id'")
		os.Exit(2)
	}
	if options.QueuePath != nil && options.InferIdPath != nil {
		fmt.Println("providing these flags together is not supported: `-l, --list' and `-I, --inferid'")
		os.Exit(2)
	}
	if options.Kind != nil && options.Size != nil && *options.Kind == queue.ItemKindNovelString {
		fmt.Println("cannot use `-s, --size' flag with `-t, --type' novel")
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

	download(options)
}
