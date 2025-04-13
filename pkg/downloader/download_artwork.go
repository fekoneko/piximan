package downloader

import (
	"strings"

	"github.com/fekoneko/piximan/pkg/encode"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/pathext"
	"github.com/fekoneko/piximan/pkg/storage"
	"github.com/fekoneko/piximan/pkg/work"
)

func (d *Downloader) DownloadArtworkMeta(id uint64, path string) (*work.Work, error) {
	work, err := d.fetchArtworkMeta(id)
	logext.LogSuccess(err, "fetched metadata for artwork %v", id)
	logext.LogError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	assets := []storage.Asset{}
	path, err = pathext.FormatWorkPath(path, work)
	if err == nil {
		err = storage.StoreWork(work, assets, path)
	}
	logext.LogSuccess(err, "stored metadata for artwork %v in %v", id, path)
	logext.LogError(err, "failed to store metadata for artwork %v", id)
	return work, err
}

func (d *Downloader) DownloadArtwork(id uint64, path string, size ImageSize) (*work.Work, error) {
	fetchedWork, err := d.fetchArtworkMeta(id)
	logext.LogSuccess(err, "fetched metadata for artwork %v", id)
	logext.LogError(err, "failed to fetch metadata for artwork %v", id)
	if err != nil {
		return nil, err
	}

	if fetchedWork.Kind == work.KindUgoira {
		err = d.continueUgoira(fetchedWork, id, path)
	} else {
		err = d.continueIllustOrManga(fetchedWork, id, size, path)
	}
	if err != nil {
		return nil, err
	}

	return fetchedWork, nil
}

func (d *Downloader) continueUgoira(work *work.Work, id uint64, path string) error {
	data, frames, err := d.fetchArtworkFrames(id)
	logext.LogSuccess(err, "fetched frames data for artwork %v", id)
	logext.LogError(err, "failed to fetch frames data for artwork %v", id)
	if err != nil {
		return err
	}

	archive, err := d.fetch(data)
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
	path, err = pathext.FormatWorkPath(path, work)
	if err == nil {
		err = storage.StoreWork(work, assets, path)
	}
	logext.LogSuccess(err, "stored files for artwork %v in %v", id, path)
	logext.LogError(err, "failed to store files for artwork %v", id)
	return err
}

func (d *Downloader) continueIllustOrManga(
	work *work.Work, id uint64, size ImageSize, path string,
) error {
	pages, err := d.fetchArtworkUrls(id)
	logext.LogSuccess(err, "fetched page urls for artwork %v", id)
	logext.LogError(err, "failed to fetch page urls for artwork %v", id)
	if err != nil {
		return err
	}

	assetChannel := make(chan storage.Asset, len(pages))
	errorChannel := make(chan error, len(pages))
	for i, urls := range pages {
		go func() {
			url := urls[size]
			bytes, err := d.fetch(url)
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
		}()
	}

	assets := make([]storage.Asset, len(pages))
	for i := range pages {
		assets[i] = <-assetChannel
		err = <-errorChannel
		if err != nil {
			return err
		}
	}

	path, err = pathext.FormatWorkPath(path, work)
	if err == nil {
		err = storage.StoreWork(work, assets, path)
	}
	logext.LogSuccess(err, "stored files for artwork %v in %v", id, path)
	logext.LogError(err, "failed to store files for artwork %v", id)
	return err
}
