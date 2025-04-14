package downloader

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/downloader/queue"
)

const CHANNEL_SIZE = 10
const PENDING_CAP = 2

type Downloader struct {
	client          http.Client
	channel         chan *work.Work
	queue           queue.Queue
	numPending      int
	numPendingMutex sync.Mutex
}

type AuthorizedDownloader Downloader

func New() *Downloader {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	channel := make(chan *work.Work, CHANNEL_SIZE)
	return &Downloader{client, channel, queue.Queue{}, 0, sync.Mutex{}}
}

// TODO: use this one to get bookmarked works
func NewAuthorized(sessionId string) *AuthorizedDownloader {
	url, _ := url.Parse("https://www.pixiv.net")
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(url, []*http.Cookie{
		{Name: "PHPSESSID", Value: sessionId},
	})
	client := http.Client{Jar: jar}
	channel := make(chan *work.Work, CHANNEL_SIZE)
	return &AuthorizedDownloader{client, channel, queue.Queue{}, 0, sync.Mutex{}}
}
