package queue

import "github.com/fekoneko/piximan/internal/downloader/image"

type Item struct {
	Id       uint64
	Kind     ItemKind
	Size     image.Size
	OnlyMeta bool
	Paths    []string
}
