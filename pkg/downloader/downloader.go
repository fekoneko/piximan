package downloader

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Downloader struct {
	client http.Client
}

type AuthorizedDownloader struct {
	Downloader
	sessionId string
}

func New() *Downloader {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	return &Downloader{client}
}

func NewAuthorized(sessionId string) *AuthorizedDownloader {
	url, _ := url.Parse("https://www.pixiv.net")
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(url, []*http.Cookie{
		{Name: "PHPSESSID", Value: sessionId},
	})
	client := http.Client{Jar: jar}
	return &AuthorizedDownloader{Downloader{client}, sessionId}
}
