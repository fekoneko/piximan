package downloader

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
)

// Schedule download. Run() to start downloading.
func (d *Downloader) Schedule(
	ids []uint64, kind queue.ItemKind, size image.Size, onlyMeta bool, paths []string,
) {
	d.downloadQueueMutex.Lock()
	defer d.downloadQueueMutex.Unlock()

	for _, id := range ids {
		d.downloadQueue.Push(queue.Item{
			Id:       id,
			Kind:     kind,
			Size:     size,
			OnlyMeta: onlyMeta,
			Paths:    paths,
		})
	}
	d.logger.ExpectWorks(len(ids))
}

// Schedule download with additional work metadata if available. Run() to start downloading.
func (d *Downloader) ScheduleWithKnown(
	ids []uint64, kind queue.ItemKind, size image.Size, onlyMeta bool, paths []string,
	work *work.Work, imageUrl *string, lowMeta bool,
) {
	d.downloadQueueMutex.Lock()
	defer d.downloadQueueMutex.Unlock()

	for _, id := range ids {
		d.downloadQueue.Push(queue.Item{
			Id:       id,
			Kind:     kind,
			Size:     size,
			OnlyMeta: onlyMeta,
			Paths:    paths,
			Work:     work,
			ImageUrl: imageUrl,
			LowMeta:  lowMeta,
		})
	}
	d.logger.ExpectWorks(len(ids))
}

// Append specified queue to the downloader queue. Run() to start downloading.
func (d *Downloader) ScheduleQueue(q *queue.Queue) {
	d.downloadQueueMutex.Lock()
	defer d.downloadQueueMutex.Unlock()

	d.downloadQueue.Push(*q...)
	d.logger.ExpectWorks(len(*q))
}

// Run the downloader.
// The operation must be waited for with WaitNext() or WaitDone() after this method.
func (d *Downloader) Run() {
	d.downloadingMutex.Lock()
	downloading := d.downloading
	d.downloadingMutex.Unlock()

	if !downloading {
		go d.superviseDownload()
		go d.superviseCrawl()
	}
}

// Block until next work is downloaded. Returns nil if there are no more works to download.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (d *Downloader) WaitNext() *work.Work {
	// TODO: return errors as well
	return <-d.channel
}

// Block until all works are downloaded.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (d *Downloader) WaitDone() {
	// TODO: return errors as well
	for d.WaitNext() != nil {
	}
}

// TODO: make supervisers prettier

// Meant to be run in a separate goroutine. Spawns download goroutines from downloadQueue
// until it is empty and no crawling is happening. Sets d.downloading to false when done.
func (d *Downloader) superviseDownload() {
	d.downloadingMutex.Lock()
	d.downloading = true
	d.downloadingMutex.Unlock()

	for {
		d.numDownloadingCond.L.Lock()
		for d.numDownloading >= downloadPendingLimit {
			d.numDownloadingCond.Wait()
		}

		d.downloadQueueMutex.Lock()
		item := d.downloadQueue.Pop()
		d.downloadQueueMutex.Unlock()

		if item == nil {
			d.numDownloadingCond.L.Unlock()

			if d.waitNextCrawled() {
				continue
			} else {
				break
			}
		}

		d.numDownloading++
		d.numDownloadingCond.L.Unlock()

		go func() {
			d.downloadItem(item)

			d.numDownloadingCond.L.Lock()
			d.numDownloading--
			d.numDownloadingCond.Broadcast()
			d.numDownloadingCond.L.Unlock()
		}()
	}

	d.numDownloadingCond.L.Lock()
	for d.numDownloading > 0 {
		d.numDownloadingCond.Wait()
	}
	d.numDownloadingCond.L.Unlock()

	d.downloadingMutex.Lock()
	d.downloading = false
	d.downloadingMutex.Unlock()

	d.channel <- nil
}

// Meant to be run in a separate goroutine.
// Spawns crawl goroutines from crawlQueue until it is empty
func (d *Downloader) superviseCrawl() {
	d.crawlingCond.L.Lock()
	defer d.crawlingCond.L.Unlock()

	for {
		for d.crawling {
			d.crawlingCond.Wait()
		}

		d.crawlQueueMutex.Lock()
		if len(d.crawlQueue) == 0 {
			d.crawlQueueMutex.Unlock()
			break
		}
		crawl := d.crawlQueue[0]
		d.crawlQueue = d.crawlQueue[1:]
		d.crawlQueueMutex.Unlock()
		d.crawling = true

		go func() {
			if err := crawl(); err == nil {
				d.logger.AddSuccessfulCrawl()
			} else if err == ErrSkipped {
				d.logger.AddSkippedCrawl()
			} else {
				d.logger.AddFailedCrawl()
			}

			d.crawlingCond.L.Lock()
			d.crawling = false
			d.crawlingCond.Broadcast()
			d.crawlingCond.L.Unlock()
		}()
	}
}

// Returns false if already crawled, otherwise waits for next crawl task to finish and returns true.
func (d *Downloader) waitNextCrawled() bool {
	d.crawlingCond.L.Lock()
	defer d.crawlingCond.L.Unlock()

	d.crawlQueueMutex.Lock()
	numCrawlTasks := len(d.crawlQueue)
	d.crawlQueueMutex.Unlock()

	if !d.crawling && numCrawlTasks == 0 {
		return false
	}
	d.crawlingCond.Wait()

	return true
}

func (d *Downloader) downloadItem(item *queue.Item) {
	var w *work.Work
	var err error

	isArtwork := item.Kind == queue.ItemKindArtwork
	isNovel := item.Kind == queue.ItemKindNovel
	withWork := item.Work != nil
	withImage := item.ImageUrl != nil
	lowMeta, onlyMeta := item.LowMeta, item.OnlyMeta

	if !isArtwork && !isNovel {
		err = fmt.Errorf("invalid work type: %v", uint8(item.Kind))
		d.logger.Error("failed to pick work %v for download: %v", item.Id, err)
		return
	}

	switch {
	case isNovel && !onlyMeta && !withImage:
		w, err = d.Novel(item.Id, item.Paths)
	case isNovel && !onlyMeta && withImage:
		w, err = d.NovelWithKnown(item.Id, *item.ImageUrl, item.Paths)
	case isNovel && onlyMeta && !(withWork && lowMeta):
		w, err = d.NovelMeta(item.Id, item.Paths)
	case isNovel && onlyMeta && withWork && lowMeta:
		w, err = d.LowNovelMetaWithKnown(item.Id, item.Work, item.Paths)
	case isArtwork && !onlyMeta && !(withWork && withImage):
		w, err = d.Artwork(item.Id, item.Size, item.Paths)
	case isArtwork && !onlyMeta && withWork && withImage && !lowMeta:
		w, err = d.ArtworkWithKnown(item.Id, item.Size, item.Work, *item.ImageUrl, item.Paths)
	case isArtwork && !onlyMeta && withWork && withImage && lowMeta:
		w, err = d.LowArtworkWithKnown(item.Id, item.Size, item.Work, *item.ImageUrl, item.Paths)
	case isArtwork && onlyMeta && (!withWork || !lowMeta):
		w, err = d.ArtworkMeta(item.Id, item.Paths)
	case isArtwork && onlyMeta && withWork && lowMeta:
		w, err = d.LowArtworkMetaWithKnown(item.Id, item.Work, item.Paths)
	default:
		err = fmt.Errorf("impossible combination of work type, known metadata, low-meta and only-meta")
		d.logger.Error("failed to pick work %v for download: %v", item.Id, err)
	}

	if err == nil && w != nil {
		d.channel <- w
		d.logger.AddSuccessfulWork()
	} else if err == ErrSkipped {
		d.logger.AddSkippedWork()
	} else if w == nil {
		d.logger.Error("received empty work %v from the task", item.Id)
		d.logger.AddFailedWork(item.Id)
	} else {
		d.logger.AddFailedWork(item.Id)
	}
}
