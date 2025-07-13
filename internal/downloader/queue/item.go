package queue

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
)

type Item struct {
	Id       uint64
	Kind     ItemKind
	Size     imageext.Size
	OnlyMeta bool
	Paths    []string

	// Partial work metadata if available, may be used to reduce the number of requests
	Work *work.Work

	// Thumbnail / cover url if available, may be used to reduce the number of requests
	ImageUrl *string

	// Whether to download full metadata or store partial metadata available in Work
	LowMeta bool
}
