package downloader

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/utils"
)

const BOOKMARK_PAGE_SIZE = 100

// TODO: private bookmarks
// Schedule bookmarks of authorized user for download. Run() to start downloading.
func (d *Downloader) ScheduleMyBookmarks(
	kind queue.ItemKind, tag *string, from *uint64, to *uint64, newerTHan *time.Time, olderThan *time.Time,
	private bool, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) {
	d.crawlQueueMutex.Lock()
	defer d.crawlQueueMutex.Unlock()

	d.crawlQueue = append(d.crawlQueue, func() error {
		userId, err := d.client.MyIdAutorized()
		logext.MaybeSuccess(err, "fetched authorizeed user id")
		logext.MaybeError(err, "failed to fetch authorizeed user id")
		if err != nil {
			return err
		}

		d.ScheduleBookmarks(
			userId, kind, tag, from, to, newerTHan, olderThan,
			private, size, onlyMeta, lowMeta, paths,
		)
		return nil
	})
	logext.Info("created crawl task to fetch authorizeed user id")
}

// Schedule bookmarks for download. Run() to start downloading.
// First page of bookmarks is fetched to determine the total count of bookmarks, then the task for
// crawling each page is created, so that it can be run in parallel.
// If older or newer time boundary is specified, all the work will be done sequentially in one crawl task.
func (d *Downloader) ScheduleBookmarks(
	userId uint64, kind queue.ItemKind, tag *string, from *uint64, to *uint64, newerThan *time.Time,
	olderThan *time.Time, private bool, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) {
	if newerThan != nil && olderThan != nil && newerThan.After(*olderThan) {
		err := fmt.Errorf(
			"newer time boundary (%v) represents the time after older time boundary (%v)",
			newerThan, olderThan,
		)
		logext.Error("%v: %v", bookmarksLogMessage("failed to schedule bookmarks", userId, tag, from), err)
		return
	}

	d.crawlQueueMutex.Lock()
	defer d.crawlQueueMutex.Unlock()

	d.crawlQueue = append(d.crawlQueue, func() error {
		fromOffset, toOffset := utils.FromPtr(from, 0), utils.FromPtr(to, math.MaxUint64)
		limit := min(BOOKMARK_PAGE_SIZE, toOffset-fromOffset)
		total, hasOlderTime, hasNewerTime, err := d.fetchBookmarksAndScheduleWorks(
			userId, kind, tag, fromOffset, limit, newerThan, olderThan,
			private, size, onlyMeta, lowMeta, paths,
		)
		if err != nil {
			return err
		}
		fromOffset += BOOKMARK_PAGE_SIZE
		toOffset = min(toOffset, total)

		if newerThan == nil && olderThan == nil {
			d.processBookmarksParallel(
				userId, kind, tag, fromOffset, toOffset, newerThan,
				olderThan, private, size, onlyMeta, lowMeta, paths,
			)
		} else if newerThan == nil && olderThan != nil {
			if !hasNewerTime {
				d.processBookmarksSeqReverse(
					userId, kind, tag, fromOffset, toOffset, newerThan,
					olderThan, private, size, onlyMeta, lowMeta, paths,
				)
			}
		} else {
			if !hasOlderTime {
				d.processBookmarksSeq(
					userId, kind, tag, fromOffset, toOffset, newerThan,
					olderThan, private, size, onlyMeta, lowMeta, paths,
				)
			}
		}
		return nil
	})
	logext.Info(bookmarksLogMessage("created bookmarks crawl task", userId, tag, nil))
}

// Create bookmark crawl tasks for all bookmark pages in range in parallel
func (d *Downloader) processBookmarksParallel(
	userId uint64, kind queue.ItemKind, tag *string, fromOffset uint64, toOffset uint64, olderTime *time.Time,
	newerTime *time.Time, private bool, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) {
	d.crawlQueueMutex.Lock()
	defer d.crawlQueueMutex.Unlock()

	offset, numTasks := fromOffset, 0
	for offset < toOffset {
		limit := min(BOOKMARK_PAGE_SIZE, toOffset-offset)
		d.crawlQueue = append(d.crawlQueue, func() error {
			_, _, _, err := d.fetchBookmarksAndScheduleWorks(
				userId, kind, tag, offset, limit, olderTime, newerTime,
				private, size, onlyMeta, lowMeta, paths,
			)
			return err
		})
		numTasks++

		if offset+limit > offset {
			offset += limit
		} else {
			break
		}
	}

	if numTasks > 0 {
		logext.Info(
			bookmarksLogMessage("created %v bookmarks crawl %v", userId, tag, nil),
			numTasks, utils.If(numTasks == 1, "task", "tasks"),
		)
	}
}

// Fetch bookmark pages sequentially and schedule the works for download.
// Stop when the bookmarks older than the specified time were found.
func (d *Downloader) processBookmarksSeq(
	userId uint64, kind queue.ItemKind, tag *string, fromOffset uint64, toOffset uint64, newerThan *time.Time,
	olderThan *time.Time, private bool, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) error {
	hasOlderTime, err := false, error(nil)
	for fromOffset < toOffset && !hasOlderTime {
		limit := min(BOOKMARK_PAGE_SIZE, toOffset-fromOffset)
		if _, hasOlderTime, _, err = d.fetchBookmarksAndScheduleWorks(
			userId, kind, tag, fromOffset, limit, newerThan, olderThan,
			private, size, onlyMeta, lowMeta, paths,
		); err != nil {
			return err
		} else if fromOffset+limit > fromOffset {
			fromOffset += limit
		} else {
			break
		}
	}
	return nil
}

// Fetch bookmark pages sequentially in reverse order and schedule the works for download.
// Stop when the bookmarks newer than the specified time were found.
func (d *Downloader) processBookmarksSeqReverse(
	userId uint64, kind queue.ItemKind, tag *string, fromOffset uint64, toOffset uint64, newerThan *time.Time,
	olderThan *time.Time, private bool, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) error {
	hasNewerTime, err := false, error(nil)
	for toOffset > fromOffset && !hasNewerTime {
		limit := min(BOOKMARK_PAGE_SIZE, toOffset-fromOffset)
		if _, _, hasNewerTime, err = d.fetchBookmarksAndScheduleWorks(
			userId, kind, tag, toOffset-limit, limit, newerThan, olderThan,
			private, size, onlyMeta, lowMeta, paths,
		); err != nil {
			return err
		} else if toOffset-limit < toOffset {
			toOffset -= limit
		} else {
			break
		}
	}
	return nil
}

// Fetch bookmarks and then schedule the works for download. Ignores works outside of the specified time range.
// Returns total count of bookmarks and whether the bookmarks outside of the specified time range were found.
func (d *Downloader) fetchBookmarksAndScheduleWorks(
	userId uint64, kind queue.ItemKind, tag *string, offset uint64, limit uint64, newerThan *time.Time,
	olderThan *time.Time, private bool, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) (total uint64, hasOlderTime bool, hasNewerTime bool, err error) {
	var successPrefix = fmt.Sprintf("fetched %v bookmarks page", kind)
	var errorPrefix = fmt.Sprintf("failed to fetch %v bookmarks page", kind)
	var noResultsPrefix = fmt.Sprintf("no %v bookmarks found", kind)

	if kind != queue.ItemKindArtwork && kind != queue.ItemKindNovel {
		err := fmt.Errorf("invalid work type: %v", uint8(kind))
		logext.Error("%v: %v", bookmarksLogMessage(errorPrefix, userId, tag, nil), err)
		return 0, false, false, err
	}

	var results []client.BookmarkResult
	if kind == queue.ItemKindArtwork {
		results, total, err = d.client.ArtworkBookmarksAuthorized(userId, tag, offset, limit, private)
	} else {
		results, total, err = d.client.NovelBookmarksAuthorized(userId, tag, offset, limit, private)
	}
	logext.MaybeSuccess(err, bookmarksLogMessage(successPrefix, userId, tag, &offset))
	logext.MaybeError(err, bookmarksLogMessage(errorPrefix, userId, tag, &offset))
	if err != nil {
		return 0, false, false, err
	}
	if len(results) == 0 {
		logext.Warning(bookmarksLogMessage(noResultsPrefix, userId, tag, &offset))
	}

	for _, result := range results {
		if result.BookmarkedTime == nil && newerThan != nil && olderThan != nil {
			err := fmt.Errorf("bookmarked time is missing in %v", result.Work)
			logext.Error("%v %v: %v", bookmarksLogMessage("failed to schedule", userId, tag, &offset), kind, err)
		} else if result.Work.Id == nil {
			err := fmt.Errorf("work id is missing in %v", result.Work)
			logext.Error("%v %v: %v", bookmarksLogMessage("failed to schedule", userId, tag, &offset), kind, err)
		} else if newerThan != nil && result.BookmarkedTime.Before(*newerThan) {
			hasOlderTime = true
		} else if olderThan != nil && result.BookmarkedTime.After(*olderThan) {
			hasNewerTime = true
		} else {
			d.ScheduleWithKnown(
				[]uint64{*result.Work.Id}, kind, size, onlyMeta, paths,
				result.Work, result.ImageUrl, lowMeta,
			)
		}
	}

	if hasOlderTime || hasNewerTime {
		logext.Info(
			bookmarksLogMessage("found bookmarks outside of the specified time range", userId, tag, &offset),
		)
	}
	return total, hasOlderTime, hasNewerTime, nil
}

func bookmarksLogMessage(prefix string, userId uint64, tag *string, offset *uint64) string {
	builder := strings.Builder{}
	builder.WriteString(prefix)
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
