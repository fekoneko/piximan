package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/storage"
)

// Download only novel metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelMeta(id uint64, paths []string) (*work.Work, error) {
	logext.Info("started downloading metadata for novel %v", id)

	w, _, _, err := fetch.NovelMeta(d.client(), id)
	logext.MaybeSuccess(err, "fetched metadata for novel %v", id)
	logext.MaybeError(err, "failed to fetch metadata for novel %v", id)
	if err != nil {
		return nil, err
	}
	if !w.Full() {
		logext.Warning("metadata for novel %v is incomplete", id)
	}

	assets := []storage.Asset{}
	return w, writeWork(id, queue.ItemKindNovel, w, assets, true, paths)
}

// Doesn't actually make additional requests, but stores incomplete metadata, received earlier.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowNovelMetaWithKnown(id uint64, w *work.Work, paths []string) (*work.Work, error) {
	assets := []storage.Asset{}
	return w, writeWork(id, queue.ItemKindNovel, w, assets, true, paths)
}

// Download novel with all assets and metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) Novel(id uint64, paths []string) (*work.Work, error) {
	logext.Info("started downloading novel %v", id)

	w, coverUrl, contentAsset, err := d.novelMeta(id)
	if err != nil {
		return nil, err
	}
	coverAsset, err := d.novelCoverAsset(id, *coverUrl)
	if err != nil {
		return nil, err
	}
	assets := []storage.Asset{*coverAsset, *contentAsset}
	return w, writeWork(id, queue.ItemKindNovel, w, assets, false, paths)
}

// Download novel with cover url known in advance and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelWithKnown(id uint64, coverUrl string, paths []string) (*work.Work, error) {
	logext.Info("started downloading novel %v", id)

	workChannel := make(chan *work.Work, 1)
	contentChannel := make(chan *storage.Asset, 1)
	coverChannel := make(chan *storage.Asset, 1)
	errorChannel := make(chan error)

	go d.novelMetaChannel(id, workChannel, contentChannel, errorChannel)
	go d.novelCoverAssetChannel(id, coverUrl, coverChannel, errorChannel)

	var w *work.Work
	var contentAsset *storage.Asset
	var coverAsset *storage.Asset

	for range 3 {
		select {
		case w = <-workChannel:
		case contentAsset = <-contentChannel:
		case coverAsset = <-coverChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	assets := []storage.Asset{*contentAsset, *coverAsset}
	return w, writeWork(id, queue.ItemKindNovel, w, assets, false, paths)
}
