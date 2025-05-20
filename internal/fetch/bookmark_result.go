package fetch

import (
	"time"

	"github.com/fekoneko/piximan/internal/collection/work"
)

type BookmarkResult struct {
	Work           *work.Work
	BookmarkedTime *time.Time
	ImageUrl       *string
}
