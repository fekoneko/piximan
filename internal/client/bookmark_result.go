package client

import (
	"time"

	"github.com/fekoneko/piximan/internal/work"
)

type BookmarkResult struct {
	Work           *work.Work
	BookmarkedTime *time.Time
	ImageUrl       *string
}
