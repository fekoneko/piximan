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

// TODO: test R-18(G) without authorization
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
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.MaybeSuccess(err, "stored files for artwork %v in %v", id, paths)
	logext.MaybeError(err, "failed to store files for artwork %v", id)
	return err
}

func (d *Downloader) continueIllustOrManga(
	w *work.Work,
	firstPageUrls *[4]string,
	thumbnailUrls map[uint64]string,
	id uint64,
	size image.Size,
	paths []string,
) error {
	// TODO: try infering first only if:
	//       - R-18(G)
	//       - no restriction and not original size
	pageUrls, withExtensions, err := inferPages(w, firstPageUrls, thumbnailUrls, size)
	if err == nil {
		assets, err := d.fetchAssets(id, pageUrls, withExtensions, true)
		if err == nil {
			return storeWork(w, assets, paths)
		}
	}
	logext.Warning("failed to download artwork %v with inferred page urls", id)

	// TODO: - if inferring failed and R-18(G) - this must be done with authorization
	//       - if inferring failed and no restriction - try without and then with authorization
	pageUrls, err = fetch.ArtworkPages(d.client, id, size)
	logext.MaybeSuccess(err, "fetched page urls for artwork %v", id)
	logext.MaybeError(err, "failed to fetch page urls for artwork %v", id)
	if err != nil {
		return err
	}
	assets, err := d.fetchAssets(id, pageUrls, true, false)
	if err != nil {
		return err
	}
	return storeWork(w, assets, paths)
}

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

var extensions = []string{".jpg", ".png"} // TODO: add more

func (d *Downloader) fetchAssets(id uint64, pageUrls []string, withExtensions bool, noLogErrors bool) ([]storage.Asset, error) {
	if len(pageUrls) == 0 {
		err := fmt.Errorf("no pages to download")
		logext.Error(err.Error())
		return nil, err
	}

	assetChannel := make(chan storage.Asset, len(pageUrls))
	errorChannel := make(chan error, len(pageUrls))
	indexChannel := make(chan int, len(pageUrls))

	urlSuffix := ""
	if !withExtensions {
		for _, extension := range extensions {
			go func() {
				bytes, err := fetch.Do(d.client, pageUrls[0]+extension)
				logext.MaybeSuccess(err, "fetched page 1 with guessed extension %v for artwork %v", extension, id)
				logext.MaybeWarning(err, "guessed extension %v was incorrect for artwork %v", extension, id)

				assets := storage.Asset{Bytes: bytes, Extension: extension}
				assetChannel <- assets
				errorChannel <- err
			}()

			if <-errorChannel == nil {
				errorChannel <- nil
				indexChannel <- 0
				urlSuffix = extension
				break
			}
			<-assetChannel
		}
		if urlSuffix == "" {
			return nil, fmt.Errorf("failed to guess extension for artwork %v", id)
		}
	}

	for i, url := range pageUrls {
		if !withExtensions && i == 0 {
			continue
		}
		go func() {
			bytes, err := fetch.Do(d.client, url+urlSuffix)
			logext.MaybeSuccess(err, "fetched page %v for artwork %v", i+1, id)
			if noLogErrors {
				logext.MaybeWarning(err, "failed to fetch page %v for artwork %v", i+1, id)
			} else {
				logext.MaybeError(err, "failed to fetch page %v for artwork %v", i+1, id)
			}

			dotIndex := strings.LastIndex(url, ".")
			var extension string
			if dotIndex != -1 {
				extension = url[dotIndex:]
			}
			assets := storage.Asset{Bytes: bytes, Extension: extension}
			assetChannel <- assets
			errorChannel <- err
			indexChannel <- i
		}()
	}

	assets := make([]storage.Asset, len(pageUrls))
	for range pageUrls {
		i := <-indexChannel
		assets[i] = <-assetChannel
		err := <-errorChannel
		if err != nil {
			return nil, err
		}
	}

	return assets, nil
}

func storeWork(w *work.Work, assets []storage.Asset, paths []string) error {
	paths, err := pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.MaybeSuccess(err, "stored files for artwork %v in %v", w.Id, paths)
	logext.MaybeError(err, "failed to store files for artwork %v", w.Id)
	return err
}
