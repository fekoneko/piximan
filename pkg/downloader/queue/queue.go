package queue

import (
	"fmt"
	"slices"
	"strings"

	"github.com/fekoneko/piximan/pkg/downloader/image"
)

type Queue []Item

func (q *Queue) Push(items ...Item) {
	for _, item := range items {
		if len(item.Paths) == 0 {
			continue
		}

		existing := slices.IndexFunc(*q, func(queueItem Item) bool {
			return item.Id == queueItem.Id &&
				item.Kind == queueItem.Kind &&
				item.Size == queueItem.Size &&
				item.OnlyMeta == queueItem.OnlyMeta
		})

		if existing == -1 {
			*q = append(*q, item)
			continue
		}

		for _, path := range item.Paths {
			if !slices.Contains((*q)[existing].Paths, path) {
				(*q)[existing].Paths = append((*q)[existing].Paths, path)
			}
		}
	}
}

func (q *Queue) Pop() *Item {
	if len(*q) == 0 {
		return nil
	}
	item := &(*q)[0]
	*q = (*q)[1:]

	return item
}

func (q *Queue) String() string {
	if len(*q) == 0 {
		return "empty download queue\n"
	}

	builder := strings.Builder{}
	builder.WriteString("download queue:\n")

	for i, item := range *q {
		if i >= 10 {
			line := fmt.Sprintf("... and %v more\n", len(*q)-i)
			builder.WriteString(line)
			break
		}
		kind := item.Kind.String()
		line := fmt.Sprintf("- id: %-10v type: %-7v paths: %v\n", item.Id, kind, item.Paths)
		builder.WriteString(line)
	}

	return builder.String()
}

func FromMap(m *map[uint64][]string, kind ItemKind, size image.Size, onlyMeta bool) *Queue {
	queue := make(Queue, len(*m))

	i := 0
	for id, paths := range *m {
		queue[i] = Item{id, kind, size, onlyMeta, paths}
		i++
	}

	return &queue
}
