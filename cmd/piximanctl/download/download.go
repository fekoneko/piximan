package download

import (
	"fmt"
	"os"

	"github.com/fekoneko/piximan/pkg/downloader"
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
	"github.com/fekoneko/piximan/pkg/pathext"
)

func download(options *options) {
	size := image.SizeDefault
	if options.Size != nil {
		size = image.SizeFromUint(*options.Size)
	}
	kind := queue.ItemKindDefault
	if options.King != nil {
		kind = queue.ItemKindFromString(*options.King)
	}
	onlyMeta := false
	if options.OnlyMeta != nil {
		onlyMeta = *options.OnlyMeta
	}
	path := ""
	if options.Path != nil {
		path = *options.Path
	}

	d := downloader.New()

	if options.InferId != nil {
		result, err := pathext.InferIdsFromWorkPath(*options.InferId)
		if err != nil {
			fmt.Printf("cannot infer work id from pattern %v: %v\n", *options.InferId, err)
			os.Exit(1)
		}
		q := queue.FromMap(result, kind, size, onlyMeta)
		if options.Path != nil {
			for i := range *q {
				(*q)[i].Paths = []string{path}
			}
		}
		fmt.Print(q, "\n\n")
		d.ScheduleQueue(q)
	} else {
		d.Schedule(*options.Id, kind, size, onlyMeta, []string{path})
	}

	for d.Listen() != nil {
	}
}
