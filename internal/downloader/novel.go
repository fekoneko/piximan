package downloader

import (
	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
)

// Download only novel metadata and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelMeta(id uint64, paths []string) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindNovel, false) || !d.matchNovelId(id) {
		return nil, ErrSkipped
	}
	d.logger.Info("downloading metadata for novel %v", id)

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
	if d.skipped(id, queue.ItemKindNovel, false) || !d.matchNovel(id, w, true) {
		return nil, ErrSkipped
	}
	assets := []fsext.Asset{}
	return w, d.writeWork(id, queue.ItemKindNovel, w, assets, true, paths)
}

// Download novel with all assets and metadata and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using Schedule().
func (d *Downloader) Novel(id uint64, size imageext.Size, paths []string) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindNovel, false) || !d.matchNovelId(id) {
		return nil, ErrSkipped
	}
	d.logger.Info("downloading novel %v", id)

	w, coverUrl, uploadedImages, pixivImages, pages, err := d.novelMeta(id, &size)
	if err != nil {
		return nil, err
	} else if !d.matchNovel(id, w, false) {
		return nil, ErrSkipped
	}

	coverChannel := make(chan *fsext.Asset, 1)
	imagesChannel := make(chan map[int]fsext.Asset, 1)
	errorChannel := make(chan error, 1)

	go d.novelCoverAssetChannel(id, *coverUrl, coverChannel, errorChannel)
	go d.novelImageAssetsChannel(id, size, uploadedImages, pixivImages, imagesChannel, errorChannel)

	var coverAsset *fsext.Asset
	var imageAssets map[int]fsext.Asset

	for range 2 {
		select {
		case coverAsset = <-coverChannel:
		case imageAssets = <-imagesChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	assets := combineAssets(coverAsset, imageAssets, pages)
	return w, d.writeWork(id, queue.ItemKindNovel, w, assets, false, paths)
}

// Download novel with cover url known in advance and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// Tries to start downloading assets as soon as possible, but if some rules dependent on full
// metadata are defined, it will wait until full metadata is received.
// For downloading multiple works consider using Schedule().
func (d *Downloader) NovelWithKnown(
	id uint64, size imageext.Size, coverUrl string, paths []string,
) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindNovel, false) {
		return nil, ErrSkipped
	} else if matches, needFull := d.matchNovelNeedFull(id, nil); !matches {
		return nil, ErrSkipped
	} else if needFull {
		return d.Novel(id, size, paths)
	}
	d.logger.Info("downloading novel %v", id)

	workChannel := make(chan *work.Work, 1)
	pagesChannel := make(chan dto.NovelPages, 1)
	imagesChannel := make(chan map[int]fsext.Asset, 1)
	coverChannel := make(chan *fsext.Asset, 1)
	errorChannel := make(chan error, 1)

	go d.novelMetaImageAssetsChannel(id, size, workChannel, pagesChannel, imagesChannel, errorChannel)
	go d.novelCoverAssetChannel(id, coverUrl, coverChannel, errorChannel)

	var w *work.Work
	var pages dto.NovelPages
	var imageAssets map[int]fsext.Asset
	var coverAsset *fsext.Asset

	for range 4 {
		select {
		case w = <-workChannel:
		case pages = <-pagesChannel:
		case imageAssets = <-imagesChannel:
		case coverAsset = <-coverChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	assets := combineAssets(coverAsset, imageAssets, pages)
	return w, d.writeWork(id, queue.ItemKindNovel, w, assets, false, paths)
}
