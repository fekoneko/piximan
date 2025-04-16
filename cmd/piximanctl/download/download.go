package download

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/pathext"
)

func download(flags flags, isInferIdProvided bool, isPathProvided bool) {
	d := downloader.New()
	size := image.SizeFromUint(*flags.size)
	kind := queue.ItemKindFromString(*flags.kind)

	if isInferIdProvided {
		result, err := pathext.InferIdsFromWorkPath(*flags.inferId)
		if err != nil {
			fmt.Printf("cannot infer work id from pattern %v: %v\n", *flags.inferId, err)
			os.Exit(1)
		}
		q := queue.FromMap(result, kind, size, *flags.onlyMeta)
		if isPathProvided {
			for i := range *q {
				(*q)[i].Paths = []string{*flags.path}
			}
		}
		fmt.Print(q, "\n\n")
		d.ScheduleQueue(q)
	} else {
		paths := []string{*flags.path}
		d.Schedule(*flags.id, kind, size, *flags.onlyMeta, paths)
	}

	for d.Listen() != nil {
	}
}
