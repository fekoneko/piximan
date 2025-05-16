package downloader

import (
	"fmt"
	"strings"

	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
)

// TODO: don't block and schedule fetching bookmarks after Run() as well
//       separate queues for work downloads and bookmark fetches so that we can keep goroutines
//       for each queue running at the same time to ensure that bookmark fetches don't block
//       downloading of works (work downloading still may use free i.pximg.net slots)
//       one concurrently running goroutine for bookmarks and 5 for works

// Fetch artwork bookmarks and then schedule them for download, blocks until done.
// Use Run() to start downloading after this function is finished.
func (d *Downloader) ScheduleArtworkBookmarks(
	userId uint64, tag *string, size image.Size, onlyMeta bool, lowMeta bool, paths []string,
) error {
	if d.sessionId == nil {
		err := fmt.Errorf("authorization is required")
		logext.Error("%v: %v", bookmarksLogPrefix("failed to fetch artwork bookmarks", userId, tag, nil), err)
		return err
	}

	// TODO: fetch with other offsets as well
	offset := uint(0)
	results, err := fetch.ArtworkBookmarksAuthorized(d.client, userId, tag, offset, 100, *d.sessionId)
	logext.MaybeSuccess(err, bookmarksLogPrefix("fetched artwork bookmarks", userId, tag, &offset))
	logext.MaybeError(err, bookmarksLogPrefix("failed to fetch artwork bookmarks", userId, tag, &offset))
	if err != nil {
		return err
	}

	for _, result := range results {
		d.ScheduleWithWork(
			[]uint64{result.Work.Id}, queue.ItemKindArtwork, size, onlyMeta, paths,
			result.Work, &result.ThumbnailUrl, lowMeta,
		)
	}

	return nil
}

// Fetch novel bookmarks and then schedule them for download, blocks until done.
// Use Run() to start downloading after this function is finished.
func (d *Downloader) ScheduleNovelBookmarks(
	userId uint64, tag *string, onlyMeta bool, lowMeta bool, paths []string,
) error {
	if d.sessionId == nil {
		err := fmt.Errorf("authorization is required")
		logext.Error("%v: %v", bookmarksLogPrefix("failed to fetch novel bookmarks", userId, tag, nil), err)
		return err
	}

	// TODO: fetch with other offsets as well
	offset := uint(0)
	results, err := fetch.NovelBookmarksAuthorized(d.client, userId, tag, offset, 100, *d.sessionId)
	logext.MaybeSuccess(err, bookmarksLogPrefix("fetched novel bookmarks", userId, tag, &offset))
	logext.MaybeError(err, bookmarksLogPrefix("failed to fetch novel bookmarks", userId, tag, &offset))
	if err != nil {
		return err
	}

	for _, result := range results {
		d.ScheduleWithWork(
			[]uint64{result.Work.Id}, queue.ItemKindNovel, image.SizeDefault, onlyMeta, paths,
			result.Work, &result.CoverUrl, lowMeta,
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
