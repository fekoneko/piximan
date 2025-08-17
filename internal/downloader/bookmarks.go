package downloader

import (
	"fmt"
	"math"
	"strings"

	"github.com/fekoneko/piximan/internal/client"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/syncext"
	"github.com/fekoneko/piximan/internal/utils"
)

// Schedule bookmarks of authorized user for download. Run() to start downloading.
// If untilSkipped is true, the crawler will stop fetching new bookmark pages once it encounters
// a fully skipped one. Use this to conserve requests when synching only new bookmarks.
func (d *Downloader) ScheduleMyBookmarks(
	kind queue.ItemKind, tags *[]string, from *uint64, to *uint64, private bool,
	size imageext.Size, onlyMeta bool, lowMeta bool, untilSkipped bool, paths []string,
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

		d.ScheduleBookmarks(userId, kind, tags, from, to, private, size, onlyMeta, lowMeta, untilSkipped, paths)
		return nil
	})
	d.logger.Info("created crawl task to fetch authorizeed user id")
	d.logger.ExpectCrawls(1)
}

// Schedule bookmarks for download. Run() to start downloading.
// If untilSkipped is true, the crawler will stop fetching new bookmark pages once it encounters
// a fully skipped one. Use this to conserve requests when synching only new bookmarks.
func (d *Downloader) ScheduleBookmarks(
	userId uint64, kind queue.ItemKind, tags *[]string, from *uint64, to *uint64, private bool,
	size imageext.Size, onlyMeta bool, lowMeta bool, untilSkipped bool, paths []string,
) {
	if tags == nil {
		d.scheduleBookmarks(
			userId, kind, nil, from, to, private,
			size, onlyMeta, lowMeta, untilSkipped, paths,
		)
	} else {
		seen := make(map[string]bool)
		for _, tag := range *tags {
			if seen[tag] {
				continue
			}
			d.scheduleBookmarks(
				userId, kind, &tag, from, to, private,
				size, onlyMeta, lowMeta, untilSkipped, paths,
			)
			seen[tag] = true
		}
	}
}

// Schedule the task that will fetch total bookmarks count from the first page and then create tasks
// for the rest of the pages.
func (d *Downloader) scheduleBookmarks(
	userId uint64, kind queue.ItemKind, tag *string, from *uint64, to *uint64, private bool,
	size imageext.Size, onlyMeta bool, lowMeta bool, untilSkipped bool, paths []string,
) {
	d.crawlQueueMutex.Lock()
	defer d.crawlQueueMutex.Unlock()

	fromOffset := utils.FromPtr(from, 0)

	d.crawlQueue = append(d.crawlQueue, func() error {
		limit := min(100, utils.FromPtr(to, math.MaxUint64)-fromOffset)
		total, allIngored, err := d.fetchBookmarksPageAndSchedule(
			userId, kind, tag, fromOffset, limit, private, size, onlyMeta, lowMeta, paths, nil,
		)
		if err != nil {
			return err
		} else if untilSkipped && allIngored {
			d.logger.Info(fullySkippedMessage)
			return nil
		}

		d.crawlQueueMutex.Lock()
		defer d.crawlQueueMutex.Unlock()

		toOffset := min(utils.FromPtr(to, total), total)
		offset := fromOffset + 100
		numTasks := 0
		signal := syncext.NewSignal()

		for offset < toOffset {
			limit := min(100, toOffset-offset)
			currentOffset := offset
			d.crawlQueue = append(d.crawlQueue, func() error {
				_, allIngored, err := d.fetchBookmarksPageAndSchedule(
					userId, kind, tag, currentOffset, limit, private, size, onlyMeta, lowMeta, paths, &signal,
				)
				if untilSkipped && allIngored {
					d.logger.Info(fullySkippedMessage)
					signal.Cancel()
				}
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
			d.logger.ExpectCrawls(numTasks)
		}

		return nil
	})
	d.logger.Info(bookmarksLogMessage("created bookmarks crawl task", userId, tag, &fromOffset))
	d.logger.ExpectCrawls(1)
}

// Fetch bookmarks and then schedule the works for download, returns total count of bookmarks
// Can be cancelled with provided signal until work download tasks were scheduled.
func (d *Downloader) fetchBookmarksPageAndSchedule(
	userId uint64, kind queue.ItemKind, tag *string, offset uint64, limit uint64, private bool,
	size imageext.Size, onlyMeta bool, lowMeta bool, paths []string, signal *syncext.Signal,
) (total uint64, allIngored bool, err error) {
	if signal != nil && signal.Cancelled() {
		return 0, false, ErrSkipped
	}

	var successPrefix = fmt.Sprintf("fetched %v bookmarks page", kind)
	var errorPrefix = fmt.Sprintf("failed to fetch %v bookmarks page", kind)
	var noResultsPrefix = fmt.Sprintf("no %v bookmarks found", kind)

	if kind != queue.ItemKindArtwork && kind != queue.ItemKindNovel {
		err := fmt.Errorf("invalid work type: %v", uint8(kind))
		d.logger.Error("%v: %v", bookmarksLogMessage(errorPrefix, userId, tag, nil), err)
		return 0, false, err
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
	if signal != nil && signal.Cancelled() {
		d.logger.AddSkippedCrawl()
		return 0, false, nil
	}
	d.logger.MaybeSuccess(err, bookmarksLogMessage(successPrefix, userId, tag, &offset))
	d.logger.MaybeError(err, bookmarksLogMessage(errorPrefix, userId, tag, &offset))
	if err != nil {
		return 0, false, err
	}
	if len(results) == 0 {
		d.logger.Warning(bookmarksLogMessage(noResultsPrefix, userId, tag, &offset))
		return 0, false, nil
	}

	numSkipped := 0
	for _, result := range results {
		if result.Work.Id == nil {
			err := fmt.Errorf("work id is missing in %v", result.Work)
			d.logger.Error(
				"%v %v: %v", bookmarksLogMessage("failed to schedule", userId, tag, &offset), kind, err,
			)
		} else if !d.skipped(*result.Work.Id, kind, true) {
			d.ScheduleWithKnown(
				[]uint64{*result.Work.Id}, kind, size, onlyMeta, paths,
				result.Work, result.ImageUrl, lowMeta,
			)
		} else {
			numSkipped++
		}
	}
	if numSkipped > 0 {
		d.logger.Info(
			"skipping %v %v%v as %v already downloaded", numSkipped, kind.String(),
			utils.If(numSkipped == 1, "", "s"), utils.If(numSkipped == 1, "it was", "they were"),
		)
	}
	return total, numSkipped >= len(results), nil
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

var fullySkippedMessage = "found fully skipped bookmarks page, stopped crawling new pages to conserve requests"
