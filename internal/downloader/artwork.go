package downloader

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/work"
)

// Download only artwork metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) ArtworkMeta(id uint64, paths []string) (*work.Work, error) {
	d.logger.Info("started downloading metadata for artwork %v", id)

	w, _, _, err := d.client.ArtworkMeta(id)
	d.logger.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	d.logger.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}
	if !w.Full() {
		d.logger.Warning("metadata for artwork %v is incomplete", id)
	}

	assets := []fsext.Asset{}
	return w, d.writeWork(id, queue.ItemKindArtwork, w, assets, true, paths)
}

// Doesn't actually make additional requests, but stores incomplete metadata, received earlier.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowArtworkMetaWithKnown(id uint64, w *work.Work, paths []string) (*work.Work, error) {
	assets := []fsext.Asset{}
	return w, d.writeWork(id, queue.ItemKindArtwork, w, assets, true, paths)
}

// Download artwork with all assets and metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) Artwork(id uint64, size image.Size, paths []string) (*work.Work, error) {
	d.logger.Info("started downloading artwork %v", id)

	w, firstPageUrls, thumbnailUrls, err := d.artworkMeta(id)
	if err != nil {
		return nil, err
	}

	if w.Kind == nil {
		err := fmt.Errorf("work kind is missing in %v", w)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	} else if *w.Kind == work.KindUgoira {
		assets, err := d.ugoiraAssets(id, w)
		d.writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga {
		thumbnailUrl := urlFromMap(id, thumbnailUrls)
		assets, err := d.illustMangaAssets(id, w, firstPageUrls, thumbnailUrl, size)
		d.writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else {
		err := fmt.Errorf("invalid work kind: %v", *w.Kind)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	}
}

// Download artwork with partial metadata known in advance and store it in paths. Blocks until done.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) ArtworkWithKnown(
	id uint64, size image.Size, w *work.Work, thumbnailUrl string, paths []string,
) (*work.Work, error) {
	d.logger.Info("started downloading artwork %v", id)

	workChannel := make(chan *work.Work)
	assetsChannel := make(chan []fsext.Asset)
	errorChannel := make(chan error)

	go d.artworkMetaChannel(id, workChannel, errorChannel)

	if w.Kind == nil {
		err := fmt.Errorf("work kind is missing in %v", w)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	} else if *w.Kind == work.KindUgoira {
		go d.ugoiraAssetsChannel(id, w, assetsChannel, errorChannel)
	} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga {
		go d.illustMangaAssetsChannel(id, w, nil, &thumbnailUrl, size, assetsChannel, errorChannel)
	} else {
		err := fmt.Errorf("invalid work kind: %v", *w.Kind)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	}

	var fullWork *work.Work
	var assets []fsext.Asset

	for range 2 {
		select {
		case fullWork = <-workChannel:
		case assets = <-assetsChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	return fullWork, d.writeWork(id, queue.ItemKindArtwork, fullWork, assets, false, paths)
}

// Download artwork using already available incomplete metadata and store it in paths.
// Doesn't fetch full metadata. Blocks until done.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowArtworkWithKnown(
	id uint64, size image.Size, w *work.Work, thumbnailUrl string, paths []string,
) (*work.Work, error) {
	d.logger.Info("started downloading artwork %v", id)

	if w.Kind == nil {
		err := fmt.Errorf("work kind is missing in %v", w)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	} else if *w.Kind == work.KindUgoira {
		assets, err := d.ugoiraAssets(id, w)
		d.writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga {
		assets, err := d.illustMangaAssets(id, w, nil, &thumbnailUrl, size)
		d.writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else {
		err := fmt.Errorf("invalid work kind: %v", *w.Kind)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	}
}
