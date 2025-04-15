package downloader

import (
	"fmt"
	"log"
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
	log.Printf("started downloading metadata for artwork %v", id)

	w, _, err := fetch.ArtworkMeta(d.client, id)
	logext.LogIfSuccess(err, "fetched metadata for artwork %v", id)
	logext.LogIfError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{}
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogIfSuccess(err, "stored metadata for artwork %v in %v", id, paths)
	logext.LogIfError(err, "failed to store metadata for artwork %v", id)
	return w, err
}

func (d *Downloader) DownloadArtwork(id uint64, size image.Size, paths []string) (*work.Work, error) {
	log.Printf("started downloading artwork %v", id)

	w, firstPageUrls, err := fetch.ArtworkMeta(d.client, id)
	logext.LogIfSuccess(err, "fetched metadata for artwork %v", id)
	logext.LogIfError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	if w.Kind == work.KindUgoira {
		err = d.continueUgoira(w, id, paths)
	} else {
		err = d.continueIllustOrManga(w, firstPageUrls, id, size, paths)
	}
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (d *Downloader) continueUgoira(w *work.Work, id uint64, paths []string) error {
	data, frames, err := fetch.ArtworkFrames(d.client, id)
	logext.LogIfSuccess(err, "fetched frames data for artwork %v", id)
	logext.LogIfError(err, "failed to fetch frames data for artwork %v", id)
	if err != nil {
		return err
	}

	archive, err := fetch.Do(d.client, data)
	logext.LogIfSuccess(err, "fetched frames for artwork %v", id)
	logext.LogIfError(err, "failed to fetch frames for artwork %v", id)
	if err != nil {
		return err
	}

	gif, err := encode.GifFromFrames(archive, frames)
	logext.LogIfSuccess(err, "encoded frames for artwork %v", id)
	logext.LogIfError(err, "failed to encode frames for artwork %v", id)
	if err != nil {
		return err
	}

	assets := []storage.Asset{{Bytes: gif, Extension: ".gif"}}
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogIfSuccess(err, "stored files for artwork %v in %v", id, paths)
	logext.LogIfError(err, "failed to store files for artwork %v", id)
	return err
}

func (d *Downloader) continueIllustOrManga(
	w *work.Work, firstPageUrls *[4]string, id uint64, size image.Size, paths []string,
) error {
	// TODO: it seems like for new works the `urls` contains null values and
	// `.../pages` doesn't  work - need another method
	if firstPageUrls != nil {
		pageUrls := inferPages((*firstPageUrls)[size], w.NumPages)
		if err := d.continueArtworkWithPages(w, pageUrls, id, paths); err == nil {
			return nil
		}
	}
	logext.LogError("failed to infer page urls for artwork %v - trying to fetch them", id)

	pageUrls, err := fetch.ArtworkPages(d.client, id, size)
	logext.LogIfSuccess(err, "fetched page urls for artwork %v", id)
	logext.LogIfError(err, "failed to fetch page urls for artwork %v", id)
	if err != nil {
		return err
	}
	return d.continueArtworkWithPages(w, pageUrls, id, paths)
}

func (d *Downloader) continueArtworkWithPages(
	w *work.Work, pageUrls []string, id uint64, paths []string,
) error {
	assetChannel := make(chan storage.Asset, len(pageUrls))
	errorChannel := make(chan error, len(pageUrls))
	indexChannel := make(chan int, len(pageUrls))

	for i, url := range pageUrls {
		go func() {
			bytes, err := fetch.Do(d.client, url)
			logext.LogIfSuccess(err, "fetched page %v for artwork %v", i, id)
			logext.LogIfError(err, "failed to fetch page %v for artwork %v", i, id)

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
			return err
		}
	}

	paths, err := pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogIfSuccess(err, "stored files for artwork %v in %v", id, paths)
	logext.LogIfError(err, "failed to store files for artwork %v", id)
	return err
}

func inferPages(firstPageUrl string, numPages uint64) []string {
	p0Index := strings.Index(firstPageUrl, "p0")
	if p0Index == -1 {
		return []string{firstPageUrl}
	}

	pageUrls := make([]string, numPages)
	pageUrls[0] = firstPageUrl
	for i := uint64(1); i < numPages; i++ {
		pageUrls[i] = fmt.Sprintf("%vp%v%v", firstPageUrl[:p0Index], i, firstPageUrl[p0Index+2:])
	}

	return pageUrls
}
