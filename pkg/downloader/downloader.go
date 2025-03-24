package downloader

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/fekoneko/piximan/pkg/collection/work"
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

func (d *Downloader) DownloadWork(id uint64, path string) (*work.Work, error) {
	work, err := d.fetchWork(id)
	if err != nil {
		return nil, err
	}

	if err := storeWork(work, path); err != nil {
		return nil, err
	}

	return work, nil
}
