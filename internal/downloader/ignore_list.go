package downloader

import "github.com/fekoneko/piximan/internal/downloader/queue"

func (d *Downloader) SetSkipList(list *queue.SkipList) {
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
