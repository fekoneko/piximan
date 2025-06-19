package client

import (
	"github.com/fekoneko/piximan/internal/work"
)

type BookmarkResult struct {
	Work     *work.Work
	ImageUrl *string
}
