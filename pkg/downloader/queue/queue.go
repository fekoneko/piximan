package queue

import (
	"fmt"
	"strings"
)

type Queue []Item

func (q *Queue) String() string {
	fmt.Println("download queue:")
	builder := strings.Builder{}

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

func FromMap(m *map[uint64][]string, kind ItemKind) *Queue {
	queue := make(Queue, len(*m))

	i := 0
	for id, paths := range *m {
		queue[i] = Item{id, kind, paths}
		i++
	}

	return &queue
}
