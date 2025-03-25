package downloader

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/fekoneko/piximan/pkg/collection/work"
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

func (d *Downloader) DownloadWork(id uint64, size ImageSize, path string) *work.Work {
	workChannel := make(chan *work.Work)
	go func() {
		work, err := d.fetchWork(id)
		logext.LogSuccess(err, "fetched metadata for work %v", id)
		logext.LogError(err, "failed to fetch metadata for work %v", id)
		workChannel <- work
	}()

	pagesChannel := make(chan [][4]string)
	go func() {
		pages, err := d.fetchPages(id)
		logext.LogSuccess(err, "fetched page urls for work %v", id)
		logext.LogError(err, "failed to fetch page urls for work %v", id)
		pagesChannel <- pages
	}()

	pages := <-pagesChannel

	imageChannel := make(chan storage.Image, len(pages))
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
			image := storage.Image{Bytes: bytes, Extension: extension}
			imageChannel <- image
		}()
	}

	work := <-workChannel

	images := make([]storage.Image, len(pages))
	for i := range pages {
		images[i] = <-imageChannel
	}

	err := storage.StoreWork(work, images, path)
	logext.LogSuccess(err, "wrote files for work %v", id)
	logext.LogError(err, "failed to write files for work %v", id)

	return work
}
