package downloader

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/logext"
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
}

// Schedule download with additional work metadata if available. Run() to start downloading.
func (d *Downloader) ScheduleWithWork(
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
}

// Merge queue to the downloader queue. Run() to start downloading.
func (d *Downloader) ScheduleQueue(q *queue.Queue) {
	d.downloadQueueMutex.Lock()
	defer d.downloadQueueMutex.Unlock()

	d.downloadQueue.Push(*q...)
}

// Run the downloader. Need to WaitNext() or WaitDone() to get the results.
func (d *Downloader) Run() {
	d.downloadingMutex.Lock()
	defer d.downloadingMutex.Unlock()

	if !d.downloading {
		go d.superviseDownload()
		go d.superviseCrawl()
	}
}

// Block until next work is downloaded. Returns nil if there are no more works to download.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (d *Downloader) WaitNext() *work.Work {
	return <-d.channel
}

// Block until all works are downloaded.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (d *Downloader) WaitDone() {
	for d.WaitNext() != nil {
	}
}

// Meant to be run in a separate goroutine. Spawns download goroutines from downloadQueu
// until it is empty and no crawling is happening. Sets d.downloading to false when done.
func (d *Downloader) superviseDownload() {
	d.downloadingMutex.Lock()
	d.downloading = true
	d.downloadingMutex.Unlock()

	d.numDownloadingCond.L.Lock()
	defer d.numDownloadingCond.L.Unlock()

	for {
		for d.numDownloading >= DOWNLOAD_PENDING_LIMIT {
			d.numDownloadingCond.Wait()
		}

		d.downloadQueueMutex.Lock()
		item := d.downloadQueue.Pop()
		d.downloadQueueMutex.Unlock()
		if item == nil && !d.waitCrawled() {
			break
		}
		d.numDownloading++

		go func() {
			d.downloadItem(item)

			d.numDownloadingCond.L.Lock()
			d.numDownloading--
			d.numDownloadingCond.Broadcast()
			d.numDownloadingCond.L.Unlock()
		}()
	}

	d.downloadingMutex.Lock()
	d.downloading = false
	d.downloadingMutex.Unlock()

	d.channel <- nil
}

// Meant to be run in a separate goroutine. Spawns crawl goroutines from crawlQueue
// until it is empty. Sets d.crawling to false when done.
func (d *Downloader) superviseCrawl() {
	d.crawlingCond.L.Lock()
	d.crawling = true
	d.crawlingCond.L.Unlock()

	d.numCrawlingCond.L.Lock()
	defer d.numCrawlingCond.L.Unlock()

	for {
		for d.numCrawling >= CRAWL_PENDING_LIMIT {
			d.numCrawlingCond.Wait()
		}

		d.crawlQueueMutex.Lock()
		if len(d.crawlQueue) == 0 {
			break
		}
		crawl := d.crawlQueue[0]
		d.crawlQueue = d.crawlQueue[1:]
		defer d.crawlQueueMutex.Unlock()

		go func() {
			crawl() // TODO: count these errors as well as download errors

			d.numCrawlingCond.L.Lock()
			d.numCrawling--
			d.numCrawlingCond.Broadcast()
			d.numCrawlingCond.L.Unlock()
		}()
	}

	d.crawlingCond.L.Lock()
	d.crawling = false
	d.crawlingCond.Broadcast()
	d.crawlingCond.L.Unlock()
}

// Returns false if already crawled, otherwise waits for crawling to be finished and returns true.
func (d *Downloader) waitCrawled() bool {
	d.crawlingCond.L.Lock()
	defer d.crawlingCond.L.Unlock()

	if !d.crawling {
		return false
	}

	for d.crawling {
		d.crawlingCond.Wait()
	}

	return true
}

func (d *Downloader) downloadItem(item *queue.Item) {
	var w *work.Work
	var err error

	// TODO: take in concideration queue.Item.Work and queue.Item.ImageUrl
	if item.Kind == queue.ItemKindNovel && item.OnlyMeta {
		w, err = d.DownloadNovelMeta(item.Id, item.Paths)
	} else if item.Kind == queue.ItemKindNovel {
		w, err = d.DownloadNovel(item.Id, item.Paths)
	} else if item.Kind == queue.ItemKindArtwork && item.OnlyMeta {
		w, err = d.DownloadArtworkMeta(item.Id, item.Paths)
	} else if item.Kind == queue.ItemKindArtwork {
		w, err = d.DownloadArtwork(item.Id, item.Size, item.Paths)
	} else {
		err = fmt.Errorf("invalid work type: %v", uint8(item.Kind))
		logext.Error("failed to pick work %v for download: %v", item.Id, err)
	}

	if err == nil {
		d.channel <- w
	}
}
