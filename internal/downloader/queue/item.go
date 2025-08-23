package queue

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/imageext"
)

type Item struct {
	Id       uint64
	Kind     ItemKind
	Size     imageext.Size
	Language work.Language
	OnlyMeta bool
	Paths    []string

	Work     *work.Work // Partial work metadata if available, may be used to reduce the number of requests
	ImageUrl *string    // Thumbnail / cover url, nesessary when Work is provided
	LowMeta  bool       // Whether to download full metadata or store partial metadata available in Work
}
