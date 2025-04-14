package downloader

import (
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

	w, err := fetch.ArtworkMeta(d.client, id)
	logext.LogSuccess(err, "fetched metadata for artwork %v", id)
	logext.LogError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{}
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogSuccess(err, "stored metadata for artwork %v in %v", id, paths)
	logext.LogError(err, "failed to store metadata for artwork %v", id)
	return w, err
}

func (d *Downloader) DownloadArtwork(id uint64, size image.Size, paths []string) (*work.Work, error) {
	log.Printf("started downloading artwork %v", id)

	w, err := fetch.ArtworkMeta(d.client, id)
	logext.LogSuccess(err, "fetched metadata for artwork %v", id)
	logext.LogError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	if w.Kind == work.KindUgoira {
		err = d.continueUgoira(w, id, paths)
	} else {
		err = d.continueIllustOrManga(w, id, size, paths)
	}
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (d *Downloader) continueUgoira(w *work.Work, id uint64, paths []string) error {
	data, frames, err := fetch.ArtworkFrames(d.client, id)
	logext.LogSuccess(err, "fetched frames data for artwork %v", id)
	logext.LogError(err, "failed to fetch frames data for artwork %v", id)
	if err != nil {
		return err
	}

	archive, err := fetch.Do(d.client, data)
	logext.LogSuccess(err, "fetched frames for artwork %v", id)
	logext.LogError(err, "failed to fetch frames for artwork %v", id)
	if err != nil {
		return err
	}

	gif, err := encode.GifFromFrames(archive, frames)
	logext.LogSuccess(err, "encoded frames for artwork %v", id)
	logext.LogError(err, "failed to encode frames for artwork %v", id)
	if err != nil {
		return err
	}

	assets := []storage.Asset{{Bytes: gif, Extension: ".gif"}}
	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogSuccess(err, "stored files for artwork %v in %v", id, paths)
	logext.LogError(err, "failed to store files for artwork %v", id)
	return err
}

func (d *Downloader) continueIllustOrManga(
	w *work.Work, id uint64, size image.Size, paths []string,
) error {
	pages, err := fetch.ArtworkUrls(d.client, id)
	logext.LogSuccess(err, "fetched page urls for artwork %v", id)
	logext.LogError(err, "failed to fetch page urls for artwork %v", id)
	if err != nil {
		return err
	}

	assetChannel := make(chan storage.Asset, len(pages))
	errorChannel := make(chan error, len(pages))
	indexChannel := make(chan int, len(pages))
	for i, urls := range pages {
		go func() {
			url := urls[size]
			bytes, err := fetch.Do(d.client, url)
			logext.LogSuccess(err, "fetched page %v for artwork %v", i, id)
			logext.LogError(err, "failed to fetch page %v for artwork %v", i, id)
			dotIndex := strings.LastIndex(url, ".")
			var extension string
			if dotIndex != -1 {
				extension = urls[size][dotIndex:]
			}
			assets := storage.Asset{Bytes: bytes, Extension: extension}
			assetChannel <- assets
			errorChannel <- err
			indexChannel <- i
		}()
	}

	assets := make([]storage.Asset, len(pages))
	for range pages {
		i := <-indexChannel
		assets[i] = <-assetChannel
		err = <-errorChannel
		if err != nil {
			return err
		}
	}

	paths, err = pathext.FormatWorkPaths(paths, w)
	if err == nil {
		err = storage.StoreWork(w, assets, paths)
	}
	logext.LogSuccess(err, "stored files for artwork %v in %v", id, paths)
	logext.LogError(err, "failed to store files for artwork %v", id)
	return err
}
