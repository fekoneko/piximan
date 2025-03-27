package downloader

import (
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
