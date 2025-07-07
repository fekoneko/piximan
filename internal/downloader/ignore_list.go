package downloader

import "github.com/fekoneko/piximan/internal/downloader/queue"

func (d *Downloader) SetIgnoreList(list *queue.IgnoreList) {
	d.ignoreListMutex.Lock()
	defer d.ignoreListMutex.Unlock()

	d.ignoreList = list
}

func (d *Downloader) artworkIgnored(id uint64) bool {
	d.ignoreListMutex.Lock()
	defer d.ignoreListMutex.Unlock()

	if d.ignoreList == nil {
		return false
	}
	return d.ignoreList.Contains(id, queue.ItemKindArtwork)
}

func (d *Downloader) novelIgnored(id uint64) bool {
	d.ignoreListMutex.Lock()
	defer d.ignoreListMutex.Unlock()

	if d.ignoreList == nil {
		return false
	}
	return d.ignoreList.Contains(id, queue.ItemKindNovel)
}
