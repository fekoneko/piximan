package downloader

import (
	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
)

func (d *Downloader) Schedule(
	ids []uint64, kind queue.ItemKind, size image.Size, onlyMeta bool, paths []string,
) {
	for _, id := range ids {
		d.queue.Push(queue.Item{Id: id, Kind: kind, Size: size, OnlyMeta: onlyMeta, Paths: paths})
	}
}

func (d *Downloader) ScheduleQueue(q *queue.Queue) {
	d.queue.Push(*q...)
}

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

func (d *Downloader) WaitNext() *work.Work {
	var w *work.Work
	for w == nil && d.NumRemaining() > 0 {
		w = <-d.channel
	}
	return w
}

func (d *Downloader) WaitDone() {
	for d.WaitNext() != nil {
	}
}

func (d *Downloader) NumRemaining() int {
	d.numPendingMutex.Lock()
	defer d.numPendingMutex.Unlock()

	return len(d.queue) + d.numPending
}

func (d *Downloader) downloadItem(item *queue.Item) {
	var w *work.Work
	var err error

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
