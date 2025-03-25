package downloader

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/fekoneko/piximan/pkg/collection/work"
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
	// TODO: async
	// TODO: accomulate errors

	work, err := d.fetchWork(id)
	if err != nil {
		return nil, err
	}

	pages, err := d.fetchPages(id)
	if err != nil {
		return nil, err
	}

	var images []storage.Image
	for _, urls := range pages {
		url := urls[size]
		bytes, err := d.fetch(url)
		if err != nil {
			return nil, err
		}
		dotIndex := strings.LastIndex(url, ".")
		var extension string
		if dotIndex != -1 {
			extension = urls[size][dotIndex:]
		}
		image := storage.Image{Bytes: bytes, Extension: extension}
		images = append(images, image)
	}

	if err := storage.StoreWork(work, images, path); err != nil {
		return nil, err
	}

	return work, nil
}
