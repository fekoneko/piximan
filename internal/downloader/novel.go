package downloader

import (
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
)

// Download only novel metadata and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelMeta(id uint64, paths []string) (*work.Work, error) {
	if d.ignored(id, queue.ItemKindNovel, false) || !d.matchNovelId(id) {
		return nil, ErrSkipped
	}
	d.logger.Info("started downloading metadata for novel %v", id)

	w, err := d.novelOnlyMeta(id)
	if err != nil {
		return nil, err
	} else if !d.matchNovel(id, w, false) {
		return nil, ErrSkipped
	}

	assets := []fsext.Asset{}
	return w, d.writeWork(id, queue.ItemKindNovel, w, assets, true, paths)
}

// Doesn't actually make additional requests, but stores incomplete metadata, received earlier.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowNovelMetaWithKnown(id uint64, w *work.Work, paths []string) (*work.Work, error) {
	if d.ignored(id, queue.ItemKindNovel, false) || !d.matchNovel(id, w, true) {
		return nil, ErrSkipped
	}
	assets := []fsext.Asset{}
	return w, d.writeWork(id, queue.ItemKindNovel, w, assets, true, paths)
}

// Download novel with all assets and metadata and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using Schedule().
func (d *Downloader) Novel(id uint64, paths []string) (*work.Work, error) {
	if d.ignored(id, queue.ItemKindNovel, false) || !d.matchNovelId(id) {
		return nil, ErrSkipped
	}
	d.logger.Info("started downloading novel %v", id)

	w, coverUrl, pageAssets, err := d.novelMeta(id)
	if err != nil {
		return nil, err
	} else if !d.matchNovel(id, w, false) {
		return nil, ErrSkipped
	}
	coverAsset, err := d.novelCoverAsset(id, *coverUrl)
	if err != nil {
		return nil, err
	}
	assets := make([]fsext.Asset, 0, len(*pageAssets)+1)
	assets = append(assets, *pageAssets...)
	assets = append(assets, *coverAsset)
	return w, d.writeWork(id, queue.ItemKindNovel, w, assets, false, paths)
}

// Download novel with cover url known in advance and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// Tries to start downloading assets as soon as possible, but if some rules dependent on full
// metadata are defined, it will wait until full metadata is received.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelWithKnown(id uint64, coverUrl string, paths []string) (*work.Work, error) {
	if d.ignored(id, queue.ItemKindNovel, false) {
		return nil, ErrSkipped
	} else if matches, needFull := d.matchNovelNeedFull(id, nil); !matches {
		return nil, ErrSkipped
	} else if needFull {
		return d.Novel(id, paths)
	}
	d.logger.Info("started downloading novel %v", id)

	workChannel := make(chan *work.Work, 1)
	pagesChannel := make(chan *[]fsext.Asset, 1)
	coverChannel := make(chan *fsext.Asset, 1)
	errorChannel := make(chan error)

	go d.novelMetaChannel(id, workChannel, pagesChannel, errorChannel)
	go d.novelCoverAssetChannel(id, coverUrl, coverChannel, errorChannel)

	var w *work.Work
	var pageAssets *[]fsext.Asset
	var coverAsset *fsext.Asset

	for range 3 {
		select {
		case w = <-workChannel:
		case pageAssets = <-pagesChannel:
		case coverAsset = <-coverChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	assets := make([]fsext.Asset, 0, len(*pageAssets)+1)
	assets = append(assets, *pageAssets...)
	assets = append(assets, *coverAsset)
	return w, d.writeWork(id, queue.ItemKindNovel, w, assets, false, paths)
}
