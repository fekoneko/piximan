package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
)

// Schedule download. Run() to start downloading.
func (d *Downloader) Schedule(
	ids []uint64, kind queue.ItemKind, size image.Size, onlyMeta bool, paths []string,
) {
	for _, id := range ids {
		d.queue.Push(queue.Item{
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
	work *work.Work, imageUrl *string,
) {
	for _, id := range ids {
		d.queue.Push(queue.Item{
			Id:       id,
			Kind:     kind,
			Size:     size,
			OnlyMeta: onlyMeta,
			Paths:    paths,
			Work:     work,
			ImageUrl: imageUrl,
		})
	}
}

// Merge queue to the downloader queue. Run() to start downloading.
func (d *Downloader) ScheduleQueue(q *queue.Queue) {
	d.queue.Push(*q...)
}

// Run the downloader. Need to WaitNext() or WaitDone() to get the results.
func (d *Downloader) Run() {
	d.numPendingMutex.Lock()
	defer d.numPendingMutex.Unlock()

	for d.numPending < PENDING_LIMIT {
		item := d.queue.Pop()
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

	return len(d.queue) + d.numPending
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
