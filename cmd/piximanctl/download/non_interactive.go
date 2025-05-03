package download

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/jessevdk/go-flags"
)

func nonInteractive() {
	options := &options{}
	if _, err := flags.Parse(options); err != nil {
		os.Exit(2)
	}

	if options.Ids == nil && options.InferId == nil {
		fmt.Println("one of these flags is not provided: `-i, --id` and `-I, --inferid`")
		os.Exit(2)
	}
	if options.Ids != nil && options.InferId != nil {
		fmt.Println("providing these flags together is not supported: `-i, --id` and `-I, --inferid`")
		os.Exit(2)
	}
	if options.King != nil && options.Size != nil && *options.King == queue.ItemKindNovelString {
		fmt.Println("cannot use `-s, --size` flag with `-t, --type` novel")
		os.Exit(2)
	}
	if options.King != nil && !queue.ValidItemKindString(*options.King) {
		fmt.Println("invalid argument for flag `-t, --type`")
		os.Exit(2)
	}
	if options.Size != nil && *options.Size > 3 {
		fmt.Println("invalid argument for flag `-s, --size`")
		os.Exit(2)
	}

	download(options)
}
