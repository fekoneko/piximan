package downloader

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

func (d *Downloader) DownloadWork(id uint64, path string) error {
	work, err := d.fetchWork(id)
	if err != nil {
		return err
	}

	fmt.Println(work)

	return nil
}
