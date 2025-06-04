package downloader

import (
	"fmt"
	"path"
	"strings"

	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/image"
	"github.com/fekoneko/piximan/internal/fetch"
	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/storage"
)

func urlFromMap(id uint64, urls map[uint64]string) *string {
	var urlPtr *string
	if url, ok := urls[id]; ok {
		urlPtr = &url
	}
	return urlPtr
}

// Fetch all image assets for illust or manga artwork
func (d *Downloader) illustMangaAssets(
	id uint64,
	w *work.Work,
	firstPageUrls *[4]string,
	thumbnailUrl *string,
	size image.Size,
) ([]storage.Asset, error) {
	pageUrls, withExtensions, err := inferPages(id, w, firstPageUrls, thumbnailUrl, size)
	if err != nil {
		logext.Warning("failed to infer page urls for artwork %v: %v", id, err)
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
	assets, err := d.fetchAssets(id, pageUrls, true, false)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// For illustrations and manga it's possible to omit the request for fetching page urls
// and infer them. This way we can also avoid authorization if the work has age restriction.
// This function receives the data derived from metadata request:
// - urls for different sizes of the first page (available only for works without restriction)
// - thumbnail url (available for all works but cannot be used to derive the extension)
// ! The extension for restricted images in original size cannot be derived, thus we'll have to
// ! try each one later.
func inferPages(
	id uint64, w *work.Work, firstPageUrls *[4]string, thumbnailUrl *string, size image.Size,
) ([]string, bool, error) {
	if w.NumPages == nil {
		err := fmt.Errorf("number of pages is missing in %v", w)
		return nil, false, err
	}

	if firstPageUrls != nil {
		firstPageUrl := (*firstPageUrls)[size]
		if *w.NumPages <= 1 {
			return []string{firstPageUrl}, true, nil
		}
		pageUrls, err := inferPagesFromFirstUrl(firstPageUrl, *w.NumPages)
		if err == nil {
			return pageUrls, true, nil
		}
	}

	if thumbnailUrl == nil {
		return nil, false, fmt.Errorf("no thumbnail url to infer from")
	}

	const PREFIX_MASTER = "https://i.pximg.net/c/250x250_80_a2/img-master/img/"
	const PREFIX_CUSTOM = "https://i.pximg.net/c/250x250_80_a2/custom-thumb/img/"
	var prefixLength int
	if strings.HasPrefix(*thumbnailUrl, PREFIX_MASTER) {
		prefixLength = len(PREFIX_MASTER)
	} else if strings.HasPrefix(*thumbnailUrl, PREFIX_CUSTOM) {
		prefixLength = len(PREFIX_CUSTOM)
	} else {
		return nil, false, fmt.Errorf("thumbnail url has incorrect prefix")
	}

	urlDateStart := prefixLength
	urlDateEnd := prefixLength + len("0000/00/00/00/00/00")
	if len(*thumbnailUrl) < urlDateEnd {
		return nil, false, fmt.Errorf("thumbnail url is too short")
	}
	urlDate := (*thumbnailUrl)[urlDateStart:urlDateEnd]

	var firstPageUrl string
	withExtensions := true
	switch size {
	case image.SizeThumbnail:
		firstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/c/128x128/img-master/img/%v/%v_p0_square1200.jpg",
			urlDate, id,
		)
	case image.SizeSmall:
		firstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/c/540x540_70/img-master/img/%v/%v_p0_master1200.jpg",
			urlDate, id,
		)
	case image.SizeMedium:
		firstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/img-master/img/%v/%v_p0_master1200.jpg",
			urlDate, id,
		)
	case image.SizeOriginal:
		firstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/img-original/img/%v/%v_p0",
			urlDate, id,
		)
		withExtensions = false
	}
	pageUrls, _ := inferPagesFromFirstUrl(firstPageUrl, *w.NumPages)
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
func (d *Downloader) fetchPages(w *work.Work, id uint64, size image.Size) ([]string, error) {
	sessionId, withSessionId := d.sessionId()
	withUnauthorized := w.Restriction == nil ||
		*w.Restriction == work.RestrictionNone || !withSessionId
	if withUnauthorized {
		pageUrls, err := fetch.ArtworkPages(d.client(), id, size)
		if err == nil {
			logext.Success("fetched page urls for artwork %v", id)
			return pageUrls, nil
		} else if !withSessionId {
			logext.Error("failed to fetch page urls for artwork %v (authorization could be required): %v", id, err)
			return nil, err
		} else {
			logext.Warning("failed to fetch page urls for artwork %v (authorization could be required): %v", id, err)
		}
	}

	if withSessionId {
		if withUnauthorized {
			logext.Info("retrying fetching pages with authorization for artwork %v", id)
		}
		pageUrls, err := fetch.ArtworkPagesAuthorized(d.client(), id, size, *sessionId)
		logext.MaybeSuccess(err, "fetched page urls for artwork %v", id)
		logext.MaybeError(err, "failed to fetch page urls for artwork %v", id)
		if err != nil {
			return nil, err
		}
		return pageUrls, nil
	}

	err := fmt.Errorf("authorization could be required")
	logext.Error("failed to fetch page urls for artwork %v: %v", id, err)
	return nil, err
}

var extensions = []string{".jpg", ".png", ".gif"}

// The function fetches the assets (images) using inferred or fetched page urls.
// If the url was inferred and the extension is not known, the function will try to fetch first
// page with different extensions until it finds the correct one. The list of guessed extensions
// is small and contains only the extensions that Pixiv accepts to be uploaded.
// Work cannot have different extensions for different pages as Pixiv does not allow it.
func (d *Downloader) fetchAssets(id uint64, pageUrls []string, withExtensions bool, noLogErrors bool) ([]storage.Asset, error) {
	logErrorOrWarning := logext.Error
	if noLogErrors {
		logErrorOrWarning = logext.Warning
	}
	if len(pageUrls) == 0 {
		err := fmt.Errorf("no pages to download")
		logErrorOrWarning("failed to fetch assets for artwork %v: %v", id, err)
		return nil, err
	}

	assetChannel := make(chan storage.Asset, len(pageUrls))
	errorChannel := make(chan error)

	guessedExtension := ""
	if !withExtensions {
		for _, extension := range extensions {
			bytes, _, err := fetch.Do(d.client(), pageUrls[0]+extension, nil)
			if err != nil {
				logext.Info("guessed extension %v was incorrect for artwork %v: %v", extension, id, err)
				continue
			}

			logext.Success("fetched page 1 with guessed extension %v for artwork %v", extension, id)
			assets := storage.Asset{Bytes: bytes, Extension: extension, Page: 1}
			assetChannel <- assets
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
			bytes, _, err := fetch.Do(d.client(), url+guessedExtension, nil)
			if err != nil {
				logErrorOrWarning("failed to fetch page %v for artwork %v: %v", i+1, id, err)
				errorChannel <- err
				return
			}

			logext.Success("fetched page %v for artwork %v", i+1, id)
			var extension = guessedExtension
			if withExtensions {
				extension = path.Ext(url)
			}
			assets := storage.Asset{Bytes: bytes, Extension: extension, Page: uint64(i + 1)}
			assetChannel <- assets
		}()
	}

	assets := make([]storage.Asset, len(pageUrls))
	for i := range pageUrls {
		select {
		case assets[i] = <-assetChannel:
		case err := <-errorChannel:
			return nil, err
		}
	}

	return assets, nil
}

func (d *Downloader) illustMangaAssetsChannel(
	id uint64,
	w *work.Work,
	firstPageUrls *[4]string,
	thumbnailUrl *string,
	size image.Size,
	assetsChannel chan []storage.Asset,
	errorChannel chan error,
) {
	if assets, err := d.illustMangaAssets(id, w, firstPageUrls, thumbnailUrl, size); err == nil {
		assetsChannel <- assets
	} else {
		errorChannel <- err
	}
}
