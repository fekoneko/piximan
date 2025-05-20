package downloader

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/storage"
)

// TODO: when downloading bookmarks we can fetch metadata in parallel with images
//       if we even need to fetch full metadata

// Download only artwork metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) ArtworkMeta(id uint64, paths []string) (*work.Work, error) {
	logext.Info("started downloading metadata for artwork %v", id)

	w, _, _, err := fetch.ArtworkMeta(d.client(), id)
	logext.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	logext.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}
	if !w.Full() {
		logext.Warning("metadata for artwork %v is incomplete", id)
	}

	assets := []storage.Asset{}
	return w, writeWork(id, queue.ItemKindArtwork, w, assets, true, paths)
}

// Doesn't actually make additional requests, but stores incomplete metadata, received earlier.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowArtworkMetaWithKnown(
	id uint64, w *work.Work, paths []string,
) (*work.Work, error) {
	assets := []storage.Asset{}
	return w, writeWork(id, queue.ItemKindArtwork, w, assets, true, paths)
}

// Download artwork with all assets and metadata and store it in paths. Blocks until done.
// For downloading multiple works consider using Schedule().
func (d *Downloader) Artwork(id uint64, size image.Size, paths []string) (*work.Work, error) {
	logext.Info("started downloading artwork %v", id)

	w, firstPageUrls, thumbnailUrls, err := d.artworkMeta(id)
	if err != nil {
		return nil, err
	}

	if w.Kind == nil {
		err := fmt.Errorf("work kind is missing in %v", w)
		logext.Error("failed to download artwork %v: %v", id, err)
		return w, err
	} else if *w.Kind == work.KindUgoira {
		assets, err := d.ugoiraAssets(id, w)
		writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga {
		thumbnailUrl := urlFromMap(id, thumbnailUrls)
		assets, err := d.illustMangaAssets(id, w, firstPageUrls, thumbnailUrl, size)
		writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else {
		err := fmt.Errorf("invalid work kind: %v", *w.Kind)
		logext.Error("failed to download artwork %v: %v", id, err)
		return w, err
	}
}

// Download artwork with partial metadata known in advance and store it in paths. Blocks until done.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) ArtworkWithKnown(
	id uint64, size image.Size, w *work.Work, thumbnailUrl string, paths []string,
) (*work.Work, error) {
	logext.Info("started downloading artwork %v", id)

	workChannel := make(chan *work.Work)
	assetsChannel := make(chan []storage.Asset)
	errorChannel := make(chan error)

	go d.artworkMetaChannel(id, workChannel, errorChannel)

	if w.Kind == nil {
		err := fmt.Errorf("work kind is missing in %v", w)
		logext.Error("failed to download artwork %v: %v", id, err)
		return w, err
	} else if *w.Kind == work.KindUgoira {
		go d.ugoiraAssetsChannel(id, w, assetsChannel, errorChannel)
	} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga {
		go d.illustMangaAssetsChannel(id, w, nil, &thumbnailUrl, size, assetsChannel, errorChannel)
	} else {
		err := fmt.Errorf("invalid work kind: %v", *w.Kind)
		logext.Error("failed to download artwork %v: %v", id, err)
		return w, err
	}

	var fullWork *work.Work
	var assets []storage.Asset

	for range 2 {
		select {
		case fullWork = <-workChannel:
		case assets = <-assetsChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	return fullWork, writeWork(id, queue.ItemKindArtwork, fullWork, assets, false, paths)
}

// Download artwork using already available incomplete metadata and store it in paths.
// Doesn't fetch full metadata. Blocks until done.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowArtworkWithKnown(
	id uint64, size image.Size, w *work.Work, thumbnailUrl string, paths []string,
) (*work.Work, error) {
	logext.Info("started downloading artwork %v", id)
	panic("unimplemented")
}
