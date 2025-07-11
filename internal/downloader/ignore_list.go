package downloader

import "github.com/fekoneko/piximan/internal/downloader/queue"

func (d *Downloader) SetIgnoreList(list *queue.IgnoreList) {
	d.ignoreListMutex.Lock()
	defer d.ignoreListMutex.Unlock()

	d.ignoreList = list
}

func (d *Downloader) ignored(id uint64, kind queue.ItemKind, silent bool) bool {
	d.ignoreListMutex.Lock()
	defer d.ignoreListMutex.Unlock()

	if d.ignoreList == nil {
		return false
	}
	if d.ignoreList.Contains(id, kind) {
		if !silent {
			d.logger.Info("skipping %v %v as it was already downloaded", kind.String(), id)
		}
		return true
	}
	return false
}
