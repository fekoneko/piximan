package downloader

import (
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/downloader/skiplist"
)

// Get skip list that is used to skip previously downloaded works. Thread-safe, may be nil.
func (d *Downloader) SkipList() *skiplist.SkipList {
	d.skipListMutex.Lock()
	defer d.skipListMutex.Unlock()
	return d.skipList
}

// Set skip list that is used to skip previously downloaded works. Thread-safe.
func (d *Downloader) SetSkipList(list *skiplist.SkipList) {
	d.skipListMutex.Lock()
	defer d.skipListMutex.Unlock()
	d.skipList = list
}

func (d *Downloader) skipped(id uint64, kind queue.ItemKind, silent bool) bool {
	d.skipListMutex.Lock()
	defer d.skipListMutex.Unlock()

	if d.skipList == nil {
		return false
	}
	if d.skipList.Contains(id, kind) {
		if !silent {
			d.logger.Info("skipping %v %v as it was already downloaded", kind.String(), id)
		}
		return true
	}
	return false
}
