package downloader

import (
	"fmt"
	"math"
	"strings"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/utils"
)

// Schedule bookmarks of authorized user for download. Run() to start downloading.
func (d *Downloader) ScheduleMyBookmarks(
	kind queue.ItemKind, tag *string, from *uint64, to *uint64, private bool,
	size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) {
	d.crawlQueueMutex.Lock()
	defer d.crawlQueueMutex.Unlock()

	d.crawlQueue = append(d.crawlQueue, func() error {
		userId, err := d.client.MyIdAutorized()
		d.logger.MaybeSuccess(err, "fetched authorizeed user id")
		d.logger.MaybeError(err, "failed to fetch authorizeed user id")
		if err != nil {
			return err
		}

		d.ScheduleBookmarks(userId, kind, tag, from, to, private, size, onlyMeta, lowMeta, paths)
		return nil
	})
	d.logger.Info("created crawl task to fetch authorizeed user id")
}

// Schedule bookmarks for download. Run() to start downloading.
func (d *Downloader) ScheduleBookmarks(
	userId uint64, kind queue.ItemKind, tag *string, from *uint64, to *uint64, private bool,
	size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) {
	d.crawlQueueMutex.Lock()
	defer d.crawlQueueMutex.Unlock()

	fromOffset := utils.FromPtr(from, 0)

	d.crawlQueue = append(d.crawlQueue, func() error {
		limit := min(100, utils.FromPtr(to, math.MaxUint64)-fromOffset)
		total, err := d.scheduleBookmarksPage(
			userId, kind, tag, fromOffset, limit, private, size, onlyMeta, lowMeta, paths,
		)
		if err != nil {
			return err
		}

		d.crawlQueueMutex.Lock()
		defer d.crawlQueueMutex.Unlock()

		toOffset := min(utils.FromPtr(to, total), total)
		offset := fromOffset + 100
		numTasks := 0

		for offset < toOffset {
			limit := min(100, toOffset-offset)
			d.crawlQueue = append(d.crawlQueue, func() error {
				_, err := d.scheduleBookmarksPage(
					userId, kind, tag, offset, limit, private, size, onlyMeta, lowMeta, paths,
				)
				return err
			})
			offset += limit
			numTasks++
		}

		if numTasks > 0 {
			d.logger.Info(
				bookmarksLogMessage("created %v bookmarks crawl %v", userId, tag, nil),
				numTasks, utils.If(numTasks == 1, "task", "tasks"),
			)
		}

		return nil
	})
	d.logger.Info(bookmarksLogMessage("created bookmarks crawl task", userId, tag, &fromOffset))
}

// Fetch bookmarks and then schedule the works for download, returns total count of bookmarks
func (d *Downloader) scheduleBookmarksPage(
	userId uint64, kind queue.ItemKind, tag *string, offset uint64, limit uint64, private bool,
	size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) (total uint64, err error) {
	var successPrefix = fmt.Sprintf("fetched %v bookmarks page", kind)
	var errorPrefix = fmt.Sprintf("failed to fetch %v bookmarks page", kind)
	var noResultsPrefix = fmt.Sprintf("no %v bookmarks found", kind)

	if kind != queue.ItemKindArtwork && kind != queue.ItemKindNovel {
		err := fmt.Errorf("invalid work type: %v", uint8(kind))
		d.logger.Error("%v: %v", bookmarksLogMessage(errorPrefix, userId, tag, nil), err)
		return 0, err
	}

	var results []client.BookmarkResult
	if kind == queue.ItemKindArtwork {
		results, total, err = d.client.ArtworkBookmarksAuthorized(
			userId, tag, offset, limit, private,
		)
	} else {
		results, total, err = d.client.NovelBookmarksAuthorized(
			userId, tag, offset, limit, private,
		)
	}
	d.logger.MaybeSuccess(err, bookmarksLogMessage(successPrefix, userId, tag, &offset))
	d.logger.MaybeError(err, bookmarksLogMessage(errorPrefix, userId, tag, &offset))
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		d.logger.Warning(bookmarksLogMessage(noResultsPrefix, userId, tag, &offset))
	}

	for _, result := range results {
		if result.Work.Id == nil {
			err := fmt.Errorf("work id is missing in %v", result.Work)
			d.logger.Error("%v %v: %v", bookmarksLogMessage("failed to schedule", userId, tag, &offset), kind, err)
			continue
		}
		d.ScheduleWithKnown(
			[]uint64{*result.Work.Id}, kind, size, onlyMeta, paths,
			result.Work, result.ImageUrl, lowMeta,
		)
	}

	return total, nil
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
