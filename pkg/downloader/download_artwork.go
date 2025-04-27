package downloader

import (
	"fmt"
	"strings"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/downloader/image"
	"github.com/fekoneko/piximan/pkg/encode"
	"github.com/fekoneko/piximan/pkg/fetch"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/fekoneko/piximan/pkg/storage"
)

func (d *Downloader) DownloadArtworkMeta(id uint64, paths []string) (*work.Work, error) {
	logext.Info("started downloading metadata for artwork %v", id)

	w, _, _, err := fetch.ArtworkMeta(d.client, id)
	logext.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	logext.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{}
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.MaybeSuccess(err, "stored metadata for artwork %v in %v", id, paths)
	logext.MaybeError(err, "failed to store metadata for artwork %v", id)
	return w, err
}

func (d *Downloader) DownloadArtwork(id uint64, size image.Size, paths []string) (*work.Work, error) {
	logext.Info("started downloading artwork %v", id)

	w, firstPageUrls, thumbnailUrls, err := fetch.ArtworkMeta(d.client, id)
	logext.MaybeSuccess(err, "fetched metadata for artwork %v", id)
	logext.MaybeError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	if w.Kind == work.KindUgoira {
		err = d.continueUgoira(w, id, paths)
	} else {
		err = d.continueIllustOrManga(w, firstPageUrls, thumbnailUrls, id, size, paths)
	}
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (d *Downloader) continueUgoira(w *work.Work, id uint64, paths []string) error {
	data, frames, err := fetch.ArtworkFrames(d.client, id)
	logext.MaybeSuccess(err, "fetched frames data for artwork %v", id)
	logext.MaybeError(err, "failed to fetch frames data for artwork %v", id)
	if err != nil {
		return err
	}

	archive, err := fetch.Do(d.client, data)
	logext.MaybeSuccess(err, "fetched frames for artwork %v", id)
	logext.MaybeError(err, "failed to fetch frames for artwork %v", id)
	if err != nil {
		return err
	}

	gif, err := encode.GifFromFrames(archive, frames)
	logext.MaybeSuccess(err, "encoded frames for artwork %v", id)
	logext.MaybeError(err, "failed to encode frames for artwork %v", id)
	if err != nil {
		return err
	}

	assets := []storage.Asset{{Bytes: gif, Extension: ".gif"}}
	return storeArtwork(w, id, assets, paths)
}

func (d *Downloader) continueIllustOrManga(
	w *work.Work,
	firstPageUrls *[4]string,
	thumbnailUrls map[uint64]string,
	id uint64,
	size image.Size,
	paths []string,
) error {
	pageUrls, withExtensions, err := inferPages(w, firstPageUrls, thumbnailUrls, size)
	if err == nil {
		assets, err := d.fetchAssets(id, pageUrls, withExtensions, true)
		if err == nil {
			return storeArtwork(w, id, assets, paths)
		}
	}
	logext.Warning("failed to download artwork %v with inferred page urls", id)

	pageUrls, err = d.tryFetchPages(w, id, size)
	if err != nil {
		return err
	}
	assets, err := d.fetchAssets(id, pageUrls, true, false)
	if err != nil {
		return err
	}
	return storeArtwork(w, id, assets, paths)
}

// For illustrations and manga it's possible to omit the request for fetching page urls
// and infer them. This way we can also avoid authorization if the work has age restriction.
// This function receives the data derived from metadata request:
// - urls for different sizes of the first page (available only for works without restriction)
// - thumbnail url (available for all works but cannot be used to derive the extension)
// ! The extension for restricted images in original size cannot be derived, thus we'll have to
// ! try each one later.
func inferPages(
	w *work.Work, firstPageUrls *[4]string, thumbnailUrls map[uint64]string, size image.Size,
) ([]string, bool, error) {
	if firstPageUrls != nil {
		firstPageUrl := (*firstPageUrls)[size]
		if w.NumPages <= 1 {
			return []string{firstPageUrl}, true, nil
		}
		pageUrls, err := inferPagesFromFirstUrl(firstPageUrl, w.NumPages)
		if err == nil {
			return pageUrls, true, nil
		}
	}

	thumbnailUrl, ok := thumbnailUrls[w.Id]
	if !ok {
		return nil, false, fmt.Errorf("cannot find urls to infer from")
	}

	if size == image.SizeThumbnail {
		firstPageUrl := thumbnailUrl
		if w.NumPages <= 1 {
			return []string{firstPageUrl}, true, nil
		}
		pageUrls, err := inferPagesFromFirstUrl(firstPageUrl, w.NumPages)
		if err != nil {
			return nil, false, err
		}
		return pageUrls, true, nil
	}

	const prefix = "https://i.pximg.net/c/250x250_80_a2/img-master/img/"
	const urlDateStart = len(prefix)
	const urlDateEnd = len(prefix) + len("0000/00/00/00/00/00")
	if !strings.HasPrefix(thumbnailUrl, prefix) ||
		len(thumbnailUrl) < urlDateEnd {
		return nil, false, fmt.Errorf("thumbnail url has incorrect format")
	}
	urlDate := thumbnailUrl[urlDateStart:urlDateEnd]

	var firstPageUrl string
	withExtensions := true
	switch size {
	case image.SizeSmall:
		firstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/c/540x540_70/img-master/img/%v/%v_p0_master1200.jpg",
			urlDate, w.Id,
		)
	case image.SizeMedium:
		firstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/img-master/img/%v/%v_p0_master1200.jpg",
			urlDate, w.Id,
		)
	case image.SizeOriginal:
		firstPageUrl = fmt.Sprintf(
			"https://i.pximg.net/img-original/img/%v/%v_p0",
			urlDate, w.Id,
		)
		withExtensions = false
	}
	pageUrls, _ := inferPagesFromFirstUrl(firstPageUrl, w.NumPages)
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

// The function is called if the images were not awailable by inferred page urls.
// Firstly the function will try to make the request without authorization and then with one.
// If the work has age restriction, there's no point in fetching page urls without authorization,
// so unauthoried request will be tried only if session id is unknown, otherwise - skipped.
func (d *Downloader) tryFetchPages(w *work.Work, id uint64, size image.Size) ([]string, error) {
	if w.Restriction == work.RestrictionNone || d.sessionId == nil {
		pageUrls, err := fetch.ArtworkPages(d.client, id, size)
		if err == nil {
			logext.Success("fetched page urls for artwork %v", id)
			return pageUrls, nil
		} else if d.sessionId == nil {
			logext.Error("failed to fetch page urls for artwork %v (authorization could be required): %v", id, err)
			return nil, err
		} else {
			logext.Warning("failed to fetch page urls for artwork %v (authorization could be required): %v", id, err)
		}
	}

	if d.sessionId != nil {
		pageUrls, err := fetch.ArtworkPagesAuthorized(d.client, id, size, *d.sessionId)
		logext.MaybeSuccess(err, "fetched page urls for artwork %v", id)
		logext.MaybeError(err, "failed to fetch page urls for artwork %v", id)
		if err != nil {
			return nil, err
		}
		return pageUrls, nil
	}

	return nil, fmt.Errorf("authorization could be required")
}

var extensions = []string{".jpg", ".png", ".gif"}

// The function fetchest the assets (images) using inferred or fetched page urls.
// If the url was inferred and the extension is not known, the function will try to fetch first
// page with different extensions until it finds the correct one. The list of guessed extensions
// is small and contains only the extensions that Pixiv accepts to be uploaded.
func (d *Downloader) fetchAssets(id uint64, pageUrls []string, withExtensions bool, noLogErrors bool) ([]storage.Asset, error) {
	logErrorOrWarning := logext.Error
	if noLogErrors {
		logErrorOrWarning = logext.Warning
	}
	if len(pageUrls) == 0 {
		err := fmt.Errorf("no pages to download")
		logErrorOrWarning(err.Error())
		return nil, err
	}

	assetChannel := make(chan storage.Asset, len(pageUrls))
	errorChannel := make(chan error)

	guessedExtension := ""
	if !withExtensions {
		for _, extension := range extensions {
			go func() {
				bytes, err := fetch.Do(d.client, pageUrls[0]+extension)
				if err != nil {
					logext.Warning("guessed extension %v was incorrect for artwork %v: %v", extension, id, err)
					errorChannel <- err
				}

				logext.Success("fetched page 1 with guessed extension %v for artwork %v", extension, id)
				assets := storage.Asset{Bytes: bytes, Extension: extension, Page: 1}
				assetChannel <- assets
			}()

			select {
			case <-errorChannel:
			case asset := <-assetChannel:
				assetChannel <- asset
				guessedExtension = extension
			}
			if guessedExtension != "" {
				break
			}
		}
		if guessedExtension == "" {
			return nil, fmt.Errorf("failed to guess extension for artwork %v", id)
		}
	}

	// TODO: The extensions could be different for each page, so we should try to fetch each page
	//       with different extensions, but prefer `guessedExtension`.
	for i, url := range pageUrls {
		if !withExtensions && i == 0 {
			continue
		}
		go func() {
			bytes, err := fetch.Do(d.client, url+guessedExtension)
			if err != nil {
				logErrorOrWarning("failed to fetch page %v for artwork %v: %v", i+1, id, err)
				errorChannel <- err
				return
			}

			logext.Success("fetched page %v for artwork %v", i+1, id)
			var extension = guessedExtension
			if withExtensions {
				dotIndex := strings.LastIndex(url, ".")
				if dotIndex != -1 {
					extension = url[dotIndex:]
				}
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

func storeArtwork(w *work.Work, id uint64, assets []storage.Asset, paths []string) error {
	paths, err := pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.MaybeSuccess(err, "stored files for artwork %v in %v", id, paths)
	logext.MaybeError(err, "failed to store files for artwork %v", id)
	return err
}
