package downloader

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/encode"
	"github.com/fekoneko/piximan/pkg/logext"
	"github.com/fekoneko/piximan/pkg/storage"
)

type Downloader struct {
	sessionId string
	client    http.Client
}

func New(sessionId string) *Downloader {
	url, _ := url.Parse("https://www.pixiv.net")
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(url, []*http.Cookie{
		{Name: "PHPSESSID", Value: sessionId},
	})
	client := http.Client{Jar: jar}
	return &Downloader{sessionId, client}
}

func (d *Downloader) DownloadWork(id uint64, size ImageSize, path string) (*work.Work, error) {
	fetchedWork, err := d.fetchArtwork(id)
	logext.LogSuccess(err, "fetched metadata for work %v", id)
	logext.LogError(err, "failed to fetch metadata for work %v", id)
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
	logext.LogSuccess(err, "fetched frames data for work %v", id)
	logext.LogError(err, "failed to fetch frames data for work %v", id)
	if err != nil {
		return err
	}

	archive, err := d.fetch(data)
	logext.LogSuccess(err, "fetched frames for work %v", id)
	logext.LogError(err, "failed to fetch frames for work %v", id)
	if err != nil {
		return err
	}

	gif, err := encode.GifFromFrames(archive, frames)
	logext.LogSuccess(err, "encoded frames for work %v", id)
	logext.LogError(err, "failed to encode frames for work %v", id)
	if err != nil {
		return err
	}

	assets := []storage.Asset{{Bytes: gif, Extension: ".gif"}}
	err = storage.StoreWork(work, assets, path)
	logext.LogSuccess(err, "wrote files for work %v", id)
	logext.LogError(err, "failed to write files for work %v", id)
	return err
}

func (d *Downloader) continueIllustOrManga(work *work.Work, id uint64, size ImageSize, path string) error {
	pages, err := d.fetchArtworkUrls(id)
	logext.LogSuccess(err, "fetched page urls for work %v", id)
	logext.LogError(err, "failed to fetch page urls for work %v", id)
	if err != nil {
		return err
	}

	assetChannel := make(chan storage.Asset, len(pages))
	errorChannel := make(chan error, len(pages))
	for i, urls := range pages {
		go func() {
			url := urls[size]
			bytes, err := d.fetch(url)
			logext.LogSuccess(err, "fetched page %v for work %v", i, id)
			logext.LogError(err, "failed to fetch page %v for work %v", i, id)
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

	err = storage.StoreWork(work, assets, path)
	logext.LogSuccess(err, "wrote files for work %v", id)
	logext.LogError(err, "failed to write files for work %v", id)
	return err
}
