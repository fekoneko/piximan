package downloader

import (
	"fmt"
	"path"
	"strings"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/fsext"
	"github.com/fekoneko/piximan/internal/imageext"
	"github.com/fekoneko/piximan/internal/utils"
)

// Fetch all image assets for illust or manga artwork
func (d *Downloader) illustMangaAssets(
	id uint64, w *work.Work, firstPageUrl *string, thumbnailUrl *string,
	size imageext.Size, onlyFirstPage bool,
) ([]fsext.Asset, error) {
	numPages := utils.If(onlyFirstPage, utils.ToPtr(uint64(1)), w.NumPages)
	pageUrls, withExtensions, err := inferPages(id, firstPageUrl, thumbnailUrl, size, numPages)
	if err != nil {
		d.logger.Warning("failed to infer page urls for artwork %v: %v", id, err)
	} else {
		assets, err := d.fetchAssets(id, pageUrls, withExtensions, true)
		if err == nil {
			return assets, nil
		}
	}

	pageUrls, err = d.fetchPages(w, id, size)
	if err != nil {
		return nil, err
	}
	if onlyFirstPage {
		pageUrls = pageUrls[:1]
	}
	assets, err := d.fetchAssets(id, pageUrls, true, false)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// For illustrations and manga it's possible to omit the request for fetching page urls
// and infer them. This way we can also avoid authorization if the work has age restriction.
// This function receives the data derived from metadata request:
// - url of the first page (available only for works without restriction)
// - thumbnail url (available for all works but cannot be used to derive the extension)
// If provided, firstPageUrl and thumbnailUls should be of the provided size.
// The extension for restricted images in original size cannot be derived, thus we'll have to
// try each one later.
func inferPages(
	id uint64, firstPageUrl *string, thumbnailUrl *string, size imageext.Size, numPages *uint64,
) (pageUrls []string, withExtensions bool, err error) {
	if numPages == nil {
		err := fmt.Errorf("page count is missing")
		return nil, false, err
	} else if *numPages == 0 {
		return nil, false, fmt.Errorf("page count is zero")
	}

	if firstPageUrl != nil {
		if *numPages == 1 {
			return []string{*firstPageUrl}, true, nil
		}
		pageUrls, err := inferPagesFromFirstUrl(*firstPageUrl, *numPages)
		if err == nil {
			return pageUrls, true, nil
		}
	}

	if thumbnailUrl == nil {
		return nil, false, fmt.Errorf("no thumbnail url to infer from")
	}

	const prefixMaster = "https://i.pximg.net/c/250x250_80_a2/img-master/img/"
	const prefixCustom = "https://i.pximg.net/c/250x250_80_a2/custom-thumb/img/"
	var prefixLength int
	if strings.HasPrefix(*thumbnailUrl, prefixMaster) {
		prefixLength = len(prefixMaster)
	} else if strings.HasPrefix(*thumbnailUrl, prefixCustom) {
		prefixLength = len(prefixCustom)
	} else {
		return nil, false, fmt.Errorf("thumbnail url has incorrect prefix")
	}

	urlDateStart := prefixLength
	urlDateEnd := prefixLength + len("0000/00/00/00/00/00")
	if len(*thumbnailUrl) < urlDateEnd {
		return nil, false, fmt.Errorf("thumbnail url is too short")
	}
	urlDate := (*thumbnailUrl)[urlDateStart:urlDateEnd]

	var inferredFirstPageUrl string
	withExtensions = true
	switch size {
	case imageext.SizeThumbnail:
		inferredFirstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/c/128x128/img-master/img/%v/%v_p0_square1200.jpg",
			urlDate, id,
		)
	case imageext.SizeSmall:
		inferredFirstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/c/540x540_70/img-master/img/%v/%v_p0_master1200.jpg",
			urlDate, id,
		)
	case imageext.SizeMedium:
		inferredFirstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/img-master/img/%v/%v_p0_master1200.jpg",
			urlDate, id,
		)
	case imageext.SizeOriginal:
		inferredFirstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/img-original/img/%v/%v_p0",
			urlDate, id,
		)
		withExtensions = false
	}
	pageUrls, _ = inferPagesFromFirstUrl(inferredFirstPageUrl, *numPages)
	return pageUrls, withExtensions, nil
}

func inferPagesFromFirstUrl(firstPageUrl string, numPages uint64) ([]string, error) {
	p0Index := strings.Index(firstPageUrl, "p0")
	if p0Index == -1 {
		return nil, fmt.Errorf("cannot find 'p0' in url")
	}
	pageUrls := make([]string, numPages)
	pageUrls[0] = firstPageUrl
	for i := uint64(1); i < numPages; i++ {
		pageUrls[i] = fmt.Sprintf("%vp%v%v", firstPageUrl[:p0Index], i, firstPageUrl[p0Index+2:])
	}

	return pageUrls, nil
}

// The function is called if the images were not available by inferred page urls.
// First the function will try to make the request without authorization and then with one.
// If the work has age restriction, there's no point in fetching page urls without authorization,
// so unauthoried request will be tried only if session id is unknown, otherwise - skipped.
func (d *Downloader) fetchPages(w *work.Work, id uint64, size imageext.Size) ([]string, error) {
	authorized := d.client.Authorized()
	withUnauthorized := w.Restriction == nil ||
		*w.Restriction == work.RestrictionNone || !authorized
	if withUnauthorized {
		pageUrls, err := d.client.IllustMangaPages(id, size)
		if err == nil {
			d.logger.Success("fetched page urls for artwork %v", id)
			return pageUrls, nil
		} else if !authorized {
			d.logger.Error("failed to fetch page urls for artwork %v (authorization could be required): %v", id, err)
			return nil, err
		} else {
			d.logger.Warning("failed to fetch page urls for artwork %v (authorization could be required): %v", id, err)
		}
	}

	if authorized {
		if withUnauthorized {
			d.logger.Info("retrying fetching pages with authorization for artwork %v", id)
		}
		pageUrls, err := d.client.IllustMangaPagesAuthorized(id, size)
		d.logger.MaybeSuccess(err, "fetched page urls for artwork %v", id)
		d.logger.MaybeError(err, "failed to fetch page urls for artwork %v", id)
		if err != nil {
			return nil, err
		}
		return pageUrls, nil
	}

	err := fmt.Errorf("authorization could be required")
	d.logger.Error("failed to fetch page urls for artwork %v: %v", id, err)
	return nil, err
}

var extensions = []string{".jpg", ".png", ".gif"}

// The function fetches the assets (images) using inferred or fetched page urls.
// If the url was inferred and the extension is not known, the function will try to fetch first
// page with different extensions until it finds the correct one. The list of guessed extensions
// is small and contains only the extensions that Pixiv accepts to be uploaded.
// Work cannot have different extensions for different pages as Pixiv does not allow it.
func (d *Downloader) fetchAssets(
	id uint64, pageUrls []string, withExtensions bool, noLogErrors bool,
) ([]fsext.Asset, error) {
	logErrorOrWarning := d.logger.Error
	if noLogErrors {
		logErrorOrWarning = d.logger.Warning
	}
	if len(pageUrls) == 0 {
		err := fmt.Errorf("no pages to download")
		logErrorOrWarning("failed to fetch assets for artwork %v: %v", id, err)
		return nil, err
	}

	assetsChannel := make(chan fsext.Asset, 1)
	errorChannel := make(chan error, 1)

	guessedExtension := ""
	if !withExtensions {
		for _, extension := range extensions {
			bytes, _, err := d.client.Do(pageUrls[0]+extension, nil)
			if err != nil {
				d.logger.Info("guessed extension %v was incorrect for artwork %v: %v", extension, id, err)
				continue
			}

			d.logger.Success("fetched page 1 with guessed extension %v for artwork %v", extension, id)
			name := fsext.IllustMangaAssetName(1, extension)
			assets := fsext.Asset{Bytes: bytes, Name: name}
			assetsChannel <- assets
			guessedExtension = extension
			break
		}
		if guessedExtension == "" {
			err := fmt.Errorf("all tried extensions were incorrect")
			logErrorOrWarning("failed to guess extension for artwork %v: %v", id, err)
			return nil, err
		}
	}

	for i, url := range pageUrls {
		if !withExtensions && i == 0 {
			continue
		}
		go func() {
			bytes, _, err := d.client.Do(url+guessedExtension, nil)
			if err != nil {
				logErrorOrWarning("failed to fetch page %v for artwork %v: %v", i+1, id, err)
				errorChannel <- err
				return
			}

			d.logger.Success("fetched page %v for artwork %v", i+1, id)
			var extension = guessedExtension
			if withExtensions {
				extension = path.Ext(url)
			}
			name := fsext.IllustMangaAssetName(uint64(i+1), extension)
			assets := fsext.Asset{Bytes: bytes, Name: name}
			assetsChannel <- assets
		}()
	}

	assets := make([]fsext.Asset, len(pageUrls))
	for i := range pageUrls {
		select {
		case assets[i] = <-assetsChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	return assets, nil
}

func (d *Downloader) illustMangaAssetsChannel(
	id uint64, w *work.Work, firstPageUrl *string, thumbnailUrl *string, size imageext.Size,
	assetsChannel chan []fsext.Asset, errorChannel chan error,
) {
	if assets, err := d.illustMangaAssets(id, w, firstPageUrl, thumbnailUrl, size, false); err == nil {
		assetsChannel <- assets
	} else {
		errorChannel <- err
	}
}
