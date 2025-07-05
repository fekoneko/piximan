package client

import (
	"github.com/fekoneko/piximan/internal/collection/work"
)

type BookmarkResult struct {
	Work     *work.Work
	ImageUrl *string
}
