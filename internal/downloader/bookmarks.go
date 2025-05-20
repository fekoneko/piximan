package downloader

import (
	"fmt"
	"strings"

	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
)

// Schedule artwork bookmarks for download. Run() to start downloading.
func (d *Downloader) ScheduleBookmarks(
	userId uint64, kind queue.ItemKind, tag *string, size image.Size,
	onlyMeta bool, lowMeta bool, paths []string,
) {
	d.crawlQueueMutex.Lock()
	defer d.crawlQueueMutex.Unlock()

	// TODO: fetch other offsets as well, but in other task in crawlQueue not to block downloading
	//       schedule other crawl fetches after first one is done and total count is known
	//       don't forget to log
	offset := uint(0)
	d.crawlQueue = append(d.crawlQueue, func() error {
		return d.scheduleBookmarks(userId, kind, tag, offset, size, onlyMeta, lowMeta, paths)
	})
	logext.Info(bookmarksLogPrefix("created bookmarks crawl task", userId, tag, &offset))
}

// fetch bookmarks and then schedule them for download
func (d *Downloader) scheduleBookmarks(
	userId uint64, kind queue.ItemKind, tag *string, offset uint,
	size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) error {
	var successMessage = fmt.Sprintf("failed to fetch %v bookmarks", kind)
	var errorMessage = fmt.Sprintf("failed to fetch %v bookmarks", kind)

	if kind != queue.ItemKindArtwork && kind != queue.ItemKindNovel {
		err := fmt.Errorf("invalid work type: %v", uint8(kind))
		logext.Error("%v: %v", bookmarksLogPrefix(errorMessage, userId, tag, nil), err)
		return err
	}

	sessionId, withSessionId := d.sessionId()
	if !withSessionId {
		err := fmt.Errorf("authorization is required")
		logext.Error("%v: %v", bookmarksLogPrefix(errorMessage, userId, tag, nil), err)
		return err
	}

	var results []fetch.BookmarkResult
	var err error
	if kind == queue.ItemKindArtwork {
		results, err = fetch.ArtworkBookmarksAuthorized(d.client(), userId, tag, offset, 100, *sessionId)
	} else {
		results, err = fetch.NovelBookmarksAuthorized(d.client(), userId, tag, offset, 100, *sessionId)
	}
	logext.MaybeSuccess(err, bookmarksLogPrefix(successMessage, userId, tag, &offset))
	logext.MaybeError(err, bookmarksLogPrefix(errorMessage, userId, tag, &offset))
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Work.Id == nil {
			err := fmt.Errorf("work id is missing in %v", result.Work)
			logext.Error("%v %v: %v", bookmarksLogPrefix("failed to schedule", userId, tag, &offset), kind, err)
			continue
		}
		d.ScheduleWithKnown(
			[]uint64{*result.Work.Id}, kind, size, onlyMeta, paths,
			result.Work, result.ImageUrl, lowMeta,
		)
	}

	return nil
}

func bookmarksLogPrefix(message string, userId uint64, tag *string, offset *uint) string {
	builder := strings.Builder{}
	builder.WriteString(message)
	builder.WriteString(fmt.Sprintf(" for user %v", userId))
	if tag != nil && offset == nil {
		builder.WriteString(fmt.Sprintf(" with tag %v", *tag))
	} else if tag == nil && offset != nil {
		builder.WriteString(fmt.Sprintf(" with offset %v", *offset))
	} else if tag != nil && offset != nil {
		builder.WriteString(fmt.Sprintf(" with tag %v and offset %v", *tag, *offset))
	}
	return builder.String()
}
