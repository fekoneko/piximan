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
		if kind == queue.ItemKindArtwork {
			return d.scheduleArtworkBookmarks(userId, tag, size, onlyMeta, lowMeta, paths)
		} else if kind == queue.ItemKindNovel {
			return d.scheduleNovelBookmarks(userId, tag, onlyMeta, lowMeta, paths)
		} else {
			err := fmt.Errorf("invalid work type: %v", uint8(kind))
			logext.Error("%v: %v", bookmarksLogPrefix("failed to fetch bookmarks", userId, tag, nil), err)
			return err
		}
	})
	logext.Info(bookmarksLogPrefix("created bookmarks crawl task", userId, tag, &offset))
}

// Fetch artwork bookmarks and then schedule them for download, blocks until done.
func (d *Downloader) scheduleArtworkBookmarks(
	userId uint64, tag *string, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) error {
	sessionId, withSessionId := d.sessionId()
	if !withSessionId {
		err := fmt.Errorf("authorization is required")
		logext.Error("%v: %v", bookmarksLogPrefix("failed to fetch artwork bookmarks", userId, tag, nil), err)
		return err
	}

	offset := uint(0)
	results, err := fetch.ArtworkBookmarksAuthorized(d.client(), userId, tag, offset, 100, *sessionId)
	logext.MaybeSuccess(err, bookmarksLogPrefix("fetched artwork bookmarks", userId, tag, &offset))
	logext.MaybeError(err, bookmarksLogPrefix("failed to fetch artwork bookmarks", userId, tag, &offset))
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Work.Id == nil {
			err := fmt.Errorf("work id is missing in %v", result.Work)
			logext.Error("%v: %v", bookmarksLogPrefix("failed to schedule artwork bookmark", userId, tag, &offset), err)
			continue
		}
		d.ScheduleWithKnown(
			[]uint64{*result.Work.Id}, queue.ItemKindArtwork, size, onlyMeta, paths,
			result.Work, result.ThumbnailUrl, lowMeta,
		)
	}

	return nil
}

// Fetch novel bookmarks and then schedule them for download, blocks until done.
func (d *Downloader) scheduleNovelBookmarks(
	userId uint64, tag *string, onlyMeta bool, lowMeta bool, paths []string,
) error {
	sessionId, withSessionId := d.sessionId()
	if !withSessionId {
		err := fmt.Errorf("authorization is required")
		logext.Error("%v: %v", bookmarksLogPrefix("failed to fetch novel bookmarks", userId, tag, nil), err)
		return err
	}

	offset := uint(0)
	results, err := fetch.NovelBookmarksAuthorized(d.client(), userId, tag, offset, 100, *sessionId)
	logext.MaybeSuccess(err, bookmarksLogPrefix("fetched novel bookmarks", userId, tag, &offset))
	logext.MaybeError(err, bookmarksLogPrefix("failed to fetch novel bookmarks", userId, tag, &offset))
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Work.Id == nil {
			err := fmt.Errorf("work id is missing in %v", result.Work)
			logext.Error("%v: %v", bookmarksLogPrefix("failed to schedule novel bookmark", userId, tag, &offset), err)
			continue
		}
		d.ScheduleWithKnown(
			[]uint64{*result.Work.Id}, queue.ItemKindNovel, image.SizeDefault, onlyMeta, paths,
			result.Work, result.CoverUrl, lowMeta,
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
