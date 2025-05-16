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
	d.downloadQueue.Push(*q...)
}

// TODO: get one task from crawlQueue and run it, don't end running until both queues are empty
//       d.crawlWaitGroup.Wait()
//       d.crawlWaitGroup.Add(1)
//       defer d.crawlWaitGroup.Done()

// Run the downloader. Need to WaitNext() or WaitDone() to get the results.
func (d *Downloader) Run() {
	d.numPendingMutex.Lock()
	defer d.numPendingMutex.Unlock()

	for d.numPending < PENDING_LIMIT {
		item := d.downloadQueue.Pop()
		if item == nil {
			break
		}
		go d.downloadItem(item)
		d.numPending++
	}
}

// Block until next work is downloaded. Returns nil if there are no more works to download.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (d *Downloader) WaitNext() *work.Work {
	var w *work.Work
	for w == nil && d.NumRemaining() > 0 {
		w = <-d.channel
	}
	return w
}

// Block until all works are downloaded.
// Use WaitNext() or WaitDone() only in one place at a time to receive all the results.
func (d *Downloader) WaitDone() {
	for d.WaitNext() != nil {
	}
}

// Number of pending and queued works.
func (d *Downloader) NumRemaining() int {
	d.numPendingMutex.Lock()
	defer d.numPendingMutex.Unlock()

	return len(d.downloadQueue) + d.numPending
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
	} else {
		d.channel <- nil
	}

	d.numPendingMutex.Lock()
	d.numPending--
	d.numPendingMutex.Unlock()

	d.Run()
}
