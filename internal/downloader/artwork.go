package downloader

import (
	"fmt"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

// Download only artwork metadata and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using Schedule().
func (d *Downloader) ArtworkMeta(
	id uint64, language work.Language, paths []string,
) (*work.Work, error) {
	return d.internalArtworkMeta(id, language, nil, paths)
}

// Download only artwork metadata and store it in paths when some fields are known in advance.
// Blocks until done. Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using Schedule().
func (d *Downloader) ArtworkMetaWithKnown(
	id uint64, language work.Language, w *work.Work, paths []string,
) (*work.Work, error) {
	return d.internalArtworkMeta(id, language, w, paths)
}

// Will fetch work title and description translations only if they are not present in w.
func (d *Downloader) internalArtworkMeta(
	id uint64, language work.Language, w *work.Work, paths []string,
) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindArtwork, false) || !d.matchArtworkId(id) {
		return nil, ErrSkipped
	}
	d.logger.Info("downloading metadata for artwork %v", id)

	// Omit language if translation is not required to avoid authorization.
	needsTranslation := needsTranslation(w)
	doLanguage := utils.If(needsTranslation, &language, nil)

	fetchedWork, _, _, err := d.artworkMeta(id, nil, doLanguage)
	maybeAddTranslation(needsTranslation, fetchedWork, w)
	if err != nil {
		return nil, err
	} else if !d.matchArtwork(id, fetchedWork, false) {
		return nil, ErrSkipped
	}

	assets := []fsext.Asset{}
	return fetchedWork, d.writeWork(id, queue.ItemKindArtwork, fetchedWork, assets, true, paths)
}

// Doesn't actually make additional requests, but stores incomplete metadata, received earlier.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowArtworkMetaWithKnown(id uint64, w *work.Work, paths []string) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindArtwork, false) || !d.matchArtwork(id, w, true) {
		return nil, ErrSkipped
	}
	assets := []fsext.Asset{}
	return w, d.writeWork(id, queue.ItemKindArtwork, w, assets, true, paths)
}

// Download artwork with all assets and metadata and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using Schedule().
func (d *Downloader) Artwork(
	id uint64, size imageext.Size, language work.Language, paths []string,
) (*work.Work, error) {
	return d.internalArtwork(id, size, language, nil, paths)
}

// Will fetch work title and description translations only if they are not present in w.
func (d *Downloader) internalArtwork(
	id uint64, size imageext.Size, language work.Language, w *work.Work, paths []string,
) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindArtwork, false) || !d.matchArtworkId(id) {
		return nil, ErrSkipped
	}
	d.logger.Info("downloading artwork %v", id)

	// Omit language if translation is not required to avoid authorization.
	needsTranslation := needsTranslation(w)
	doLanguage := utils.If(needsTranslation, &language, nil)

	fetchedWork, firstPageUrl, thumbnailUrl, err := d.artworkMeta(id, &size, doLanguage)
	maybeAddTranslation(needsTranslation, fetchedWork, w)
	if err != nil {
		return nil, err
	} else if !d.matchArtwork(id, fetchedWork, false) {
		return nil, ErrSkipped
	}

	if fetchedWork.Kind == nil {
		err := fmt.Errorf("work kind is missing in %v", fetchedWork)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return fetchedWork, err
	} else if *fetchedWork.Kind == work.KindUgoira {
		asset, err := d.ugoiraAsset(id, fetchedWork)
		assets := []fsext.Asset{*asset}
		d.writeWork(id, queue.ItemKindArtwork, fetchedWork, assets, false, paths)
		return fetchedWork, err
	} else if *fetchedWork.Kind == work.KindIllust || *fetchedWork.Kind == work.KindManga {
		assets, err := d.illustMangaAssets(id, fetchedWork, firstPageUrl, thumbnailUrl, size, false)
		d.writeWork(id, queue.ItemKindArtwork, fetchedWork, assets, false, paths)
		return fetchedWork, err
	} else {
		err := fmt.Errorf("invalid work kind: %v", *fetchedWork.Kind)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return fetchedWork, err
	}
}

// Download artwork with partial metadata known in advance and store it in paths. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// Tries to start downloading assets as soon as possible, but if some rules dependent on full
// metadata are defined, it will wait until full metadata is received.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) ArtworkWithKnown(
	id uint64, size imageext.Size, language work.Language, w *work.Work, thumbnailUrl string, paths []string,
) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindArtwork, false) {
		return nil, ErrSkipped
	} else if matches, needFull := d.matchArtworkNeedFull(id, w); !matches {
		return nil, ErrSkipped
	} else if needFull {
		return d.internalArtwork(id, size, language, w, paths)
	}
	d.logger.Info("downloading artwork %v", id)

	workChannel := make(chan *work.Work, 1)
	assetsChannel := make(chan []fsext.Asset, 1)
	errorChannel := make(chan error, 1)
	var fetchedWork *work.Work
	var assets []fsext.Asset

	// Omit language if translation is not required to avoid authorization.
	needsTranslation := needsTranslation(w)
	doLanguage := utils.If(needsTranslation, &language, nil)

	go d.artworkMetaChannel(id, doLanguage, workChannel, errorChannel)

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

	for range 2 {
		select {
		case fetchedWork = <-workChannel:
		case assets = <-assetsChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}
	maybeAddTranslation(needsTranslation, fetchedWork, w)

	return fetchedWork, d.writeWork(id, queue.ItemKindArtwork, fetchedWork, assets, false, paths)
}

// Download artwork using already available incomplete metadata and store it in paths.
// Doesn't fetch full metadata. Blocks until done.
// Skips downloading if the work doesn't match download rules.
// For downloading multiple works consider using ScheduleWithKnown().
func (d *Downloader) LowArtworkWithKnown(
	id uint64, size imageext.Size, w *work.Work, thumbnailUrl string, paths []string,
) (*work.Work, error) {
	if d.skipped(id, queue.ItemKindArtwork, false) || !d.matchArtwork(id, w, true) {
		return nil, ErrSkipped
	}
	if w.Kind == nil {
		err := fmt.Errorf("work kind is missing in %v", w)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	} else if *w.Kind == work.KindUgoira {
		asset, err := d.ugoiraAsset(id, w)
		assets := []fsext.Asset{*asset}
		d.writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga {
		assets, err := d.illustMangaAssets(id, w, nil, &thumbnailUrl, size, false)
		d.writeWork(id, queue.ItemKindArtwork, w, assets, false, paths)
		return w, err
	} else {
		err := fmt.Errorf("invalid work kind: %v", *w.Kind)
		d.logger.Error("failed to download artwork %v: %v", id, err)
		return w, err
	}
}
