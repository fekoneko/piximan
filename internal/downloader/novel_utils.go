package downloader

import (
	"fmt"
	"path"
	"sync"

	"github.com/fekoneko/piximan/internal/client/dto"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

// TODO: don't forget to change the docs

// Fetch novel metadata, cover url and information about embedded illustrations.
// Retry authorized if something is missing.
func (d *Downloader) novelMeta(
	id uint64, size *imageext.Size,
) (
	w *work.Work, coverUrl *string, upladedImages dto.NovelUpladedImages,
	pixivImages dto.NovelPixivImages, pages dto.NovelPages, err error,
) {
	authorized := d.client.Authorized()
	do := d.client.NovelMeta
	logError := utils.If(authorized, d.logger.Warning, d.logger.Error)
	triedAuthorized := false

	for {
		w, coverUrl, upladedImages, pixivImages, pages, withPages, err := do(id, size)
		d.logger.MaybeSuccess(err, "fetched metadata for novel %v", id)
		if err != nil {
			logError("failed to fetch metadata for novel %v: %v", id, err)
		} else if !withPages {
			err = fmt.Errorf("pages are missing")
			logError("failed to download novel %v: %v", id, err)
		} else if coverUrl == nil {
			err = fmt.Errorf("cover url is missing")
			logError("failed to download novel %v: %v", id, err)
		} else {
			if !w.Full() {
				d.logger.Warning("metadata for novel %v is incomplete", id)
			}
			return w, coverUrl, upladedImages, pixivImages, pages, nil
		}

		if triedAuthorized || !authorized {
			return nil, nil, nil, nil, nil, err
		}
		d.logger.Info("retrying fetching metadata with authorization for novel %v", id)
		do = d.client.NovelMetaAuthorized
		logError = d.logger.Error
		triedAuthorized = true
	}
}

// Fetch novel metadata and ignore if anything else is missing.
func (d *Downloader) novelOnlyMeta(id uint64) (*work.Work, error) {
	w, _, _, _, _, _, err := d.client.NovelMeta(id, nil)
	d.logger.MaybeSuccess(err, "fetched metadata for novel %v", id)
	d.logger.MaybeError(err, "failed to fetch metadata for novel %v", id)
	if !w.Full() {
		d.logger.Warning("metadata for novel %v is incomplete", id)
	}
	return w, err
}

// Fetch novel cover asset.
func (d *Downloader) novelCoverAsset(id uint64, coverUrl string) (*fsext.Asset, error) {
	cover, _, err := d.client.Do(coverUrl, nil)
	d.logger.MaybeSuccess(err, "fetched cover for novel %v", id)
	d.logger.MaybeError(err, "failed to fetch cover for novel %v", id)

	name := fsext.NovelCoverAssetName(path.Ext(coverUrl))
	asset := fsext.Asset{Bytes: cover, Name: name}
	return &asset, nil
}

// Fetch all novel illustrations as assets. uploadedImages urls should have the correct size.
func (d *Downloader) novelImageAssets(
	id uint64, size imageext.Size,
	uploadedImages dto.NovelUpladedImages, pixivImages dto.NovelPixivImages,
) (map[int]fsext.Asset, error) {
	assets := make(map[int]fsext.Asset, len(uploadedImages)+len(pixivImages))
	assetsMutex := sync.Mutex{}
	errorChannel := make(chan error, 1)

	for index, url := range uploadedImages {
		go func() {
			if bytes, _, err := d.client.Do(url, nil); err == nil {
				d.logger.Success("fetched illustration %v for novel %v", index, id)
				name := fsext.NovelImageAssetName(index, path.Ext(url))
				asset := fsext.Asset{Bytes: bytes, Name: name}
				assetsMutex.Lock()
				assets[index] = asset
				assetsMutex.Unlock()
				errorChannel <- nil
			} else {
				d.logger.Error("failed to fetch illustration %v for novel %v: %v", index, id, err)
				errorChannel <- err
			}
		}()
	}

	for index, artworkId := range pixivImages {
		go func() {
			if asset, err := d.novelPixivImage(id, index, artworkId, size); err == nil {
				assetsMutex.Lock()
				assets[index] = *asset
				assetsMutex.Unlock()
				errorChannel <- nil
			} else {
				errorChannel <- err
			}
		}()
	}

	for range len(uploadedImages) + len(pixivImages) {
		if err := <-errorChannel; err != nil {
			return nil, err
		}
	}

	return assets, nil
}

// coverAsset() but returs results through channels.
func (d *Downloader) novelCoverAssetChannel(
	id uint64, coverUrl string, coverChannel chan *fsext.Asset, errorChannel chan error,
) {
	if coverAsset, err := d.novelCoverAsset(id, coverUrl); err == nil {
		coverChannel <- coverAsset
	} else {
		errorChannel <- err
	}
}

// novelImageAssets() but returs results through channels.
func (d *Downloader) novelImageAssetsChannel(
	id uint64, size imageext.Size, uploadedImages dto.NovelUpladedImages,
	pixivImages dto.NovelPixivImages, imagesChannel chan map[int]fsext.Asset,
	errorChannel chan error,
) {
	if assets, err := d.novelImageAssets(id, size, uploadedImages, pixivImages); err == nil {
		imagesChannel <- assets
	} else {
		errorChannel <- err
	}
}

// novelMeta() + novelImageAssets() but returs results through channels.
func (d *Downloader) novelMetaImageAssetsChannel(
	id uint64, size imageext.Size, workChannel chan *work.Work,
	pagesChannel chan dto.NovelPages, imagesChannel chan map[int]fsext.Asset, errorChannel chan error,
) {
	if w, _, uploadedImages, pixivImages, pages, err := d.novelMeta(id, &size); err != nil {
		errorChannel <- err
	} else if assets, err := d.novelImageAssets(id, size, uploadedImages, pixivImages); err != nil {
		errorChannel <- err
	} else {
		workChannel <- w
		pagesChannel <- pages
		imagesChannel <- assets
	}
}

func combineAssets(
	coverAsset *fsext.Asset, imageAssets map[int]fsext.Asset, pages dto.NovelPages,
) []fsext.Asset {
	imageName := func(index int) string {
		return imageAssets[index].Name
	}
	pageName := func(page int) string {
		return fsext.NovelPageAssetName(page)
	}
	pageAssets := pages(imageName, pageName)

	imageAssetsSlice := make([]fsext.Asset, 0, len(imageAssets))
	for _, asset := range imageAssets {
		imageAssetsSlice = append(imageAssetsSlice, asset)
	}

	assets := make([]fsext.Asset, 0, len(pageAssets)+len(imageAssetsSlice)+2)
	assets = append(assets, pageAssets...)
	assets = append(assets, imageAssetsSlice...)
	assets = append(assets, *coverAsset)

	return assets
}

// Get only the first page of artwork. Used for novel embedded illustrations.
// Doesn't store anything, just returns the work and the asset.
// This operation cannot be ignored with download rules or ignore list.
func (d *Downloader) novelPixivImage(
	novelId uint64, imageIndex int, artworkId uint64, size imageext.Size,
) (*fsext.Asset, error) {
	d.logger.Info("getting artwork %v for illustration %v in novel %v", artworkId, imageIndex, novelId)

	w, firstPageUrl, thumbnailUrl, err := d.artworkMeta(artworkId, &size)
	if err != nil {
		return nil, err
	}

	var asset *fsext.Asset
	if w.Kind == nil {
		err = fmt.Errorf("work kind is missing in %v", w)
		d.logger.Error("failed to download artwork %v: %v", artworkId, err)
	} else if *w.Kind == work.KindUgoira {
		asset, err = d.ugoiraAsset(artworkId, w)
	} else if *w.Kind == work.KindIllust || *w.Kind == work.KindManga {
		var assets []fsext.Asset
		assets, err = d.illustMangaAssets(artworkId, w, firstPageUrl, thumbnailUrl, size, true)
		if err == nil && len(assets) > 0 {
			asset = &assets[0]
		}
	} else {
		err = fmt.Errorf("invalid work kind: %v", *w.Kind)
		d.logger.Error("failed to download artwork %v: %v", artworkId, err)
	}

	if err != nil {
		return nil, err
	} else {
		ext := path.Ext(asset.Name)
		asset.Name = fsext.NovelImageAssetName(imageIndex, ext)
		return asset, nil
	}
}
